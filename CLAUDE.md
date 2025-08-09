# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project implementing a lambda calculus parser and evaluator. The codebase consists of:

- **Module**: `github.com/gusbicalho/go-lambda`
- **Go Version**: 1.24 (managed via mise tool configuration)
- **Architecture**: Multi-package modular structure with complete parsing pipeline

## Code Architecture

### Package Structure

- `main` package: Entry point that orchestrates parsing and locally nameless conversion
- `parser` package: Complete recursive descent parser implementation with error handling
- `parse_tree` package: AST definitions with sealed union types and pretty printing
- `parse_tree_to_locally_nameless` package: Conversion logic from named to De Bruijn representation
- `locally_nameless` package: De Bruijn index representation for lambda expressions
- `stack` package: Functional stack data structure with iterator support
- `tokenizer` package: Lexical analysis with peek/consume interface
- `token` package: Token types and constructors for lambda calculus syntax
- `position` package: Source location tracking (line/column positions)
- `runes_reader` package: Unicode-aware input reader with position tracking
- `pretty` package: Context-aware pretty printing system for formatted output

### Parser Architecture

The system implements a complete lambda calculus interpreter with these key components:

1. **Tokenization Pipeline**: `RunesReader` → `Tokenizer` → `Token` stream
2. **Parse Tree Construction**: Tokens are parsed into a named AST representation
3. **Locally Nameless Conversion**: Named variables converted to De Bruijn indices
4. **Context-Aware Pretty Printing**: AST nodes implement `Pretty[context]` interface

### Core Types

#### Parse Tree (`parse_tree` package)
- `ParseTree`: Root structure containing `InputLocation` and `ParseItem`
- `ParseItem`: Sealed interface with four implementations:
  - `Var`: Variables with name
  - `Lambda`: Lambda expressions with argument name and body (`\x.body`)
  - `App`: Function applications with callee and arguments
  - `Parens`: Parenthesized expressions for grouping
- All types implement `Pretty[any]` interface for context-aware formatting

#### Locally Nameless (`locally_nameless` package)
Implementation of the locally nameless representation using De Bruijn indices:
- `Expr`: Sealed interface for lambda expressions without name capture issues
- `FreeVar`: Free variables represented by their original names
- `BoundVar`: Bound variables represented by De Bruijn indices (distance to binder)
- `Lambda`: Lambda abstraction with original parameter name for display
- `App`: Binary function application
- All types implement `Pretty[Stack[string]]` for context-aware display with bound variable names

#### Stack Data Structure (`stack` package)
- `Stack[T]`: Immutable functional stack with generic type parameter
- `Push(value)`: Returns new stack with value on top
- `Pop()`: Returns `*Popped[T]` with value and remaining stack
- `Nth(index, default)`: Access element by index with fallback
- `Items()` and `IndexedItems()`: Iterator support using Go 1.23 `iter` package

#### Pretty Printing (`pretty` package)
- `Pretty[context]`: Generic interface supporting context-dependent formatting
- `PrettyDoc`: Document structure for consistent indentation and layout
- Context parameter allows passing state (e.g., bound variable names) through printing
- `String()` method provides final string rendering

#### Parser Implementation (`parser` package)
- `Parse(tokenizer)` function: Main entry point that ensures complete input consumption
- `ParseResult[T]` monad pattern for error handling:
  - Tracks whether input was consumed (for better error recovery)
  - Implements left-associative function application parsing
  - Handles precedence through separate parsing functions for different syntactic categories
- Uses pointer semantics (`*parse_tree.ParseTree`) for efficient memory management
- Validates complete input consumption and rejects leftover tokens

#### Variable Binding Conversion (`parse_tree_to_locally_nameless` package)
Dedicated package for converting named AST to De Bruijn indexed representation:
- `ToLocallyNameless()`: Main entry point for conversion
- `toLocallyNameless()`: Internal recursive conversion with bound variable tracking
- Uses functional stack to track bound variable names during traversal
- Free variables remain as names, bound variables become indices
- Handles all parse tree node types: `Parens`, `Var`, `Lambda`, `App`

## Common Commands

### Building
```bash
go build
```

### Running
```bash
go run .
```
The main program reads lambda calculus expressions from stdin and outputs both:
1. The original parse tree with named variables
2. The locally nameless representation with De Bruijn indices

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