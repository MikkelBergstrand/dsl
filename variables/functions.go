package variables

type Argument struct {
	Identifier string
	Type       Type
}

type FunctionDefinition struct {
	ArgumentList []Argument
	ReturnType   Type
}

// Holds the label of the function it is referring to.
type FunctionPointer string

// Verify an argument list of symbols.
// If the argument list length and type of each argument does not match, return false.
func (fn FunctionDefinition) ValidateArgumentList(symbols []Symbol) bool {
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
