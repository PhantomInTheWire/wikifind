package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PhantomInTheWire/wikifind/indexer"
	"github.com/PhantomInTheWire/wikifind/search"
)

// Main application
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: wikifind <command> <args>")
		fmt.Println("Commands:")
		fmt.Println("  index <xml_file> <index_path>")
		fmt.Println("  search <index_path>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "index":
		if len(os.Args) != 4 {
			fmt.Println("Usage: wikifind index <xml_file> <index_path>")
			os.Exit(1)
		}

		xmlFile := os.Args[2]
		indexPath := os.Args[3]

		fmt.Printf("Parsing Wikipedia XML dump: %s\n", xmlFile)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle Ctrl+C gracefully
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigChan
			fmt.Println("\nReceived interrupt signal, cancelling...")
			cancel()
		}()

		parser := indexer.NewWikiXMLParser(indexPath)

		if err := parser.Parse(ctx, xmlFile); err != nil {
			log.Fatalf("Error parsing XML: %v", err)
		}

		fmt.Println("Indexing completed successfully!")

	case "search":
		if len(os.Args) != 3 {
			fmt.Println("Usage: wikifind search <index_path>")
			os.Exit(1)
		}

		indexPath := os.Args[2]

		fmt.Println("Initializing search engine...")
		engine := search.NewSearchEngine(indexPath)

		if err := engine.Initialize(); err != nil {
			log.Fatalf("Error initializing search engine: %v", err)
		}
		defer engine.Close()

		fmt.Println("Search engine ready. Enter queries (Ctrl+C to exit):")

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}

			query := scanner.Text()
			if query == "" {
				continue
			}

			results, err := engine.Search(query, 10)
			if err != nil {
				fmt.Printf("Search error: %v\n", err)
				continue
			}

			if len(results) == 0 {
				fmt.Println("No results found.")
				continue
			}

			fmt.Printf("Found %d results:\n", len(results))
			for i, result := range results {
				fmt.Printf("%d. DocID: %s (Score: %.4f)\n", i+1, result.DocID, result.Score)
			}
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
