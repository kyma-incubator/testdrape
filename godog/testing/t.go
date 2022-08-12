package testing

import (
	"fmt"
	"strings"
)

// testingT is an interface wrapper around *testing.T
type testingT interface {
	Errorf(format string, args ...interface{})
}

// T is the minimal testify-conforming testing object.
type T struct {
	err error
}

// Logf formats its arguments according to the format, analogous to Printf. A final newline is added if not provided.
func (t *T) Logf(format string, args ...any) {
	fmt.Printf(strings.Join(strings.Split(fmt.Sprintf(format, args...), "\n"), "\n"))
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
