package indexer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// Interfaces for extensibility and testing
type Parser interface {
	Parse(ctx context.Context, r io.Reader) ([]Document, error)
}

type Indexer interface {
	Index(ctx context.Context, docs []Document) error
	Close() error
}

type TextProcessor interface {
	Process(ctx context.Context, doc Document) (map[string]Posting, error)
}

type Document struct {
	ID       string
	Title    string
	Content  string
	Metadata map[string]string
}

type Posting struct {
	Fields    FieldMask
	Frequency int
}

type xmlPage struct {
	ID    string `xml:"id"`
	Title string `xml:"title"`
	Text  string `xml:"revision>text"`
}

// Field mask type for better type safety
type FieldMask byte

type WikiXMLParser struct {
	indexPath string
	index     *InvertedIndex
}

func (p Posting) String() string {
	return fmt.Sprintf("%d$%d", p.Fields, p.Frequency)
}

type InvertedIndex struct {
	Index map[string]map[string]Posting
	mutex sync.RWMutex
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Index: make(map[string]map[string]Posting),
	}
}

func (idx *InvertedIndex) Add(term, docID string, posting Posting) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	if idx.Index[term] == nil {
		idx.Index[term] = make(map[string]Posting)
	}

	existing := idx.Index[term][docID]
	existing.Fields |= posting.Fields
	existing.Frequency += posting.Frequency
	idx.Index[term][docID] = existing
}

// IndexWriter handles writing the inverted index to disk
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

func NewWikiXMLParser(indexPath string) *WikiXMLParser {
	return &WikiXMLParser{
		indexPath: indexPath,
		index:     NewInvertedIndex(),
	}
}

// Field constants (bit flags) - matching Java order
const (
	GEOBOX   FieldMask = 1 << 0 // 1
	CATEGORY FieldMask = 1 << 4 // 16
	TITLE    FieldMask = 1 << 5 // 32
	BODY     FieldMask = 1 << 3 // 8
	LINKS    FieldMask = 1 << 2 // 4
	INFOBOX  FieldMask = 1 << 1 // 2
)
