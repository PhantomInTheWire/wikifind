package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStemmer_Stem(t *testing.T) {
	s := NewStemmer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"running", "running", "run"},
		{"cats", "cats", "cat"},
		{"jumped", "jumped", "jump"},
		{"beautiful", "beautiful", "beauti"},
		{"empty", "", ""},
		{"single", "a", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.Stem(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
