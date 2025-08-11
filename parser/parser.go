package parser

import (
	"errors"
	"fmt"

	"github.com/gusbicalho/go-lambda/parse_tree"
	"github.com/gusbicalho/go-lambda/token"
	"github.com/gusbicalho/go-lambda/tokenizer"
)

type ParseResult[v any] struct {
	value            v
	hasConsumedInput bool
	error            error
}

func (r ParseResult[v]) consumedInput() ParseResult[v] {
	r.hasConsumedInput = true
	return r
}

func Parse(tokenizer *tokenizer.Tokenizer) (*parse_tree.ParseTree, error) {
	result := parseTree(tokenizer)
	if result.error != nil {
		return nil, result.error
	}
	if tok := tokenizer.Next(); tok.Type() != token.EOF {
		return nil, errors.New(fmt.Sprint("Expected EOF, found ", tok))
	}
	return result.value, nil
}

func parseTree(tokenizer *tokenizer.Tokenizer) ParseResult[*parse_tree.ParseTree] {
	calleeResult := parseApplicable(tokenizer)
	if calleeResult.error != nil {
		return calleeResult
	}
	return parsePossibleApp(tokenizer, *calleeResult.value)
}

func parseApplicable(tokenizer *tokenizer.Tokenizer) ParseResult[*parse_tree.ParseTree] {
	tok := tokenizer.Peek()
	switch tok.Type() {
	case token.Lambda:
		tokenizer.Next()
		return parseLambda(tokenizer, tok).consumedInput()
	case token.Identifier:
		tokenizer.Next()
		callee := &parse_tree.ParseTree{
			InputLocation: tok.Position,
			Item:          parse_tree.Var{Name: tok.Value},
		}
		return ParseResult[*parse_tree.ParseTree]{value: callee, hasConsumedInput: true}
	case token.LeftParen:
		tokenizer.Next()
		return parseParenTree(tokenizer, tok).consumedInput()

	default:
		return ParseResult[*parse_tree.ParseTree]{error: errors.New(fmt.Sprint("Unexpected token ", tok))}
	}
}

func parseParenTree(tokenizer *tokenizer.Tokenizer, leftParen token.Token) ParseResult[*parse_tree.ParseTree] {
	child := parseTree(tokenizer)
	if child.error != nil {
		return child.consumedInput()
	}
	nextTok := tokenizer.Next()
	if nextTok.Type() == token.RightParen {
		return ParseResult[*parse_tree.ParseTree]{
			value: &parse_tree.ParseTree{
				InputLocation: leftParen.Position,
				Item: parse_tree.Parens{
					Child: *child.value,
				},
			},
			hasConsumedInput: true,
		}
	}
	return ParseResult[*parse_tree.ParseTree]{error: errors.New(fmt.Sprint("Expected ), found ", nextTok))}
}

func parsePossibleApp(tokenizer *tokenizer.Tokenizer, callee parse_tree.ParseTree) ParseResult[*parse_tree.ParseTree] {
	result := parseArgs(tokenizer)
	if len(result.value) == 0 {
		return ParseResult[*parse_tree.ParseTree]{
			value:            &callee,
			hasConsumedInput: result.hasConsumedInput,
			error:            result.error,
		}
	} else {
		app := &parse_tree.ParseTree{
			InputLocation: callee.InputLocation,
			Item: parse_tree.App{
				Callee: callee,
				Args: parse_tree.AppArgs{
					First: result.value[0],
					More:  result.value[1:],
				},
			},
		}
		return ParseResult[*parse_tree.ParseTree]{
			value:            app,
			hasConsumedInput: result.hasConsumedInput,
			error:            result.error,
		}
	}
}

func parseArgs(tokenizer *tokenizer.Tokenizer) ParseResult[[]parse_tree.ParseTree] {
	trees := make([]parse_tree.ParseTree, 0)
	hasConsumedInput := false
	for {
		result := parseApplicable(tokenizer)
		if result.hasConsumedInput {
			hasConsumedInput = true
		}
		if result.error != nil {
			if result.hasConsumedInput {
				return ParseResult[[]parse_tree.ParseTree]{
					hasConsumedInput: true,
					error:            result.error,
				}
			}
			return ParseResult[[]parse_tree.ParseTree]{
				value:            trees,
				hasConsumedInput: hasConsumedInput,
				error:            nil,
			}
		}
		trees = append(trees, *result.value)
	}
}

func parseLambda(tokenizer *tokenizer.Tokenizer, lambdaTok token.Token) ParseResult[*parse_tree.ParseTree] {
	argNameTok := tokenizer.Next()
	if argNameTok.Type() != token.Identifier {
		return ParseResult[*parse_tree.ParseTree]{error: errors.New(fmt.Sprint("Expected identifier, found ", argNameTok))}
	}
	dotTok := tokenizer.Next()
	if dotTok.Type() != token.Dot {
		return ParseResult[*parse_tree.ParseTree]{error: errors.New(fmt.Sprint("Expected ., found ", dotTok))}
	}
	bodyResult := parseTree(tokenizer)
	if bodyResult.error != nil {
		return bodyResult
	}
	return ParseResult[*parse_tree.ParseTree]{
		value: &parse_tree.ParseTree{
			InputLocation: lambdaTok.Position,
			Item: parse_tree.Lambda{
				ArgName: argNameTok.Value,
				Body:    *bodyResult.value,
			},
		},
	}
}
