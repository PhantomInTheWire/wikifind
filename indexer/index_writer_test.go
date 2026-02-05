package indexer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexWriter_WriteIndex(t *testing.T) {
	tempDir := t.TempDir()
	indexPath := filepath.Join(tempDir, "index")

	idx := NewInvertedIndex()
	// Add terms starting with different letters
	idx.Add("apple", "doc1", Posting{Fields: TITLE, Frequency: 1})
	idx.Add("ant", "doc2", Posting{Fields: BODY, Frequency: 2})
	idx.Add("banana", "doc1", Posting{Fields: BODY, Frequency: 1})
	idx.Add("cherry", "doc3", Posting{Fields: TITLE | BODY, Frequency: 3})

	writer := NewIndexWriter(indexPath)
	err := writer.WriteIndex(idx)
	require.NoError(t, err)

	// Verify files a-z are created
	for char := 'a'; char <= 'z'; char++ {
		filename := filepath.Join(indexPath, fmt.Sprintf("index%c.idx", char))
		assert.FileExists(t, filename)
	}

	// Verify content of indexa.idx (should have ant and apple, sorted)
	contentA := readIndexFile(t, filepath.Join(indexPath, "indexa.idx"))
	require.Len(t, contentA, 2)
	assert.True(t, strings.HasPrefix(contentA[0], "ant:"))
	assert.True(t, strings.HasPrefix(contentA[1], "apple:"))

	// Verify specific format of an entry
	// Format: term:docID$fields$freq
	// ant:doc2$8$2
	assert.Equal(t, "ant:doc2$8$2", contentA[0])
	assert.Equal(t, "apple:doc1$32$1", contentA[1])

	// Verify indexb.idx
	contentB := readIndexFile(t, filepath.Join(indexPath, "indexb.idx"))
	require.Len(t, contentB, 1)
	assert.Equal(t, "banana:doc1$8$1", contentB[0])
}

func readIndexFile(t *testing.T, path string) []string {
	file, err := os.Open(path)
	require.NoError(t, err)
	defer func() { _ = file.Close() }()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
