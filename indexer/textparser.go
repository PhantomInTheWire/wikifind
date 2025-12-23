package indexer

import (
	"regexp"
	"strings"
)

// Stop words - using empty struct for memory-efficient set
var stopWords = map[string]struct{}{
	"a": {}, "an": {}, "and": {}, "are": {}, "as": {}, "at": {},
	"be": {}, "by": {}, "for": {}, "from": {}, "has": {}, "he": {},
	"in": {}, "is": {}, "it": {}, "its": {}, "of": {}, "on": {},
	"that": {}, "the": {}, "to": {}, "was": {}, "were": {}, "will": {},
	"with": {}, "would": {}, "you": {}, "your": {},
}

func IsStopWord(word string) bool {
	if len(word) <= 1 {
		return true
	}
	_, exists := stopWords[strings.ToLower(word)]
	return exists
}

type WikiTextParser struct {
	stemmer *Stemmer
	doc     *Document
	terms   map[string]Posting
}

func NewWikiTextParser(doc *Document) *WikiTextParser {
	return &WikiTextParser{
		stemmer: NewStemmer(),
		doc:     doc,
		terms:   make(map[string]Posting),
	}
}

func (p *WikiTextParser) Parse() map[string]Posting {
	defer p.stemmer.Release() // Return stemmer to pool

	p.parseText(p.doc.Title, TITLE)

	p.parseWikiText(p.doc.Content)

	return p.terms
}

func (p *WikiTextParser) parseWikiText(text string) {
	text = strings.ToLower(text)

	categoryRegex := regexp.MustCompile(`\[\[category:([^\]]+)\]\]`)
	categories := categoryRegex.FindAllStringSubmatch(text, -1)
	for _, match := range categories {
		if len(match) > 1 {
			p.parseText(match[1], CATEGORY)
		}
	}

	infoboxRegex := regexp.MustCompile(`\{\{infobox([^}]*)\}\}`)
	infoboxes := infoboxRegex.FindAllStringSubmatch(text, -1)
	for _, match := range infoboxes {
		if len(match) > 1 {
			p.parseInfobox(match[1])
		}
	}

	geoboxRegex := regexp.MustCompile(`\{\{geobox[^}]*\}\}`)
	geoboxes := geoboxRegex.FindAllString(text, -1)
	for _, geobox := range geoboxes {
		p.parseText(geobox, GEOBOX)
	}

	linkRegex := regexp.MustCompile(`\[\[([^\]|]+)`)
	links := linkRegex.FindAllStringSubmatch(text, -1)
	for _, match := range links {
		if len(match) > 1 {
			p.parseText(match[1], LINKS)
		}
	}

	text = p.removeWikiMarkup(text)

	text = categoryRegex.ReplaceAllString(text, "")
	text = infoboxRegex.ReplaceAllString(text, "")
	text = geoboxRegex.ReplaceAllString(text, "")
	text = linkRegex.ReplaceAllString(text, "")

	p.parseText(text, BODY)
}

func (p *WikiTextParser) removeWikiMarkup(text string) string {
	commentRegex := regexp.MustCompile(`<!--.*?-->`)
	text = commentRegex.ReplaceAllString(text, "")

	refRegex := regexp.MustCompile(`<ref[^>]*>.*?</ref>`)
	text = refRegex.ReplaceAllString(text, "")

	templateRegex := regexp.MustCompile(`\{\{[^}]*\}\}`)
	text = templateRegex.ReplaceAllString(text, "")

	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	text = htmlRegex.ReplaceAllString(text, "")

	return text
}

func (p *WikiTextParser) parseInfobox(infoboxText string) {
	parts := strings.Split(infoboxText, "|")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				if key != "" && value != "" {
					p.doc.Metadata[strings.ToLower(key)] = strings.ToLower(value)
					p.parseText(key, INFOBOX)
					p.parseText(value, INFOBOX)
				}
			}
		}
	}
}

func (p *WikiTextParser) parseText(text string, field FieldMask) {
	wordRegex := regexp.MustCompile(`[a-z]+`)
	words := wordRegex.FindAllString(strings.ToLower(text), -1)

	for _, word := range words {
		if len(word) > 1 && !IsStopWord(word) {
			stemmed := p.stemmer.Stem(word)

			term := p.terms[stemmed]
			term.Fields |= field
			term.Frequency++
			p.terms[stemmed] = term
		}
	}
}
