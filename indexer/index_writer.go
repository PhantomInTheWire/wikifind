package indexer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type IndexWriter struct {
	indexPath string
}

func NewIndexWriter(indexPath string) *IndexWriter {
	return &IndexWriter{indexPath: indexPath}
}

func (w *IndexWriter) WriteIndex(index *InvertedIndex) error {
	if err := os.MkdirAll(w.indexPath, 0755); err != nil {
		return err
	}

	for char := 'a'; char <= 'z'; char++ {
		if err := w.writeIndexFile(char, index); err != nil {
			return err
		}
	}

	return nil
}

func (w *IndexWriter) writeIndexFile(char rune, index *InvertedIndex) error {
	filename := filepath.Join(w.indexPath, fmt.Sprintf("index%c.idx", char))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	writer := bufio.NewWriter(file)
	defer func() { _ = writer.Flush() }()

	var terms []string
	index.mutex.RLock()
	for term := range index.Index {
		if len(term) > 0 && rune(term[0]) == char {
			terms = append(terms, term)
		}
	}
	index.mutex.RUnlock()

	sort.Strings(terms)

	for _, term := range terms {
		index.mutex.RLock()
		postings := index.Index[term]
		index.mutex.RUnlock()

		var docIDs []string
		for docID := range postings {
			docIDs = append(docIDs, docID)
		}
		sort.Strings(docIDs)

		_, _ = fmt.Fprintf(writer, "%s", term)
		for _, docID := range docIDs {
			posting := postings[docID]
			_, _ = fmt.Fprintf(writer, ":%s$%s", docID, posting.String())
		}
		_, _ = fmt.Fprintln(writer)
	}

	return nil
}
