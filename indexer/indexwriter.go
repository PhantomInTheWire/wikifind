package indexer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func (parser *WikiXMLParser) writeIndex() error {
	// Create index directory
	if err := os.MkdirAll(parser.indexPath, 0755); err != nil {
		return err
	}

	// Write inverted index files (a-z)
	for char := 'a'; char <= 'z'; char++ {
		if err := parser.writeIndexFile(char); err != nil {
			return err
		}
	}

	return nil
}

func (parser *WikiXMLParser) writeIndexFile(char rune) error {
	filename := filepath.Join(parser.indexPath, fmt.Sprintf("index%c.idx", char))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Collect terms starting with this character
	var terms []string
	parser.index.mutex.RLock()
	for term := range parser.index.Index {
		if len(term) > 0 && rune(term[0]) == char {
			terms = append(terms, term)
		}
	}
	parser.index.mutex.RUnlock()

	sort.Strings(terms)

	// Write sorted terms and their posting lists
	for _, term := range terms {
		parser.index.mutex.RLock()
		postings := parser.index.Index[term]
		parser.index.mutex.RUnlock()

		var docIDs []string
		for docID := range postings {
			docIDs = append(docIDs, docID)
		}
		sort.Strings(docIDs)

		// Write term and postings
		fmt.Fprintf(writer, "%s", term)
		for _, docID := range docIDs {
			termObj := postings[docID]
			fmt.Fprintf(writer, ":%s$%s", docID, termObj.String())
		}
		fmt.Fprintln(writer)
	}

	return nil
}
