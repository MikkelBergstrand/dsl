package functions

import "dsl/variables"

type Argument struct {
	Type variables.Type
}

type Function struct {
	ArgumentList []Argument
}
