package backup

import (
	"sync"

	"github.com/zero-os/0-Disk/testdata"
)

var validConfigs = []Config{
	// Explicit Example
	Config{
		VdiskID:         "foo",
		SnapshotID:      "foo",
		BlockSize:       DefaultBlockSize,
		JobCount:        0,
		CompressionType: LZ4Compression,
		CryptoKey:       CryptoKey{4, 2},
	},
	// implicit version of first example
	Config{
		VdiskID:   "foo",
		CryptoKey: CryptoKey{4, 2},
	},
	// full (FTP) example
	Config{
		VdiskID:         "foo",
		SnapshotID:      "bar",
		BlockSize:       4096,
		JobCount:        1,
		CompressionType: XZCompression,
		CryptoKey: CryptoKey{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			0, 1},
	},
}

var invalidConfigs = []Config{
	// Nothing Given,
	Config{},
	// Invalid BlockSize
	Config{
		VdiskID:   "foo",
		BlockSize: 2000,
	},
	// Missing VdiskID
	Config{},
}

func getLedeImageBlocks() map[int64][]byte {
	fetchLedeImageBlocksOnce.Do(func() {
		var err error
		ledeImageBlocks, err = testdata.ReadAllLedeBlocks()
		if err != nil {
			panic(err)
		}
	})

	return ledeImageBlocks
}

var fetchLedeImageBlocksOnce sync.Once
var ledeImageBlocks map[int64][]byte
