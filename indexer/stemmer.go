package indexer

import "strings"

// Stemmer - simplified Porter stemmer
type Stemmer struct {
	word []rune
}

func NewStemmer() *Stemmer {
	return &Stemmer{}
}

func (s *Stemmer) Stem(word string) string {
	s.word = []rune(strings.ToLower(word))
	if len(s.word) <= 2 {
		return string(s.word)
	}

	s.step1()
	s.step2()
	s.step3()

	return string(s.word)
}

func (s *Stemmer) step1() {
	// Remove plurals and -ed/-ing
	if s.endsWith("sses") {
		s.word = s.word[:len(s.word)-2]
	} else if s.endsWith("ies") {
		s.word = s.word[:len(s.word)-3]
		s.word = append(s.word, 'i')
	} else if s.endsWith("s") && len(s.word) > 1 && s.word[len(s.word)-2] != 's' {
		s.word = s.word[:len(s.word)-1]
	}

	if s.endsWith("eed") {
		if s.measure() > 0 {
			s.word = s.word[:len(s.word)-1]
		}
	} else if (s.endsWith("ed") || s.endsWith("ing")) && s.containsVowel() {
		if s.endsWith("ed") {
			s.word = s.word[:len(s.word)-2]
		} else {
			s.word = s.word[:len(s.word)-3]
		}
	}
}

func (s *Stemmer) step2() {
	if s.endsWith("y") && s.containsVowel() && len(s.word) > 1 {
		s.word[len(s.word)-1] = 'i'
	}
}

func (s *Stemmer) step3() {
	// Simplified step3 - handle common suffixes
	if s.endsWith("ational") && s.measure() > 0 {
		s.word = s.word[:len(s.word)-7]
		s.word = append(s.word, []rune("ate")...)
	} else if s.endsWith("tion") && s.measure() > 0 {
		s.word = s.word[:len(s.word)-4]
		s.word = append(s.word, []rune("tion")...)
	}
}

func (s *Stemmer) endsWith(suffix string) bool {
	suffixRunes := []rune(suffix)
	if len(s.word) < len(suffixRunes) {
		return false
	}

	for i := 0; i < len(suffixRunes); i++ {
		if s.word[len(s.word)-len(suffixRunes)+i] != suffixRunes[i] {
			return false
		}
	}
	return true
}

func (s *Stemmer) isVowel(i int) bool {
	if i < 0 || i >= len(s.word) {
		return false
	}
	c := s.word[i]
	return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' ||
		(c == 'y' && (i == 0 || !s.isVowel(i-1)))
}

func (s *Stemmer) containsVowel() bool {
	for i := 0; i < len(s.word); i++ {
		if s.isVowel(i) {
			return true
		}
	}
	return false
}

func (s *Stemmer) measure() int {
	n := 0
	i := 0
	for i < len(s.word) && !s.isVowel(i) {
		i++
	}
	for i < len(s.word) {
		for i < len(s.word) && s.isVowel(i) {
			i++
		}
		if i >= len(s.word) {
			break
		}
		n++
		for i < len(s.word) && !s.isVowel(i) {
			i++
		}
	}
	return n
}
