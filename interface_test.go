package gounrar

import (
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestSeek(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		pos           int64
		expectFile    string
		expectContent string
		expectError   bool
	}{
		{
			name:          "rar",
			path:          "testdata/test-rar.rar",
			pos:           67,
			expectFile:    "中文.txt",
			expectContent: "world",
		},
		{
			name:          "rar5",
			path:          "testdata/test-rar5.rar",
			pos:           23,
			expectFile:    "dir/2.txt",
			expectContent: "world",
		},
		{
			name:        "seek error",
			path:        "testdata/test-rar5.rar",
			pos:         20,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := Open(tt.path)
			require.NoError(t, err)
			hdr, err := a.SeekPos(tt.pos)
			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectFile, hdr.FileName)
			data, err := a.ReadAll()
			require.Equal(t, tt.expectContent, string(data))
		})
	}
}

func TestOpenNotRARFile(t *testing.T) {
	_, err := Open("testdata/plain.txt")
	require.Error(t, err)
}

func TestSimple(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectFiles map[string]string
		expectDirs  []string
	}{
		{
			name: "rar",
			path: "testdata/test-rar.rar",
			expectFiles: map[string]string{
				"1.txt":  "hello",
				"中文.txt": "world",
			},
		},
		{
			name:       "rar5",
			path:       "testdata/test-rar5.rar",
			expectDirs: []string{"dir"},
			expectFiles: map[string]string{
				"1.txt":     "hello",
				"dir/2.txt": "world",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := Open(tt.path)
			require.NoError(t, err)
			//goland:noinspection GoUnhandledErrorResult
			defer a.Close()
			actualFiles := make(map[string]string)
			var actualDirs []string

			for {
				hdr, err := a.Next()
				if err == io.EOF {
					break
				}

				t.Logf("hdr %s blk %d", hdr.FileName, hdr.BlockPos)
				require.NoError(t, err)
				if hdr.IsDir() {
					actualDirs = append(actualDirs, hdr.FileName)
				} else {
					data, err := a.ReadAll()
					require.NoError(t, err)
					actualFiles[hdr.FileName] = string(data)
				}
			}

			require.Equal(t, tt.expectDirs, actualDirs)
			require.Equal(t, tt.expectFiles, actualFiles)
		})
	}
}
