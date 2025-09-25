package indexer

import "testing"

func TestInvertedIndex_Add(t *testing.T) {
	idx := NewInvertedIndex()

	termObj1 := TermObject{Fields: BODY, Frequency: 1}
	termObj2 := TermObject{Fields: BODY, Frequency: 2}

	idx.Add("test", "doc1", termObj1)
	idx.Add("test", "doc1", termObj2) // Should merge

	if postings, ok := idx.Index["test"]; ok {
		if obj, exists := postings["doc1"]; exists {
			if obj.Frequency != 3 {
				t.Errorf("Expected frequency 3, got %d", obj.Frequency)
			}
			if obj.Fields != BODY {
				t.Errorf("Expected fields %d, got %d", BODY, obj.Fields)
			}
		} else {
			t.Error("Expected doc1 in postings")
		}
	} else {
		t.Error("Expected 'test' term in index")
	}
}

func TestTermObject_String(t *testing.T) {
	obj := TermObject{Fields: 8, Frequency: 5}
	expected := "8$5"
	if obj.String() != expected {
		t.Errorf("String() = %q, want %q", obj.String(), expected)
	}
}
