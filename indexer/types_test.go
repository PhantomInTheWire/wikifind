package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvertedIndex_Add(t *testing.T) {
	tests := []struct {
		name     string
		term     string
		docID    string
		postings []Posting
		expected map[string]map[string]Posting
	}{
		{
			name:  "add single term",
			term:  "test",
			docID: "doc1",
			postings: []Posting{
				{Fields: BODY, Frequency: 1},
			},
			expected: map[string]map[string]Posting{
				"test": {
					"doc1": {Fields: BODY, Frequency: 1},
				},
			},
		},
		{
			name:  "merge terms",
			term:  "test",
			docID: "doc1",
			postings: []Posting{
				{Fields: BODY, Frequency: 1},
				{Fields: BODY, Frequency: 2},
			},
			expected: map[string]map[string]Posting{
				"test": {
					"doc1": {Fields: BODY, Frequency: 3},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewInvertedIndex()
			for _, posting := range tt.postings {
				idx.Add(tt.term, tt.docID, posting)
			}
			assert.Equal(t, tt.expected, idx.Index)
		})
	}
}

func TestPosting_String(t *testing.T) {
	tests := []struct {
		name     string
		posting  Posting
		expected string
	}{
		{"basic", Posting{Fields: BODY, Frequency: 5}, "8$5"},
		{"zero", Posting{Fields: 0, Frequency: 0}, "0$0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.posting.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}
