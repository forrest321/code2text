
# code2text

## Introduction
`code2text` is a CLI tool designed to walk through directories, collect code files based on specified filters, and consolidate their contents into a single text document. This tool is particularly useful for gathering code for analysis or archiving.

## Installation

### Prerequisites
- Go (1.15 or later)

### Install with Go  
   ```bash
   go install github.com/forrest321/code2text
   ```


### Build from Source
1. Clone the repository:
   ```bash
   git clone https://github.com/forrest321/code2text
   ```
2. Navigate to the project directory:
   ```bash
   cd code2text
   ```
3. Build the tool:
   ```bash
   go build -o code2text
   ```
4. Move the created binary to a location within your PATH
   ```bash
   sudo mv code2text /usr/bin/
   ```

## Usage

Run the tool using the following syntax:

```bash
./code2text [flags]
```

### Flags
- `-o, --output <file>`: Specify the output file name (default is `output.txt`).
- `-i, --include <extensions>`: Specify file extensions to include (default includes `.go`, `.js`).
- `-e, --exclude <extensions>`: Specify file extensions to exclude.
- `-g, --ignore <directories>`: Specify directories to ignore (default ignores `.git`, `.idea`).

### Example
```bash
./code2text --output result.txt --include .go --include .js --exclude .test --ignore .git
```

## Output Details

The output file begins with the directory structure, followed by the content of each file included in the scan. Each file's content is clearly separated and includes metadata about the file path and the number of lines.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
