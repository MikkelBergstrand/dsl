package instructions

import (
	"dsl/functions"
	"dsl/storage"
	"dsl/variables"
	"fmt"
)

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
	A        variables.Symbol
	B        variables.Symbol
	Operator Operator
	Result   variables.Symbol
}

type InstrAssign struct {
	Dest   variables.Symbol
	Source variables.Symbol
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

func (instr *InstrAssign) Execute(storage *storage.Storage) {
	storage.SetInt(instr.Dest, storage.GetInt(instr.Source))
}

type InstructionEcho struct {
	A variables.Symbol
}

func (instr *InstructionEcho) Execute(storage *storage.Storage) {
	fmt.Println(storage.AddressFromSymbol(instr.A))
}

type InstrCallFunction struct {
	Func         *functions.FunctionDefinition
	ArgumentList []int
}

func (instr *InstrCallFunction) Execute(storage *storage.Storage) {
	storage.NewScope()
	for i := range instr.ArgumentList {
		storage.NewIntVariable(instr.Func.ArgumentList[i].Identifier)
	}
}

type InstrExitFunction struct {
	RetVal int
}

func (instr *InstrExitFunction) Execute(storage *storage.Storage) {
	storage.DestroyScope()
	storage.CurrentScope.RetVal = instr.RetVal
}

type InstrDeclareFunction struct {
	AddressPointer int
	Identifier     string
	ArgumentList   []functions.Argument
	ReturnType     variables.Type
}

func (instr *InstrDeclareFunction) Execute(storage *storage.Storage) {
	fmt.Println("Declaring new function", instr)
	storage.Functions[instr.Identifier] = functions.FunctionDefinition{
		ArgumentList:   instr.ArgumentList,
		AddressPointer: instr.AddressPointer,
	}
}

type InstrEndDeclareFunction struct {
}

func (instr *InstrEndDeclareFunction) Execute(storage *storage.Storage) {
	fmt.Println("End declaration of function.")
}
