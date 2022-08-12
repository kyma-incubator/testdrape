package testing

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type (
	tokenKind int
	tokenRE   string
	tokens    = map[string]tokenKind
)

const (
	// The regular expression pattern for the token in the step definition.
	tokenRegExp = `([\w][\w\s]*)`

	// Supported token types.
	stringToken tokenKind = iota
	wordToken
	intToken
)

var (
	ErrUnknownToken = errors.New("token is unknown")

	tokenRegExps = map[tokenKind]tokenRE{
		stringToken: `("([^"\\]*(\\.[^"\\]*)*)"|'([^'\\]*(\\.[^'\\]*)*)')`,
		wordToken:   `([^\s]+)`,
		intToken:    `(-?\d+)`,
	}
)

// tokenizer provides the support of the tokens in the step definitions.
// Provides the list of methods to register tokens and an Interpolate method
// that replaces the tokens in the string with regular expressions.
type tokenizer struct {
	tokens         tokens
	leftDelimiter  string
	rightDelimiter string
}

func NewTokenizer(opts ...tokenizerOption) *tokenizer {
	tokenizer := tokenizer{
		leftDelimiter:  "{",
		rightDelimiter: "}",
		tokens:         make(tokens),
	}

	// Register predefined tokens.
	tokenizer.tokens["string"] = stringToken
	tokenizer.tokens["int"] = intToken
	tokenizer.tokens["word"] = wordToken

	options := tokenizerOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	// Process the token delimiters arguments.
	if options.leftDelimiter != "" {
		tokenizer.leftDelimiter = options.leftDelimiter
	}
	if options.rightDelimiter != "" {
		tokenizer.rightDelimiter = options.rightDelimiter
	}

	return &tokenizer
}

// Interpolate performs the replacement of the tokens in the string with
// respective regular expressions.
//
// Usage:
//   NewTokenizer().Int("Token").Interpolate("string with a {Token}")
//       > `string with a (-?\d+)`
func (g *tokenizer) Interpolate(stepDef string) (string, error) {
	tokenRegExp := regexp.MustCompile(fmt.Sprintf(`%s%s%s`, g.leftDelimiter, tokenRegExp, g.rightDelimiter))
	for _, match := range tokenRegExp.FindAllStringSubmatch(stepDef, -1) {
		token := match[1]
		expr := match[0]
		kind, ok := g.tokens[token]
		if !ok {
			return "", fmt.Errorf("%w: %s is not registered", ErrUnknownToken, expr)
		}

		stepDef = strings.ReplaceAll(stepDef, expr, string(tokenRegExps[kind]))
	}

	return stepDef, nil
}

func (g *tokenizer) String(n string) *tokenizer {
	g.tokens[n] = stringToken
	return g
}

func (g *tokenizer) Int(n string) *tokenizer {
	g.tokens[n] = intToken
	return g
}

func (g *tokenizer) Word(n string) *tokenizer {
	g.tokens[n] = wordToken
	return g
}
