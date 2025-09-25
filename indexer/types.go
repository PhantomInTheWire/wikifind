package indexer

import (
	"fmt"
	"sync"
)

// Core data structures
type WikiPage struct {
	ID    string
	Title string
	Text  string
}

type TermObject struct {
	Fields    byte
	Frequency int
}

func (t TermObject) String() string {
	return fmt.Sprintf("%d$%d", t.Fields, t.Frequency)
}

type InvertedIndex struct {
	Index map[string]map[string]TermObject
	mutex sync.RWMutex
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Index: make(map[string]map[string]TermObject),
	}
}

func (idx *InvertedIndex) Add(term, docID string, termObj TermObject) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	if idx.Index[term] == nil {
		idx.Index[term] = make(map[string]TermObject)
	}

	existing := idx.Index[term][docID]
	existing.Fields |= termObj.Fields
	existing.Frequency += termObj.Frequency
	idx.Index[term][docID] = existing
}

// Field constants (bit flags)
const (
	BODY     = 1 << 3 // 8
	LINKS    = 1 << 2 // 4
	INFOBOX  = 1 << 1 // 2
	GEOBOX   = 1 << 0 // 1
	CATEGORY = 1 << 4 // 16
	TITLE    = 1 << 5 // 32
)
