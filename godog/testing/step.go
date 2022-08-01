package testing

import (
	"github.com/cucumber/godog"
)

type ctx struct {
	ct *godog.ScenarioContext
}

func (c *ctx) Register(expr, stepFunc interface{}) *ctx {
	c.ct.Step(expr, AdaptTestFunction(stepFunc))

	return c
}

func NewContext(c *godog.ScenarioContext) *ctx {
	return &ctx{
		ct: c,
	}
}
