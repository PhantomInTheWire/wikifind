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
	os.MkdirAll(indexPath, 0755)

	// Create dummy index files for a-z
	for char := 'a'; char <= 'z'; char++ {
		filename := filepath.Join(indexPath, fmt.Sprintf("index%c.idx", char))
		file, err := os.Create(filename)
		if err != nil {
			t.Fatal(err)
		}
		if char == 't' {
			file.WriteString("test:doc1$8$1:doc2$32$2\n")
		}
		file.Close()
	}

	se := NewSearchEngine(indexPath)
	se.Initialize()
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
}
