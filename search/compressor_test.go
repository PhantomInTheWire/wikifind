package search

import "testing"

func TestCompressor_CompressDecompress(t *testing.T) {
	c := &Compressor{}

	original := []byte("Hello, world! This is a test string for compression.")

	compressed, err := c.Compress(original)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	decompressed, err := c.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if string(decompressed) != string(original) {
		t.Error("Decompressed data does not match original")
	}
}
