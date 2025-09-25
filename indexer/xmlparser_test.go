package indexer

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWikiXMLParser_Parse(t *testing.T) {
	tempDir := t.TempDir()
	indexPath := filepath.Join(tempDir, "index")

	parser := NewWikiXMLParser(indexPath)

	xmlFile := "../cmd/test_data.xml"

	err := parser.Parse(xmlFile)
	require.NoError(t, err)
}
