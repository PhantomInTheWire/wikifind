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

// Search engine
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
	// Open all index files
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

	// Get posting lists for each term
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

	// Sort results by score
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

	// Convert to search results
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

func editDistance(str1, str2 string) int {
	distance := make([][]int, len(str1)+1)
	for i := range distance {
		distance[i] = make([]int, len(str2)+1)
	}

	for i := 0; i <= len(str1); i++ {
		distance[i][0] = i
	}
	for j := 1; j <= len(str2); j++ {
		distance[0][j] = j
	}

	for i := 1; i <= len(str1); i++ {
		for j := 1; j <= len(str2); j++ {
			cost := 0
			if str1[i-1] != str2[j-1] {
				cost = 1
			}
			distance[i][j] = min(
				distance[i-1][j]+1,      // deletion
				distance[i][j-1]+1,      // insertion
				distance[i-1][j-1]+cost, // substitution
			)
		}
	}

	return distance[len(str1)][len(str2)]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

func (se *SearchEngine) parseQuery(query string) []string {
	stemmer := indexer.NewStemmer()
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

func (se *SearchEngine) getPostings(term string) (map[string]indexer.TermObject, error) {
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

	// Reset file pointer
	_, _ = file.Seek(0, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		if parts[0] == term {
			// Parse postings
			postings := make(map[string]indexer.TermObject)
			for i := 1; i < len(parts); i++ {
				posting := parts[i]
				dollarIdx := strings.Index(posting, "$")
				if dollarIdx == -1 {
					continue
				}

				docID := posting[:dollarIdx]
				termData := posting[dollarIdx+1:]

				// Parse term object
				termParts := strings.Split(termData, "$")
				if len(termParts) != 2 {
					continue
				}

				fields, _ := strconv.Atoi(termParts[0])
				freq, _ := strconv.Atoi(termParts[1])

				postings[docID] = indexer.TermObject{
					Fields:    byte(fields),
					Frequency: freq,
				}
			}
			return postings, nil
		}
	}

	return make(map[string]indexer.TermObject), nil
}
