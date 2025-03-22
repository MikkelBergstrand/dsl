package functions

import "dsl/variables"

type Argument struct {
	Identifier string
	Type       variables.Type
}

type FunctionDefinition struct {
	ArgumentList       []Argument
	InstructionPointer int
	ReturnType         variables.Type
}
