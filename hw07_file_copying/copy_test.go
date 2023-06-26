package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const dataPath = "testdata"

var srcPath = filepath.Join(dataPath, "input.txt")

func TestCopy(t *testing.T) {
	for _, tst := range []struct {
		offset int64
		limit  int64
	}{
		{offset: 0, limit: 0},
		{offset: 0, limit: 10},
		{offset: 0, limit: 1000},
		{offset: 0, limit: 10000},
		{offset: 100, limit: 1000},
		{offset: 6000, limit: 1000},
	} {
		tst := tst
		t.Run(fmt.Sprintf("%v_%v", offset, limit), func(t *testing.T) {
			dst, err := os.CreateTemp(os.TempDir(), "go-copy")
			require.NoError(t, err)
			defer dst.Close()
			defer os.Remove(dst.Name())
			err = Copy(srcPath, dst.Name(), tst.offset, tst.limit)

			require.NoError(t, err)
			expPath := filepath.Join(dataPath, fmt.Sprintf("out_offset%d_limit%d.txt", tst.offset, tst.limit))
			fmt.Println(expPath)
			expContent, err := os.ReadFile(expPath)
			require.NoError(t, err)
			dstContent, err := os.ReadFile(dst.Name())
			require.NoError(t, err)
			require.Zero(t, bytes.Compare(expContent, dstContent))
		})
	}

	t.Run("offset exceeds file size", func(t *testing.T) {
		offset := 70000
		limit := 0

		dst, err := os.CreateTemp(os.TempDir(), "go-copy")
		require.NoError(t, err)
		defer dst.Close()
		defer os.Remove(dst.Name())
		err = Copy(srcPath, dst.Name(), int64(offset), int64(limit))
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("err same file", func(t *testing.T) {
		offset := 0
		limit := 0

		dst, err := os.CreateTemp(os.TempDir(), "go-copy")
		require.NoError(t, err)
		defer dst.Close()
		defer os.Remove(dst.Name())
		err = Copy(dst.Name(), dst.Name(), int64(offset), int64(limit))
		require.ErrorIs(t, err, ErrSameFile)
	})

	t.Run("block device copy support", func(t *testing.T) {
		srcPath := "/dev/urandom"
		offset := 0
		limit := 10_000

		dst, err := os.CreateTemp(os.TempDir(), "go-copy")
		require.NoError(t, err)
		defer dst.Close()
		defer os.Remove(dst.Name())
		err = Copy(srcPath, dst.Name(), int64(offset), int64(limit))
		require.NoError(t, err)
	})

	t.Run("block device copy", func(t *testing.T) {
		srcPath := "/dev/urandom"
		offset := 0
		limit := 0

		dst, err := os.CreateTemp(os.TempDir(), "go-copy")
		require.NoError(t, err)
		defer dst.Close()
		defer os.Remove(dst.Name())
		err = Copy(srcPath, dst.Name(), int64(offset), int64(limit))
		fmt.Println(err)
		require.ErrorIs(t, err, ErrNoLimitedDeviceOperation)
	})
}
