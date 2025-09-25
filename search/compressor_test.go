package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompressor_CompressDecompress(t *testing.T) {
	tests := []struct {
		name     string
		original []byte
	}{
		{"short string", []byte("Hello, world!")},
		{"long string", []byte("Hello, world! This is a test string for compression. It should be longer to test compression effectiveness.")},
		{"empty", []byte("")},
		{"single byte", []byte("a")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compressor{}

			compressed, err := c.Compress(tt.original)
			require.NoError(t, err)

			decompressed, err := c.Decompress(compressed)
			require.NoError(t, err)

			assert.Equal(t, tt.original, decompressed)
		})
	}
}
