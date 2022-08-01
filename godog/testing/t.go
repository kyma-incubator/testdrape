package testing

import (
	"fmt"
)

// testingT is an interface wrapper around *testing.T
type testingT interface {
	Errorf(format string, args ...interface{})
}

// T is the minimal testify-conforming testing object.
type T struct {
	err error
}

func (t *T) Errorf(format string, args ...interface{}) {
	t.err = fmt.Errorf(format, args...)
}
func (t *T) Err() error {
	return t.err
}

func newT() *T {
	return &T{}
}
