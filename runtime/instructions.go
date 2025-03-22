package runtime

import (
	"dsl/functions"
	"dsl/variables"
	"fmt"
)

const (
	ADD Operator = iota
	SUB
	MULT
	DIV
)

type Instruction interface {
	Execute(*Runtime)
}

type Operator int

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

func (instr *InstrArithmetic) Execute(runtime *Runtime) {
	switch instr.Operator {
	case ADD:
		runtime.SetInt(instr.Result, runtime.GetInt(instr.A)+runtime.GetInt(instr.B))
	case MULT:
		runtime.SetInt(instr.Result, runtime.GetInt(instr.A)*runtime.GetInt(instr.B))
	case DIV:
		runtime.SetInt(instr.Result, runtime.GetInt(instr.A)/runtime.GetInt(instr.B))
	case SUB:
		runtime.SetInt(instr.Result, runtime.GetInt(instr.A)-runtime.GetInt(instr.B))
	}
}

type InstrJmp struct {
	NewPC int
}

func (instr *InstrJmp) Execute(runtime *Runtime) {
	runtime.Programcounter = instr.NewPC - 1
}

type InstrLoadImmediate struct {
	Dest  variables.Symbol
	Value any
}

func (instr *InstrLoadImmediate) Execute(runtime *Runtime) {
	runtime.SetInt(instr.Dest, instr.Value.(int))
}

func (instr *InstrAssign) Execute(runtime *Runtime) {
	runtime.SetInt(instr.Dest, runtime.GetInt(instr.Source))
}

type InstructionEcho struct {
	A variables.Symbol
}

func (instr *InstructionEcho) Execute(runtime *Runtime) {
	fmt.Println("Echo:", runtime.Variables[runtime.AddressFromSymbol(instr.A)])
}

type InstrCallFunction struct {
	Func         functions.FunctionDefinition
	ArgumentList []variables.Symbol
	AddressStart int
}

func (instr *InstrCallFunction) Execute(runtime *Runtime) {
	fmt.Println("Calling function: ", instr)
	runtime.NewAR(instr.AddressStart, len(instr.ArgumentList)+2)
}

type InstrExitFunction struct {
	RetVal int
}

func (instr *InstrExitFunction) Execute(runtime *Runtime) {
	fmt.Println("Exiting function")
	runtime.PopAR()
}

type InstrDeclareFunction struct {
	AddressPointer int
	Identifier     string
	ArgumentList   []functions.Argument
	ReturnType     variables.Type
}

func (instr *InstrDeclareFunction) Execute(runtime *Runtime) {
	fmt.Println("Declaring new function", instr)
}

type InstrEndDeclareFunction struct {
}

func (instr *InstrEndDeclareFunction) Execute(runtime *Runtime) {
	fmt.Println("End declaration of function.")
}
