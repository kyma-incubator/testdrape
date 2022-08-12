package testing

type tokenizerOption func(*tokenizerOptions)

type tokenizerOptions struct {
	leftDelimiter  string
	rightDelimiter string
}

// WithDelimiters exposes the option to configure token's left and right delimiters:
//   Usage:
//       NewTokenizer(WithDelimiters("<", "}")).Int("Token").Interpolate("a string with a <Token}")
func WithDelimiters(left, right string) tokenizerOption {
	return func(opts *tokenizerOptions) {
		opts.leftDelimiter = left
		opts.rightDelimiter = right
	}
}
