# WikiFind

WikiFind is a command-line search engine for Wikipedia XML dumps. It allows you to index large Wikipedia datasets and perform fast searches on them.

## Features

- Index Wikipedia XML dumps
- Perform keyword searches
- Stemming for better search results
- Support for various Wikipedia markup elements (categories, infoboxes, links)
- Efficient compression for index storage

## Prerequisites

- Go 1.24 or later
- A Wikipedia XML dump file (e.g., from https://dumps.wikimedia.org/)

## Installation

Clone the repository:

```bash
git clone https://github.com/PhantomInTheWire/wikifind.git
cd wikifind
```

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

## Testing

Run the tests:

```bash
go test ./...
```

For test coverage:

```bash
go test -cover ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run the pre-commit hooks
6. Submit a pull request

## License

This project is licensed under the MIT License.
