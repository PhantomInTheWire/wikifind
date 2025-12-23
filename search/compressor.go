package search

import (
	"bytes"
	"compress/gzip"
	"io"
)

type Compressor struct{}

func (c *Compressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(data); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer func() { _ = gz.Close() }()

	return io.ReadAll(gz)
}
