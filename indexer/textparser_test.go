package indexer

import "testing"

func TestIsStopWord(t *testing.T) {
	tests := []struct {
		word     string
		expected bool
	}{
		{"the", true},
		{"apple", false},
		{"a", true},
		{"", true},
		{"run", false},
	}

	for _, test := range tests {
		result := IsStopWord(test.word)
		if result != test.expected {
			t.Errorf("IsStopWord(%q) = %v, want %v", test.word, result, test.expected)
		}
	}
}

func TestWikiTextParser_Parse(t *testing.T) {
	page := &WikiPage{
		ID:    "1",
		Title: "Test Page",
		Text:  "This is a test page with some content.",
	}

	parser := NewWikiTextParser(page)
	terms := parser.Parse()

	if len(terms) == 0 {
		t.Error("Expected some terms, got none")
	}

	// Check if title term is present
	if _, ok := terms["test"]; !ok {
		t.Error("Expected 'test' term from title")
	}
}
