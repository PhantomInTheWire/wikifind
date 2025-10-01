package indexer

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func (parser *WikiXMLParser) Parse(xmlPath string) error {
	file, err := os.Open(xmlPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	decoder := xml.NewDecoder(file)
	pageCount := 0

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if se, ok := token.(xml.StartElement); ok && se.Name.Local == "page" {
			var xmlPage XMLPage
			if err := decoder.DecodeElement(&xmlPage, &se); err != nil {
				continue
			}

			page := &WikiPage{
				ID:      xmlPage.ID,
				Title:   xmlPage.Title,
				Text:    xmlPage.Text,
				Infobox: make(map[string]string),
			}

			parser.processPage(page)
			pageCount++

			if pageCount%1000 == 0 {
				fmt.Printf("Processed %d pages\n", pageCount)
			}
		}
	}

	return parser.writeIndex()
}

func (parser *WikiXMLParser) processPage(page *WikiPage) {
	textParser := NewWikiTextParser(page)
	terms := textParser.Parse()

	for term, termObj := range terms {
		parser.index.Add(term, page.ID, termObj)
	}
}
