package functions

import "dsl/variables"

type Argument struct {
	Identifier string
	Type       variables.Type
}

type FunctionDefinition struct {
	ArgumentList       []Argument
	ReturnType         variables.Type
}

type FullFunctionDefinition struct {
	Name               string
	FunctionDefinition FunctionDefinition
}

// Verify an argument list of symbols.
// If the argument list length and type of each argument does not match, return false.
func (fn FunctionDefinition) ValidateArgumentList(symbols []variables.Symbol) bool {
	if len(symbols) != len(fn.ArgumentList) {
		return false
	}

	for i := range fn.ArgumentList {
		if symbols[i].Type != fn.ArgumentList[i].Type {
			return false
		}
	}

	return true
}
