package search

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestSearchEngine_parseQuery(t *testing.T) {
	se := &SearchEngine{}

	terms := se.parseQuery("hello world")

	if len(terms) != 2 {
		t.Errorf("Expected 2 terms, got %d", len(terms))
	}

	expected := []string{"hello", "world"}
	for i, term := range terms {
		if term != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], term)
		}
	}
}

func TestSearchEngine_getPostings(t *testing.T) {
	// Create temp index file
	tempDir := t.TempDir()
	indexPath := filepath.Join(tempDir, "index")

	// Create index directory
	_ = os.MkdirAll(indexPath, 0755)

	// Create dummy index files for a-z
	for char := 'a'; char <= 'z'; char++ {
		filename := filepath.Join(indexPath, fmt.Sprintf("index%c.idx", char))
		file, err := os.Create(filename)
		if err != nil {
			t.Fatal(err)
		}
		if char == 't' {
			_, _ = file.WriteString("test:doc1$8$1:doc2$32$2\n")
		}
		_ = file.Close()
	}

	se := NewSearchEngine(indexPath)
	_ = se.Initialize()
	defer se.Close()

	postings, err := se.getPostings("test")
	if err != nil {
		t.Fatalf("getPostings failed: %v", err)
	}

	if len(postings) != 2 {
		t.Errorf("Expected 2 postings, got %d", len(postings))
	}

	if obj, ok := postings["doc1"]; ok {
		if obj.Fields != 8 || obj.Frequency != 1 {
			t.Errorf("Wrong term object for doc1: %+v", obj)
		}
	} else {
		t.Error("Missing doc1")
	}

	// Test Search
	results, err := se.Search("test", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestEditDistance(t *testing.T) {
	tests := []struct {
		str1, str2 string
		expected   int
	}{
		{"kitten", "kitten", 0},
		{"kitten", "sitting", 3},
		{"", "", 0},
		{"a", "b", 1},
		{"abc", "def", 3},
	}

	for _, test := range tests {
		result := editDistance(test.str1, test.str2)
		if result != test.expected {
			t.Errorf("editDistance(%q, %q) = %d, want %d", test.str1, test.str2, result, test.expected)
		}
	}
}
