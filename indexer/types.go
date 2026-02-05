package indexer

import (
	"context"
	"io"
)

type Parser interface {
	Parse(ctx context.Context, r io.Reader) ([]Document, error)
}

type Indexer interface {
	Index(ctx context.Context, docs []Document) error
	Close() error
}

type TextProcessor interface {
	Process(ctx context.Context, doc Document) (map[string]Posting, error)
}

type Document struct {
	ID       string
	Title    string
	Content  string
	Metadata map[string]string
}

type Posting struct {
	Fields    FieldMask
	Frequency int
}

type xmlPage struct {
	ID    string `xml:"id"`
	Title string `xml:"title"`
	Text  string `xml:"revision>text"`
}

type FieldMask byte

type WikiXMLParser struct {
	indexPath string
	index     *InvertedIndex
}

const (
	GEOBOX   FieldMask = 1 << 0 // 1
	CATEGORY FieldMask = 1 << 4 // 16
	TITLE    FieldMask = 1 << 5 // 32
	BODY     FieldMask = 1 << 3 // 8
	LINKS    FieldMask = 1 << 2 // 4
	INFOBOX  FieldMask = 1 << 1 // 2
)
