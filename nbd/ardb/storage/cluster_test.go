package storage

import (
	"context"
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zero-os/0-Disk/config"
	"github.com/zero-os/0-Disk/errors"
	"github.com/zero-os/0-Disk/log"
	"github.com/zero-os/0-Disk/nbd/ardb"
	"github.com/zero-os/0-Disk/redisstub"
)

// TODO:
// Test for NewPrimaryCluster:
//   + that we get panic when no primary cluster is defined when we create it;
// Test for NewSlaveCluster and NewTemplateCluster:
//   + that we can create an optional cluster with no servers, and hot-configure them afterwards;
// Test for NewPrimaryCluster, NewSlaveCluster and NewTemplateCluster:
//   + that we can update the configuration of servers;
//   + that we can hot-swap 2 online servers (and automatically copy the data from the old one to the new one);
//      + add also a test to ensure that there cannot be a window between the hot-swap
//        and another command still writing to the previous server
//   + that we cannot grow or shrink a cluster in terms of server size;

func TestPrimaryServerRIP(t *testing.T) {
	mr := redisstub.NewMemoryRedis()
	defer mr.Close()

	const (
		vdiskID    = "foo"
		clusterID  = "foo"
		blockSize  = 8
		blockCount = 8
	)

	source := config.NewStubSource()
	sourceClusterConfig := config.StorageClusterConfig{
		Servers: []config.StorageServerConfig{
			mr.StorageServerConfig(),
			config.StorageServerConfig{State: config.StorageServerStateRIP},
		},
	}
	source.SetPrimaryStorageCluster(vdiskID, clusterID, &sourceClusterConfig)

	ctx := context.Background()
	require := require.New(t)

	cluster, err := NewPrimaryCluster(ctx, vdiskID, source)
	require.NoError(err)
	defer cluster.Close()

	// NonDedupedStorage is the easiest to use for this kind of testing purpose
	storage, err := NonDeduped(vdiskID, "", blockSize, cluster, nil)
	require.NoError(err)
	defer storage.Close()

	var contentSlice [][]byte

	// store blocks, this should be fine
	for index := int64(0); index < blockCount; index++ {
		content := make([]byte, blockSize)
		rand.Read(content)
		contentSlice = append(contentSlice, content)

		err = storage.SetBlock(index, content)
		require.NoError(err)
	}

	// load blocks, this should  be fine as well
	for index := int64(0); index < blockCount; index++ {
		content, err := storage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}
}

func TestPrimaryServerFails(t *testing.T) {
	slice := redisstub.NewMemoryRedisSlice(2)
	defer slice.Close()

	const (
		vdiskID    = "foo"
		clusterID  = "foo"
		blockSize  = 8
		blockCount = 8
	)

	source := config.NewStubSource()
	sourceClusterConfig := slice.StorageClusterConfig()
	source.SetPrimaryStorageCluster(vdiskID, clusterID, &sourceClusterConfig)

	ctx := context.Background()
	require := require.New(t)

	cluster, err := NewPrimaryCluster(ctx, vdiskID, source)
	require.NoError(err)
	defer cluster.Close()

	// NonDedupedStorage is the easiest to use for this kind of testing purpose
	storage, err := NonDeduped(vdiskID, "", blockSize, cluster, nil)
	require.NoError(err)
	defer storage.Close()

	var contentSlice [][]byte

	// store blocks, this should all still be fine
	for index := int64(0); index < blockCount; index++ {
		content := make([]byte, blockSize)
		rand.Read(content)
		contentSlice = append(contentSlice, content)

		err = storage.SetBlock(index, content)
		require.NoError(err)
	}

	// load blocks, this should all still be fine as well
	for index := int64(0); index < blockCount; index++ {
		content, err := storage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}

	// now let's disable the 2nd server
	slice.CloseServer(1)

	// getting all content from the 1st server should still work
	for index := int64(0); index < blockCount; index += 2 {
		content, err := storage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}

	// getting all content from the 2nd server should however no longer work
	for index := int64(1); index < blockCount; index += 2 {
		content, err := storage.GetBlock(index)
		require.Equal(ardb.ErrServerUnavailable, err)
		require.Nil(content)
	}
}

