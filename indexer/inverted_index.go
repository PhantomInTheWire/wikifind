package indexer

import (
	"fmt"
	"sync"
)

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

func (p Posting) String() string {
	return fmt.Sprintf("%d$%d", p.Fields, p.Frequency)
}
