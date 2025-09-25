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
		termObjs []TermObject
		expected map[string]map[string]TermObject
	}{
		{
			name:  "add single term",
			term:  "test",
			docID: "doc1",
			termObjs: []TermObject{
				{Fields: BODY, Frequency: 1},
			},
			expected: map[string]map[string]TermObject{
				"test": {
					"doc1": {Fields: BODY, Frequency: 1},
				},
			},
		},
		{
			name:  "merge terms",
			term:  "test",
			docID: "doc1",
			termObjs: []TermObject{
				{Fields: BODY, Frequency: 1},
				{Fields: BODY, Frequency: 2},
			},
			expected: map[string]map[string]TermObject{
				"test": {
					"doc1": {Fields: BODY, Frequency: 3},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewInvertedIndex()
			for _, obj := range tt.termObjs {
				idx.Add(tt.term, tt.docID, obj)
			}
			assert.Equal(t, tt.expected, idx.Index)
		})
	}
}

func TestTermObject_String(t *testing.T) {
	tests := []struct {
		name     string
		obj      TermObject
		expected string
	}{
		{"basic", TermObject{Fields: 8, Frequency: 5}, "8$5"},
		{"zero", TermObject{Fields: 0, Frequency: 0}, "0$0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.obj.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}