func TestPrimaryServerFailsByNotification(t *testing.T) {
	slice := redisstub.NewMemoryRedisSlice(2)
	defer slice.Close()

	const (
		vdiskID    = "foo"
		clusterID  = "foo"
		blockSize  = 8
		blockCount = 8
	)

	source := config.NewStubSource()
	sourceClusterConfig := slice.StorageClusterConfig()
	source.SetPrimaryStorageCluster(vdiskID, clusterID, &sourceClusterConfig)

	ctx := context.Background()
	require := require.New(t)

	cluster, err := NewPrimaryCluster(ctx, vdiskID, source)
	require.NoError(err)
	defer cluster.Close()

	// NonDedupedStorage is the easiest to use for this kind of testing purpose
	storage, err := NonDeduped(vdiskID, "", blockSize, cluster, nil)
	require.NoError(err)
	defer storage.Close()

	var contentSlice [][]byte

	// store blocks, this should all still be fine
	for index := int64(0); index < blockCount; index++ {
		content := make([]byte, blockSize)
		rand.Read(content)
		contentSlice = append(contentSlice, content)

		err = storage.SetBlock(index, content)
		require.NoError(err)
	}

	// load blocks, this should all still be fine as well
	for index := int64(0); index < blockCount; index++ {
		content, err := storage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}

	// now let's disable the 2nd primary server
	sourceClusterConfig.Servers[1].State = config.StorageServerStateOffline
	source.SetStorageCluster(clusterID, &sourceClusterConfig)
	waitForAsyncClusterUpdate(t, func() bool {
		state, err := cluster.controller.ServerStateAt(1)
		require.NoError(err)
		return state.Config.State == config.StorageServerStateOffline
	})

	// getting all content from the 1st server should still work
	for index := int64(0); index < blockCount; index += 2 {
		content, err := storage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}

	// getting all content from the 2nd server should however no longer work
	for index := int64(1); index < blockCount; index += 2 {
		content, err := storage.GetBlock(index)
		require.Equal(ardb.ErrServerUnavailable, err)
		require.Nil(content)
	}
}

func TestTemplateServerFails(t *testing.T) {
	slice := redisstub.NewMemoryRedisSlice(2)
	defer slice.Close()

	const (
		vdiskID           = "foo"
		clusterID         = "foo"
		templateClusterID = "bar"
		blockSize         = 8
		blockCount        = 8
	)

	source := config.NewStubSource()
	sourceClusterConfig := slice.StorageClusterConfig()
	source.SetPrimaryStorageCluster(vdiskID, clusterID, &sourceClusterConfig)

	templateSlice := redisstub.NewMemoryRedisSlice(2)
	defer templateSlice.Close()

	templateClusterConfig := templateSlice.StorageClusterConfig()
	source.SetTemplateStorageCluster(vdiskID, templateClusterID, &templateClusterConfig)

	ctx := context.Background()
	require := require.New(t)

	cluster, err := NewPrimaryCluster(ctx, vdiskID, source)
	require.NoError(err)
	defer cluster.Close()

	templateCluster, err := NewTemplateCluster(ctx, vdiskID, false, source)
	require.NoError(err)
	defer templateCluster.Close()

	// NonDedupedStorage is the easiest to use for this kind of testing purpose
	storage, err := NonDeduped(vdiskID, "", blockSize, cluster, templateCluster)
	require.NoError(err)
	defer storage.Close()
	templateStorage, err := NonDeduped(vdiskID, "", blockSize, templateCluster, nil)
	require.NoError(err)
	defer templateStorage.Close()

	var contentSlice [][]byte

	// store blocks in template, this should all still be fine
	for index := int64(0); index < blockCount; index++ {
		content := make([]byte, blockSize)
		rand.Read(content)
		contentSlice = append(contentSlice, content)

		err = templateStorage.SetBlock(index, content)
		require.NoError(err)
	}

	// load blocks, this should all still be fine as well
	for index := int64(0); index < blockCount; index++ {
		content, err := templateStorage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}

	// get blocks from template, this should all still be fine as well
	for index := int64(0); index < blockCount; index++ {
		content, err := storage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}

	// now let's disable the 2nd template server
	templateSlice.CloseServer(1)

	// create new primary servers
	newSlice := redisstub.NewMemoryRedisSlice(2)
	defer newSlice.Close()

	pcc := cluster.controller.(*singleClusterStateController)
	pcc.mux.Lock()
	pcc.servers = newSlice.StorageClusterConfig().Servers
	pcc.serverCount = int64(len(pcc.servers))
	pcc.mux.Unlock()

	// getting all content from the 1st server should still work
	for index := int64(0); index < blockCount; index += 2 {
		content, err := storage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}

	// getting all content from the 2nd server should however no longer work
	for index := int64(1); index < blockCount; index += 2 {
		content, err := storage.GetBlock(index)
		require.Equal(ardb.ErrServerUnavailable, err)
		require.Nil(content)
	}
}

