package instructionset

import "dsl/storage"

type Instruction interface {
	Execute(*storage.Storage)
}

type Operator int

const (
	ADD Operator = iota
	SUB
	MULT
	DIV
)

type InstrArithmetic struct {
	A        int
	B        int
	Operator Operator
	Result   int
}

func (instr *InstrArithmetic) Execute(storage *storage.Storage) {
	switch instr.Operator {
	case ADD:
		storage.SetInt(instr.Result, storage.GetInt(instr.A)+storage.GetInt(instr.B))
	case MULT:
		storage.SetInt(instr.Result, storage.GetInt(instr.A)*storage.GetInt(instr.B))
	case DIV:
		storage.SetInt(instr.Result, storage.GetInt(instr.A)/storage.GetInt(instr.B))
	case SUB:
		storage.SetInt(instr.Result, storage.GetInt(instr.A)-storage.GetInt(instr.B))
	}
}
