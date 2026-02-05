package indexer

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func NewWikiXMLParser(indexPath string) *WikiXMLParser {
	return &WikiXMLParser{
		indexPath: indexPath,
		index:     NewInvertedIndex(),
	}
}

func (parser *WikiXMLParser) Parse(ctx context.Context, xmlPath string) error {
	file, err := os.Open(xmlPath)
	if err != nil {
		return NewIOError("open file", err)
	}
	defer func() { _ = file.Close() }()

	return parser.ParseReader(ctx, file)
}

func (parser *WikiXMLParser) ParseReader(ctx context.Context, r io.Reader) error {
	decoder := xml.NewDecoder(r)
	pageCount := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return NewInvalidXMLError(err)
		}

		if se, ok := token.(xml.StartElement); ok && se.Name.Local == "page" {
			var xmlPage xmlPage
			if err := decoder.DecodeElement(&xmlPage, &se); err != nil {
				continue
			}

			doc := &Document{
				ID:       xmlPage.ID,
				Title:    xmlPage.Title,
				Content:  xmlPage.Text,
				Metadata: make(map[string]string),
			}

			if err := parser.processDocument(ctx, doc); err != nil {
				return err
			}

			pageCount++
			if pageCount%1000 == 0 {
				fmt.Printf("Processed %d pages\n", pageCount)
			}
		}
	}

	writer := NewIndexWriter(parser.indexPath)
	return writer.WriteIndex(parser.index)
}

func (parser *WikiXMLParser) processDocument(ctx context.Context, doc *Document) error {
	textParser := NewWikiTextParser(doc)
	terms := textParser.Parse()

	for term, posting := range terms {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		parser.index.Add(term, doc.ID, posting)
	}
	return nil
}
