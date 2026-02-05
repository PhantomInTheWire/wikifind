package search

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/PhantomInTheWire/wikifind/indexer"
)

type SearchEngine struct {
	indexPath string
	indexes   map[rune]*os.File
	mutex     sync.RWMutex
}

func NewSearchEngine(indexPath string) *SearchEngine {
	return &SearchEngine{
		indexPath: indexPath,
		indexes:   make(map[rune]*os.File),
	}
}

func (se *SearchEngine) Initialize() error {
	for char := 'a'; char <= 'z'; char++ {
		filename := filepath.Join(se.indexPath, fmt.Sprintf("index%c.idx", char))
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		se.indexes[char] = file
	}
	return nil
}

func (se *SearchEngine) Close() {
	for _, file := range se.indexes {
		_ = file.Close()
	}
}

func (se *SearchEngine) Search(query string, limit int) ([]SearchResult, error) {
	terms := se.parseQuery(query)
	if len(terms) == 0 {
		return nil, fmt.Errorf("no valid terms in query")
	}

	docScores := make(map[string]float64)

	for _, term := range terms {
		postings, err := se.getPostings(term)
		if err != nil {
			continue
		}

		idf := math.Log10(14128976.0 / float64(len(postings)))

		for docID, termObj := range postings {
			tf := 1.0 + math.Log10(float64(termObj.Frequency))
			score := tf * idf

			// Boost title matches
			if termObj.Fields&indexer.TITLE != 0 {
				score *= 2.0
			}

			docScores[docID] += score
		}
	}

	type docScore struct {
		docID string
		score float64
	}

	var results []docScore
	for docID, score := range docScores {
		results = append(results, docScore{docID, score})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	var searchResults []SearchResult
	for i, result := range results {
		if i >= limit {
			break
		}
		searchResults = append(searchResults, SearchResult{
			DocID: result.docID,
			Score: result.score,
		})
	}

	return searchResults, nil
}

func (se *SearchEngine) parseQuery(query string) []string {
	stemmer := indexer.NewStemmer()
	defer stemmer.Release()

	wordRegex := regexp.MustCompile(`[a-z]+`)
	words := wordRegex.FindAllString(strings.ToLower(query), -1)

	var terms []string
	for _, word := range words {
		if len(word) > 1 && !indexer.IsStopWord(word) {
			terms = append(terms, stemmer.Stem(word))
		}
	}

	return terms
}

func (se *SearchEngine) getPostings(term string) (map[string]indexer.Posting, error) {
	if len(term) == 0 {
		return nil, fmt.Errorf("empty term")
	}

	char := rune(term[0])
	if char < 'a' || char > 'z' {
		return nil, fmt.Errorf("invalid term")
	}

	se.mutex.RLock()
	file := se.indexes[char]
	se.mutex.RUnlock()

	if file == nil {
		return nil, fmt.Errorf("index file not found")
	}

	_, _ = file.Seek(0, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		if parts[0] == term {
			postings := make(map[string]indexer.Posting)
			for i := 1; i < len(parts); i++ {
				posting := parts[i]
				dollarIdx := strings.Index(posting, "$")
				if dollarIdx == -1 {
					continue
				}

				docID := posting[:dollarIdx]
				termData := posting[dollarIdx+1:]

				termParts := strings.Split(termData, "$")
				if len(termParts) != 2 {
					continue
				}

				fields, _ := strconv.Atoi(termParts[0])
				freq, _ := strconv.Atoi(termParts[1])

				postings[docID] = indexer.Posting{
					Fields:    indexer.FieldMask(fields),
					Frequency: freq,
				}
			}
			return postings, nil
		}
	}

	return make(map[string]indexer.Posting), nil
}