func TestTemplateServerFailsByNotification(t *testing.T) {
	slice := redisstub.NewMemoryRedisSlice(2)
	defer slice.Close()

	const (
		vdiskID           = "foo"
		clusterID         = "foo"
		templateClusterID = "bar"
		blockSize         = 8
		blockCount        = 8
	)

	source := config.NewStubSource()
	sourceClusterConfig := slice.StorageClusterConfig()
	source.SetPrimaryStorageCluster(vdiskID, clusterID, &sourceClusterConfig)

	templateSlice := redisstub.NewMemoryRedisSlice(2)
	defer templateSlice.Close()

	templateClusterConfig := templateSlice.StorageClusterConfig()
	source.SetTemplateStorageCluster(vdiskID, templateClusterID, &templateClusterConfig)

	ctx := context.Background()
	require := require.New(t)

	cluster, err := NewPrimaryCluster(ctx, vdiskID, source)
	require.NoError(err)
	defer cluster.Close()

	templateCluster, err := NewTemplateCluster(ctx, vdiskID, false, source)
	require.NoError(err)
	defer templateCluster.Close()

	// NonDedupedStorage is the easiest to use for this kind of testing purpose
	storage, err := NonDeduped(vdiskID, "", blockSize, cluster, templateCluster)
	require.NoError(err)
	defer storage.Close()
	templateStorage, err := NonDeduped(vdiskID, "", blockSize, templateCluster, nil)
	require.NoError(err)
	defer templateStorage.Close()

	var contentSlice [][]byte

	// store blocks in template, this should all still be fine
	for index := int64(0); index < blockCount; index++ {
		content := make([]byte, blockSize)
		rand.Read(content)
		contentSlice = append(contentSlice, content)

		err = templateStorage.SetBlock(index, content)
		require.NoError(err)
	}

	// load blocks, this should all still be fine as well
	for index := int64(0); index < blockCount; index++ {
		content, err := templateStorage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}

	// get blocks from template, this should all still be fine as well
	for index := int64(0); index < blockCount; index++ {
		content, err := storage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}

	// now let's disable the 2nd template server
	templateClusterConfig.Servers[1].State = config.StorageServerStateOffline
	source.SetStorageCluster(templateClusterID, &templateClusterConfig)
	waitForAsyncClusterUpdate(t, func() bool {
		cfg, err := templateCluster.controller.ServerStateAt(1)
		require.NoError(err)
		return cfg.Config.State == config.StorageServerStateOffline
	})

	// create new primary servers
	newSlice := redisstub.NewMemoryRedisSlice(2)
	defer newSlice.Close()

	pcc := cluster.controller.(*singleClusterStateController)
	pcc.mux.Lock()
	pcc.servers = newSlice.StorageClusterConfig().Servers
	pcc.serverCount = int64(len(pcc.servers))
	pcc.mux.Unlock()

	// getting all content from the 1st server should still work
	for index := int64(0); index < blockCount; index += 2 {
		content, err := storage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}

	// getting all content from the 2nd server should however no longer work
	for index := int64(1); index < blockCount; index += 2 {
		content, err := storage.GetBlock(index)
		require.Equal(ardb.ErrServerUnavailable, err)
		require.Nil(content)
	}
}

