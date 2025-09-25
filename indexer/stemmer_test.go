package indexer

import "testing"

func TestStemmer_Stem(t *testing.T) {
	s := NewStemmer()

	tests := []struct {
		input    string
		expected string
	}{
		{"running", "run"},
		{"cats", "cat"},
		{"jumped", "jump"},
		{"beautiful", "beauti"},
		{"", ""},
		{"a", "a"},
	}

	for _, test := range tests {
		result := s.Stem(test.input)
		if result != test.expected {
			t.Errorf("Stem(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}
