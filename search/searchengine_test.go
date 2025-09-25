package search

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/PhantomInTheWire/wikifind/indexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchEngine_parseQuery(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []string
	}{
		{"single word", "hello", []string{"hello"}},
		{"two words", "hello world", []string{"hello", "world"}},
		{"empty", "", nil},
		{"with stop words", "the apple is red", []string{"appl", "red"}},
		{"short words", "a an i", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			se := &SearchEngine{}
			terms := se.parseQuery(tt.query)
			assert.Equal(t, tt.expected, terms)
		})
	}
}

func TestSearchEngine_Initialize(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(string)
		expectErr bool
	}{
		{
			name: "successful initialize",
			setupFunc: func(indexPath string) {
				require.NoError(t, os.MkdirAll(indexPath, 0755))
				for char := 'a'; char <= 'z'; char++ {
					filename := filepath.Join(indexPath, fmt.Sprintf("index%c.idx", char))
					require.NoError(t, os.WriteFile(filename, []byte{}, 0644))
				}
			},
			expectErr: false,
		},
		{
			name: "missing index files",
			setupFunc: func(indexPath string) {
				require.NoError(t, os.MkdirAll(indexPath, 0755))
				// Don't create files
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			indexPath := filepath.Join(tempDir, "index")
			tt.setupFunc(indexPath)

			se := NewSearchEngine(indexPath)
			err := se.Initialize()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				se.Close()
			}
		})
	}
}

func TestSearchEngine_getPostings(t *testing.T) {
	// Create temp index file
	tempDir := t.TempDir()
	indexPath := filepath.Join(tempDir, "index")

	// Create index directory
	require.NoError(t, os.MkdirAll(indexPath, 0755))

	// Create dummy index files for a-z
	for char := 'a'; char <= 'z'; char++ {
		filename := filepath.Join(indexPath, fmt.Sprintf("index%c.idx", char))
		file, err := os.Create(filename)
		require.NoError(t, err)
		if char == 't' {
			_, err = file.WriteString("test:doc1$8$1:doc2$32$2\n")
			require.NoError(t, err)
		}
		require.NoError(t, file.Close())
	}

	se := NewSearchEngine(indexPath)
	require.NoError(t, se.Initialize())
	defer se.Close()

	// Test successful getPostings
	postings, err := se.getPostings("test")
	require.NoError(t, err)

	expectedPostings := map[string]indexer.TermObject{
		"doc1": {Fields: 8, Frequency: 1},
		"doc2": {Fields: 32, Frequency: 2},
	}
	assert.Equal(t, expectedPostings, postings)

	// Test empty term
	_, err = se.getPostings("")
	assert.Error(t, err)

	// Test invalid char
	_, err = se.getPostings("1test")
	assert.Error(t, err)

	// Test Search
	results, err := se.Search("test", 10)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Test Search with no results
	results2, err := se.Search("nonexistent", 10)
	require.NoError(t, err)
	assert.Empty(t, results2)

	// Test Search with no valid terms
	results3, err := se.Search("the a an", 10)
	assert.Error(t, err)
	assert.Nil(t, results3)
}

func TestEditDistance(t *testing.T) {
	tests := []struct {
		name     string
		str1     string
		str2     string
		expected int
	}{
		{"identical", "kitten", "kitten", 0},
		{"different", "kitten", "sitting", 3},
		{"empty", "", "", 0},
		{"single diff", "a", "b", 1},
		{"all diff", "abc", "def", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := editDistance(tt.str1, tt.str2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a, b, c  int
		expected int
	}{
		{"a min", 1, 2, 3, 1},
		{"b min", 2, 1, 3, 1},
		{"c min", 3, 2, 1, 1},
		{"equal", 1, 1, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := min(tt.a, tt.b, tt.c)
			assert.Equal(t, tt.expected, result)
		})
	}
}
