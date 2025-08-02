# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project implementing a lambda calculus parser and evaluator. The codebase consists of:

- **Module**: `github.com/gusbicalho/go-lambda`
- **Go Version**: 1.24.5 (managed via mise tool configuration)
- **Architecture**: Simple two-package structure with main entry point and parse tree definitions

## Code Architecture

### Package Structure

- `main` package: Entry point (currently empty main function)
- `parse_tree` package: Core data structures for lambda calculus AST

### Core Types

The `parse_tree` package defines a sealed union type system for lambda calculus expressions:

- `ParseTree`: Root structure containing `InputRange` and `ParseItem`
- `ParseItem`: Sealed interface with four implementations:
  - `Var`: Variables with name
  - `Lambda`: Lambda expressions with argument name and body
  - `App`: Function applications with function and arguments
  - `Parens`: Parenthesized expressions
- `Case` function: Generic pattern matching over `ParseItem` types

### Input Tracking

The parser tracks source location information:
- `InputRange`: Spans from one location to another
- `InputLocation`: Line/column position in source

## Common Commands

### Building
```bash
go build
```

### Running
```bash
go run .
```

### Testing
```bash
go test ./...
```

### Module Management
```bash
go mod tidy
go mod download
```

### Development Environment
This project uses mise for Go version management. The Go 1.24 toolchain is specified in `mise.toml`.