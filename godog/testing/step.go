package testing

import (
	"fmt"

	"github.com/cucumber/godog"
)

type ctx struct {
	ct        *godog.ScenarioContext
	tokenizer *tokenizer
}

func (c *ctx) Register(expr string, stepFunc interface{}) *ctx {
	stepExpr, err := c.tokenizer.Interpolate(expr)
	if err != nil {
		panic(fmt.Sprintf("step definition is incorrect: %v", err))
	}

	c.ct.Step(stepExpr, AdaptTestFunction(stepFunc))

	return c
}

func NewContext(c *godog.ScenarioContext) *ctx {
	return &ctx{
		ct:        c,
		tokenizer: NewTokenizer(),
	}
}

// IntToken registers a named token of int type.
func (c *ctx) IntToken(tokenName string) *ctx {
	c.tokenizer.Int(tokenName)
	return c
}

// WordToken registers a named token represented by a word (an uninterrupted
// series of alphanumerical characters).
func (c *ctx) WordToken(tokenName string) *ctx {
	c.tokenizer.Word(tokenName)
	return c
}

// StringToken registers a named token of string type.
func (c *ctx) StringToken(tokenName string) *ctx {
	c.tokenizer.String(tokenName)
	return c
}
