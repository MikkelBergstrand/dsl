package variables

type Type int

const (
	INT Type = iota
	NONE
)

type Symbol struct {
	Scope  int
	Offset int
	Type   Type
}

type SymbolTableEntry struct {
	Offset int
	Type   Type
}
