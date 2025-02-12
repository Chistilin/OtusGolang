package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestCopy(t *testing.T) {
	input := "testdata/input.txt"
	dir, _ := os.MkdirTemp("", "temp")
	destination, _ := os.CreateTemp(dir, "temp")
	defer os.RemoveAll(dir)

	tests := []struct {
		text      string
		offset    int64
		limit     int64
		input     string
		checkFile string
		err       error
	}{
		{offset: 0, limit: 0, input: input, checkFile: "testdata/out_offset0_limit0.txt"},
		{offset: 0, limit: 10, input: input, checkFile: "testdata/out_offset0_limit10.txt"},
		{offset: 100, limit: 1000, input: input, checkFile: "testdata/out_offset100_limit1000.txt"},
		{offset: 0, limit: 10000, input: input, checkFile: "testdata/out_offset0_limit10000.txt"},
		{offset: 0, limit: 1000, input: input, checkFile: "testdata/out_offset0_limit1000.txt"},
		{offset: 6000, limit: 1000, input: input, checkFile: "testdata/out_offset6000_limit1000.txt"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("Offset: %v, limit %v %s", tc.offset, tc.limit, tc.text), func(t *testing.T) {
			err := Copy(input, destination.Name(), tc.offset, tc.limit)

			if tc.err == nil {
				require.Nil(t, err)

				srcContent, _ := os.ReadFile(tc.checkFile)
				destContent, _ := os.ReadFile(destination.Name())

				require.Equal(t, string(srcContent), string(destContent))
			} else {
				require.Error(t, tc.err, err)
			}
		})
	}

	t.Run("Offset exceeds file size", func(t *testing.T) {
		err := Copy(input, destination.Name(), 10000, 0)
		require.Error(t, err, ErrUnsupportedFile)
	})

	t.Run("FromPath is undefined", func(t *testing.T) {
		err := Copy("", destination.Name(), 0, 0)
		require.Error(t, err, ErrUnsupportedFile)
	})

	t.Run("ToPath is undefined", func(t *testing.T) {
		err := Copy(input, "", 0, 0)
		require.Error(t, err, ErrUnsupportedFile)
	})

	t.Run("Offset is negative", func(t *testing.T) {
		err := Copy(input, destination.Name(), -1, 0)
		require.Error(t, err, ErrNegativeOffsetSize)
	})

	t.Run("Limit is negative", func(t *testing.T) {
		err := Copy(input, destination.Name(), 0, -1)
		require.Error(t, err, ErrNegativeLimit)
	})

	t.Run("Unsupported input file", func(t *testing.T) {
		err := Copy("/dev/urandom", destination.Name(), 0, 0)

		require.Equal(t, ErrUnsupportedFile, err)
	})
}
