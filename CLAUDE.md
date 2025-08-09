# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project implementing a lambda calculus parser and evaluator. The codebase consists of:

- **Module**: `github.com/gusbicalho/go-lambda`
- **Go Version**: 1.24 (managed via mise tool configuration)
- **Architecture**: Multi-package modular structure with complete parsing pipeline

## Code Architecture

### Package Structure

- `main` package: Entry point that orchestrates tokenization and parsing
- `parser` package: Complete recursive descent parser implementation with error handling
- `parse_tree` package: AST definitions with sealed union types and pretty printing
- `tokenizer` package: Lexical analysis with peek/consume interface
- `token` package: Token types and constructors for lambda calculus syntax
- `position` package: Source location tracking (line/column positions)
- `runes_reader` package: Unicode-aware input reader with position tracking
- `pretty` package: Pretty printing system for formatted output

### Parser Architecture

The parser is implemented as a recursive descent parser with these key components:

1. **Tokenization Pipeline**: `RunesReader` → `Tokenizer` → `Token` stream
2. **Parse Tree Construction**: Tokens are parsed into a typed AST
3. **Pretty Printing**: AST nodes implement `Pretty` interface for formatted output

### Core Types

#### Parse Tree (`parse_tree` package)
- `ParseTree`: Root structure containing `InputLocation` and `ParseItem`
- `ParseItem`: Sealed interface with four implementations:
  - `Var`: Variables with name
  - `Lambda`: Lambda expressions with argument name and body (`\x.body`)
  - `App`: Function applications with callee and arguments
  - `Parens`: Parenthesized expressions for grouping
- All types implement `Pretty` interface for consistent output formatting

#### Tokenizer (`tokenizer` package)
- `Tokenizer`: Main lexer with buffered peek/consume interface
- Supports lambda calculus syntax: `\`, `.`, `(`, `)`, identifiers
- Handles whitespace skipping and identifier recognition
- Built on `RunesReader` for Unicode-aware parsing with position tracking

#### Parser Implementation (`parser` package)
- `Parse(tokenizer)` function: Main entry point that ensures complete input consumption
- `ParseResult[T]` monad pattern for error handling:
  - Tracks whether input was consumed (for better error recovery)
  - Implements left-associative function application parsing
  - Handles precedence through separate parsing functions for different syntactic categories
- Uses pointer semantics (`*parse_tree.ParseTree`) for efficient memory management
- Validates complete input consumption and rejects leftover tokens

## Common Commands

### Building
```bash
go build
```

### Running
```bash
go run .
```
The main program reads lambda calculus expressions from stdin and outputs the parsed AST with pretty formatting.

### Testing
```bash
go test ./...
```
Note: Currently no test files exist in the project.

### Module Management
```bash
go mod tidy
go mod download
```

### Development Environment
This project uses mise for Go version management. The Go 1.24 toolchain is specified in `mise.toml`.