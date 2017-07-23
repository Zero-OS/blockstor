package ardb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testBackendReadWrite(ctx context.Context, t *testing.T, vdiskID string, blockSize int64, size uint64, storage backendStorage) {
	if !assert.NotNil(t, storage) {
		return
	}

	vComp := new(vdiskCompletion)
	backend := newBackend(vdiskID, blockSize, size, storage, nil, vComp)
	if !assert.NotNil(t, backend) {
		return
	}
	go backend.GoBackground(ctx)
	defer backend.Close(ctx)

	someContent := make([]byte, blockSize)
	for i := range someContent {
		someContent[i] = byte(i % 255)
	}

	nilContent := make([]byte, blockSize)

	// ensure the content doest not exist yet
	payload, err := backend.ReadAt(ctx, 0, blockSize)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, nilContent, payload) {
		return
	}
	payload, err = backend.ReadAt(ctx, blockSize, blockSize)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, nilContent, payload) {
		return
	}

	// write first block
	bw, err := backend.WriteAt(ctx, someContent, 0)
	if !assert.NoError(t, err) || !assert.Equal(t, blockSize, bw) {
		return
	}
	// first block should now exist
	payload, err = backend.ReadAt(ctx, 0, blockSize)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, someContent, payload) {
		return
	}
	// second block should still not exist
	payload, err = backend.ReadAt(ctx, blockSize, blockSize)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, nilContent, payload) {
		return
	}

	// write second block
	bw, err = backend.WriteAt(ctx, someContent, blockSize)
	if !assert.NoError(t, err) || !assert.Equal(t, blockSize, bw) {
		return
	}

	// both blocks should now exist
	payload, err = backend.ReadAt(ctx, 0, blockSize)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, someContent, payload) {
		return
	}
	payload, err = backend.ReadAt(ctx, blockSize, blockSize)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, someContent, payload) {
		return
	}

	midBlockOffset := blockSize / 2
	midBlockLength := blockSize - midBlockOffset

	// Let's now erase half of the first block by writing zeroes
	bytesWritten, err := backend.WriteZeroesAt(ctx, midBlockOffset, midBlockLength)
	if !assert.NoError(t, err) || !assert.Equal(t, midBlockLength, bytesWritten) {
		return
	}

	// getting the first block should now result in half of them being zeroes
	halfNilBlock := make([]byte, blockSize)
	copy(halfNilBlock[:midBlockOffset], someContent)
	payload, err = backend.ReadAt(ctx, 0, blockSize)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, halfNilBlock, payload) {
		return
	}

	// let's now delete the content by writing just zeroes
	bytesWritten, err = backend.WriteZeroesAt(ctx, 0, blockSize)
	if !assert.NoError(t, err) || !assert.Equal(t, blockSize, bytesWritten) {
		return
	}
	// first block should now be deleted
	payload, err = backend.ReadAt(ctx, 0, blockSize)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, nilContent, payload) {
		return
	}

	// let's now try to merge a full block (the 2nd one) with some other content
	bytesWritten, err = backend.WriteAt(ctx, someContent[:midBlockLength], blockSize+midBlockOffset)
	if !assert.NoError(t, err) || !assert.Equal(t, midBlockLength, bytesWritten) {
		return
	}
	// getting the content should now the first one
	mergedBlock := make([]byte, blockSize)
	copy(mergedBlock[:midBlockOffset], someContent)
	copy(mergedBlock[midBlockOffset:], someContent)
	payload, err = backend.ReadAt(ctx, blockSize, blockSize)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, mergedBlock, payload) {
		return
	}
}