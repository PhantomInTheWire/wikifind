package indexer

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWikiXMLParser_Parse(t *testing.T) {
	tempDir := t.TempDir()
	indexPath := filepath.Join(tempDir, "index")

	parser := NewWikiXMLParser(indexPath)

	xmlFile := "../cmd/test_data.xml"

	ctx := context.Background()
	err := parser.Parse(ctx, xmlFile)
	require.NoError(t, err)
}