func TestTemplateServerRIP(t *testing.T) {
	slice := redisstub.NewMemoryRedisSlice(2)
	defer slice.Close()

	const (
		vdiskID           = "foo"
		clusterID         = "foo"
		templateClusterID = "bar"
		blockSize         = 8
		blockCount        = 8
	)

	source := config.NewStubSource()
	sourceClusterConfig := slice.StorageClusterConfig()
	source.SetPrimaryStorageCluster(vdiskID, clusterID, &sourceClusterConfig)

	templateMR := redisstub.NewMemoryRedis()
	defer templateMR.Close()

	templateClusterConfig := config.StorageClusterConfig{
		Servers: []config.StorageServerConfig{
			templateMR.StorageServerConfig(),
			config.StorageServerConfig{State: config.StorageServerStateRIP},
		},
	}
	source.SetTemplateStorageCluster(vdiskID, templateClusterID, &templateClusterConfig)

	ctx := context.Background()
	require := require.New(t)

	cluster, err := NewPrimaryCluster(ctx, vdiskID, source)
	require.NoError(err)
	defer cluster.Close()

	templateCluster, err := NewTemplateCluster(ctx, vdiskID, false, source)
	require.NoError(err)
	defer templateCluster.Close()

	// NonDedupedStorage is the easiest to use for this kind of testing purpose
	storage, err := NonDeduped(vdiskID, "", blockSize, cluster, templateCluster)
	require.NoError(err)
	defer storage.Close()

	templateStorage, err := NonDeduped(vdiskID, "", blockSize, templateCluster, nil)
	require.NoError(err)
	defer templateStorage.Close()

	var contentSlice [][]byte

	// store blocks in template storage, this should be fine
	for index := int64(0); index < blockCount; index++ {
		content := make([]byte, blockSize)
		rand.Read(content)
		contentSlice = append(contentSlice, content)

		err = templateStorage.SetBlock(index, content)
		require.NoError(err)
	}

	// load blocks in primary, this should be fine as well
	for index := int64(0); index < blockCount; index++ {
		content, err := storage.GetBlock(index)
		require.NoError(err)
		require.Equal(contentSlice[index], content)
	}
}

func waitForAsyncClusterUpdate(t *testing.T, predicate func() bool) {
	timeoutTicker := time.NewTicker(30 * time.Second)
	pollTicker := time.NewTicker(5 * time.Millisecond)

	for {
		select {
		case <-pollTicker.C:
			if predicate() {
				return
			}

		case <-timeoutTicker.C:
			t.Fatal("Timed out waiting for tlog cluster ID to be updated.")
		}
	}
}

func TestMapErrorToBroadcastStatus(t *testing.T) {
	assert := assert.New(t)

	// unknown errors return true
	status := MapErrorToBroadcastStatus(errors.New("foo"))
	assert.Equal(log.StatusUnknownError, status)

	// all possible sucesfull map scenarios:

	// map EOF errors
	status = MapErrorToBroadcastStatus(io.EOF)
	assert.Equal(log.StatusServerDisconnect, status)

	// map net.Errors
	status = MapErrorToBroadcastStatus(stubNetError{false, false})
	assert.Equal(log.StatusUnknownError, status)

	status = MapErrorToBroadcastStatus(stubNetError{false, true})
	assert.Equal(log.StatusServerTempError, status)

	status = MapErrorToBroadcastStatus(stubNetError{true, false})
	assert.Equal(log.StatusServerTimeout, status)

	status = MapErrorToBroadcastStatus(stubNetError{true, true})
	assert.Equal(log.StatusServerTimeout, status)
}

type stubNetError struct {
	timeout, temporary bool
}

func (err stubNetError) Timeout() bool {
	return err.timeout
}

func (err stubNetError) Temporary() bool {
	return err.temporary
}

func (err stubNetError) Error() string {
	return "stub net error"
}
