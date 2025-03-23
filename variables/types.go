package variables

import "log"

type Type int

const (
	INT Type = iota
	BOOL
	NONE
)

func (t Type) String() string {
	switch t {
	case INT:
		return "int"
	case BOOL:
		return "bool"
	case NONE:
		return "void"
	}
	log.Panicln("Invalid variable type!")
	return ""
}

type Symbol struct {
	Scope  int
	Offset int
	Type   Type
}

type SymbolTableEntry struct {
	Offset int
	Type   Type
}
