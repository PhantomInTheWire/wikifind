package indexer

import (
	"regexp"
	"strings"
)

// Stop words
var stopWords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
	"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
	"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
	"that": true, "the": true, "to": true, "was": true, "were": true, "will": true,
	"with": true, "would": true, "you": true, "your": true,
}

func IsStopWord(word string) bool {
	return len(word) <= 1 || stopWords[strings.ToLower(word)]
}

// Text parser
type WikiTextParser struct {
	stemmer *Stemmer
	page    *WikiPage
	terms   map[string]TermObject
}

func NewWikiTextParser(page *WikiPage) *WikiTextParser {
	return &WikiTextParser{
		stemmer: NewStemmer(),
		page:    page,
		terms:   make(map[string]TermObject),
	}
}

func (p *WikiTextParser) Parse() map[string]TermObject {
	// Parse title
	p.parseText(p.page.Title, TITLE)

	// Parse main text
	p.parseWikiText(p.page.Text)

	return p.terms
}

func (p *WikiTextParser) parseWikiText(text string) {
	text = strings.ToLower(text)

	// Extract categories
	categoryRegex := regexp.MustCompile(`\[\[category:([^\]]+)\]\]`)
	categories := categoryRegex.FindAllStringSubmatch(text, -1)
	for _, match := range categories {
		if len(match) > 1 {
			p.parseText(match[1], CATEGORY)
		}
	}

	// Extract infobox
	infoboxRegex := regexp.MustCompile(`\{\{infobox([^}]*)\}\}`)
	infoboxes := infoboxRegex.FindAllStringSubmatch(text, -1)
	for _, match := range infoboxes {
		if len(match) > 1 {
			p.parseInfobox(match[1])
		}
	}

	// Extract geobox
	geoboxRegex := regexp.MustCompile(`\{\{geobox[^}]*\}\}`)
	geoboxes := geoboxRegex.FindAllString(text, -1)
	for _, geobox := range geoboxes {
		p.parseText(geobox, GEOBOX)
	}

	// Extract links
	linkRegex := regexp.MustCompile(`\[\[([^\]|]+)`)
	links := linkRegex.FindAllStringSubmatch(text, -1)
	for _, match := range links {
		if len(match) > 1 {
			p.parseText(match[1], LINKS)
		}
	}

	// Remove wiki markup
	text = p.removeWikiMarkup(text)

	// Parse remaining body text
	text = categoryRegex.ReplaceAllString(text, "")
	text = infoboxRegex.ReplaceAllString(text, "")
	text = geoboxRegex.ReplaceAllString(text, "")
	text = linkRegex.ReplaceAllString(text, "")

	p.parseText(text, BODY)
}

func (p *WikiTextParser) removeWikiMarkup(text string) string {
	// Remove comments
	commentRegex := regexp.MustCompile(`<!--.*?-->`)
	text = commentRegex.ReplaceAllString(text, "")

	// Remove references
	refRegex := regexp.MustCompile(`<ref[^>]*>.*?</ref>`)
	text = refRegex.ReplaceAllString(text, "")

	// Remove templates
	templateRegex := regexp.MustCompile(`\{\{[^}]*\}\}`)
	text = templateRegex.ReplaceAllString(text, "")

	// Remove HTML tags
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	text = htmlRegex.ReplaceAllString(text, "")

	return text
}

func (p *WikiTextParser) parseInfobox(infoboxText string) {
	// Parse key=value pairs
	parts := strings.Split(infoboxText, "|")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				if key != "" && value != "" {
					p.page.Infobox[strings.ToLower(key)] = strings.ToLower(value)
					// Index key and value
					p.parseText(key, INFOBOX)
					p.parseText(value, INFOBOX)
				}
			}
		}
	}
}

func (p *WikiTextParser) parseText(text string, field byte) {
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
