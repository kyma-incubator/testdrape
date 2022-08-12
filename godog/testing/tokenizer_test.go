package testing

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InterpolateTestSuite struct {
	suite.Suite
	emptyGlossary *tokenizer
	tokenlessExp  string
	tokenizedExp  string
}

func TestInterpolateTestSuite(t *testing.T) {
	suite.Run(t, new(InterpolateTestSuite))
}

func (s *InterpolateTestSuite) SetupSuite() {
	s.emptyGlossary = NewTokenizer()
	s.tokenlessExp = "some expression without tokens"
	s.tokenizedExp = "some expression with a {Token}"
}

func (s *InterpolateTestSuite) TestInterpolate() {
	tcs := []struct {
		scenario   string
		glossary   *tokenizer
		expression string
		res        string
		err        error
	}{
		{
			scenario:   "an empty expression :: for a nil tokens :: returns an empty expression",
			glossary:   NewTokenizer(),
			expression: "",
			res:        "",
			err:        nil,
		},
		{
			scenario:   "a tokenless expression :: for a nil tokens :: returns the original expression",
			glossary:   NewTokenizer(),
			expression: s.tokenlessExp,
			res:        s.tokenlessExp,
			err:        nil,
		},
		{
			scenario:   "a tokenized expression :: for a nil tokens :: returns an error",
			glossary:   NewTokenizer(),
			expression: "some expression with a {Token}",
			res:        "",
			err:        ErrUnknownToken,
		},
		{
			scenario:   "a tokenized expression :: for a configured tokens :: returns a transformed expression",
			glossary:   NewTokenizer().Int("Token"),
			expression: `some expression with a {Token}`,
			res:        fmt.Sprintf(`some expression with a %s`, tokenRegExps[intToken]),
			err:        nil,
		},
		{
			scenario:   "a multiple tokenized expression :: for a configured tokens :: returns a transformed expression",
			glossary:   NewTokenizer().Int("Token"),
			expression: `some expression with a {Token} and another {Token}`,
			res:        fmt.Sprintf(`some expression with a %[1]s and another %[1]s`, tokenRegExps[intToken]),
			err:        nil,
		},
		{
			scenario:   "a multiple tokenized expression :: for a configured tokens :: returns an error",
			glossary:   NewTokenizer().Int("Token"),
			expression: `some expression with a {Token} and another {Token} and unregistered {Sub Token}`,
			res:        "",
			err:        ErrUnknownToken,
		},
		{
			scenario:   "a multiple tokenized expression :: for a configured tokens :: returns a transformed expression",
			glossary:   NewTokenizer().Int("Token").String("Sub Token"),
			expression: `some expression with a {Token} and another {Sub Token} and unregistered {Token}`,
			res:        fmt.Sprintf(`some expression with a %[1]s and another %[2]s and unregistered %[1]s`, tokenRegExps[intToken], tokenRegExps[stringToken]),
			err:        nil,
		},
		{
			scenario:   "a tokenized expression :: for a primitive types :: returns a transformed expression",
			glossary:   NewTokenizer(),
			expression: `an expression with an {int}, {string} and {word} tokens`,
			res:        fmt.Sprintf(`an expression with an %s, %s and %s tokens`, tokenRegExps[intToken], tokenRegExps[stringToken], tokenRegExps[wordToken]),
			err:        nil,
		},
		{
			scenario:   "a tokenized expression :: for a primitive string type :: returns a transformed expression",
			glossary:   NewTokenizer(),
			expression: `some expression with an {string} token`,
			res:        fmt.Sprintf(`some expression with an %s token`, tokenRegExps[stringToken]),
			err:        nil,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.scenario, func() {
			glossary := tc.glossary

			res, err := glossary.Interpolate(tc.expression)

			if tc.err == nil {
				s.Nil(err)
			} else {
				s.True(errors.Is(err, tc.err))
			}
			s.Equal(tc.res, res)
		})
	}
}

func TestSteps(t *testing.T) {
	tcs := []struct {
		scenario       string
		tokenizer      *tokenizer
		step           string
		stepDefinition string
		matches        bool
	}{
		{
			scenario:       "with an int primitive",
			tokenizer:      NewTokenizer(),
			step:           "there are 34 pods",
			stepDefinition: "^there are {int} pods$",
			matches:        true,
		},
		{
			scenario:       "with a negative int primitive",
			tokenizer:      NewTokenizer(),
			step:           "there is -567",
			stepDefinition: "^there is {int}$",
			matches:        true,
		},
		{
			scenario:       "with a word primitive",
			tokenizer:      NewTokenizer(),
			step:           "there is a statement",
			stepDefinition: "^there is a {word}$",
			matches:        true,
		},
		{
			scenario:       "with a string primitive",
			tokenizer:      NewTokenizer(),
			step:           "there is a 'statement string'",
			stepDefinition: "^there is a {string}$",
			matches:        true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			re, _ := tc.tokenizer.Interpolate(tc.stepDefinition)

			match := regexp.MustCompile(re).MatchString(tc.step)

			assert.Equal(t, tc.matches, match)
		})
	}
}
