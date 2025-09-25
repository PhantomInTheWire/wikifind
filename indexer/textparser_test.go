package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsStopWord(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		expected bool
	}{
		{"the", "the", true},
		{"apple", "apple", false},
		{"a", "a", true},
		{"empty", "", true},
		{"run", "run", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsStopWord(tt.word)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWikiTextParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		page     *WikiPage
		expected map[string]bool // terms that should be present
	}{
		{
			name: "basic page",
			page: &WikiPage{
				ID:    "1",
				Title: "Test Page",
				Text:  "This is a test page with some content.",
			},
			expected: map[string]bool{
				"test": true,
				"page": true,
			},
		},
		{
			name: "page with markup",
			page: &WikiPage{
				ID:      "2",
				Title:   "Apple",
				Text:    "An apple is a fruit. [[Category:Fruits]] {{Infobox fruit|color=red|type=edible}} [[Link to something]].",
				Infobox: make(map[string]string),
			},
			expected: map[string]bool{
				"appl":     true, // stemmed
				"fruit":    true,
				"categori": true,
				"link":     true,
				"someth":   true,
				"color":    true,
				"red":      true,
				"type":     true,
				"edibl":    true,
			},
		},
		{
			name: "page with infobox no equals",
			page: &WikiPage{
				ID:      "3",
				Title:   "Test",
				Text:    "{{infobox test|key1|key2=value2}}",
				Infobox: make(map[string]string),
			},
			expected: map[string]bool{
				"test": true,
				"kei":  true, // key2 stemmed
				"valu": true, // value2 stemmed
			},
		},
		{
			name: "page with wiki markup to remove",
			page: &WikiPage{
				ID:    "3",
				Title: "Test",
				Text:  "This is <!-- comment --> text with <ref>reference</ref> and {{template}} and <b>bold</b>.",
			},
			expected: map[string]bool{
				"text": true,
				"bold": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewWikiTextParser(tt.page)
			terms := parser.Parse()

			assert.NotEmpty(t, terms, "Expected some terms")

			for term := range tt.expected {
				assert.Contains(t, terms, term, "Expected term %q to be present", term)
			}

			// Check infobox if present
			if tt.page.Infobox != nil {
				if tt.name == "page with markup" {
					assert.Equal(t, "red", tt.page.Infobox["color"])
					assert.Equal(t, "edible", tt.page.Infobox["type"])
				}
				if tt.name == "page with infobox no equals" {
					assert.Equal(t, "value2", tt.page.Infobox["key2"])
				}
			}
		})
	}
}
