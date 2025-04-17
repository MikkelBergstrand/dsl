package variables

import (
	"fmt"
	"log"
)

type Type int

const (
	INVALID Type = iota
	INT Type = iota
	BOOL
	FUNC_PTR
	NONE
)

func TypeFromString(s string) (Type, error) {
	switch s {
	case "fptr":
		return FUNC_PTR, nil
	case "int":
		return INT, nil
	case "bool":
		return BOOL, nil
	case "void":
		return NONE, nil
	}

	return NONE, fmt.Errorf("could not resolve %s to a variable type", s)
}

func (t Type) String() string {
	switch t {
	case INT:
		return "int"
	case BOOL:
		return "bool"
	case NONE:
		return "void"
	case FUNC_PTR:
		return "func"
	case INVALID:
		return "invalid"
	}
	log.Panicln("Invalid variable type!")
	return ""
}

type Symbol struct {
	Scope              int
	Offset             int
	Type               Type
	FunctionDefinition FunctionDefinition //optional, if type == FUNC_PTR
}

type SymbolTableEntry struct {
	Offset             int
	Type               Type
	FunctionDefinition FunctionDefinition //optional information, if the value is a function pointer.
}
