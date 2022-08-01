package testing

import (
	"fmt"
	"reflect"
)

// AdaptTestFunction adapts a test function to conform to the Godog runner.
//
// The Godog step runner relies on returning an error for the failed steps.
// As it contradicts the go testing paradigm, this decorator function converts
// a conventional go test function to the one the Godog runner expects.
//
// Performs a conversion:
//     func (t *testing.T, arg... any) -> func(arg... any) error
func AdaptTestFunction(testFunc interface{}) interface{} {
	validateTestFunction(reflect.ValueOf(testFunc).Type())
	nilError := reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())

	ins := fetchTestFunctionArguments(reflect.TypeOf(testFunc))
	outType := nilError.Type()
	adapterType := reflect.FuncOf(ins, []reflect.Type{outType}, false)

	return reflect.MakeFunc(adapterType, func(in []reflect.Value) []reflect.Value {
		// Inject the *testing.T argument to the argument list.
		t := newT()
		ins := append([]reflect.Value{reflect.ValueOf(t)}, in...)

		// Invoke the test function.
		if reflect.ValueOf(testFunc).Call(ins); t.Err() != nil {
			return []reflect.Value{reflect.ValueOf(t.Err())}
		}

		return []reflect.Value{nilError}
	}).Interface()
}

func validateTestFunction(t reflect.Type) {
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("handler is not a func, but: %T", t))
	}

	if outs := t.NumOut(); outs > 0 {
		panic(fmt.Sprintf("handler must not return any values, but returns %d", outs))
	}

	if ins := t.NumIn(); ins == 0 {
		panic(fmt.Sprintf("handler must have at least *testing.T argument"))
	}

	if test := t.In(0); test.Kind() != reflect.Pointer || test.Elem().Implements(reflect.TypeOf((*testingT)(nil)).Elem()) {
		panic(fmt.Sprintf("the first argument must be of type *testing.T"))
	}
}

func fetchTestFunctionArguments(t reflect.Type) []reflect.Type {
	ins := make([]reflect.Type, t.NumIn()-1)
	for i := 1; i < t.NumIn(); i++ {
		ins[i-1] = t.In(i)
	}

	return ins
}
