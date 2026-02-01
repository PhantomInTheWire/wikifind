# WikiFind

WikiFind is a command-line search engine for Wikipedia XML dumps. It allows you to index large Wikipedia datasets and perform fast searches on them.

## Prerequisites

- Go 1.24 or later
- A Wikipedia XML dump file (e.g., from https://dumps.wikimedia.org/)

## Installation

Build the project:

```bash
go build -o wikifind ./cmd
```

## Usage

### Indexing

To index a Wikipedia XML dump:

```bash
./wikifind index <xml_file> <index_path>
```

- `<xml_file>`: Path to the Wikipedia XML dump file
- `<index_path>`: Directory where the index will be stored

Example:

```bash
./wikifind index enwiki-20231201-pages-articles.xml index/
```

### Searching

To search the indexed data:

```bash
./wikifind search <index_path>
```

- `<index_path>`: Directory containing the index

This will start an interactive search prompt. Enter your queries and get results.

Example:

```bash
./wikifind search index/
> apple
Found 5 results:
1. DocID: Apple (Score: 0.95)
...
```

## Architecture

The project is organized into several packages:

- `cmd/`: Main application entry point
- `indexer/`: Indexing logic, including XML parsing, text processing, and inverted index creation
- `search/`: Search engine implementation with compression and query processing
