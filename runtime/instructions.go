package runtime

import (
	"dsl/color"
	"dsl/functions"
	"dsl/variables"
	"fmt"
	"slices"
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
type BooleanOperator int

const (
	EQUALS BooleanOperator = iota
	NOTEQUALS
	LESS
	LESSOREQUAL
	GREATER
	GREATEROREQUAL
	AND
	OR
	NOT
)

func (op BooleanOperator) IsValidFor(t variables.Type) bool {
	legalBools := []BooleanOperator{EQUALS, NOTEQUALS, AND, OR, NOT}
	legalInts := []BooleanOperator{EQUALS, NOTEQUALS, LESS, LESSOREQUAL, GREATER, GREATEROREQUAL}

	switch t {
	case variables.BOOL:
		return slices.Contains(legalBools, op)
	case variables.INT:
		return slices.Contains(legalInts, op)
	}

	return false
}

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
		runtime.Set(instr.Result, runtime.GetInt(instr.A)+runtime.GetInt(instr.B))
	case MULT:
		runtime.Set(instr.Result, runtime.GetInt(instr.A)*runtime.GetInt(instr.B))
	case DIV:
		runtime.Set(instr.Result, runtime.GetInt(instr.A)/runtime.GetInt(instr.B))
	case SUB:
		runtime.Set(instr.Result, runtime.GetInt(instr.A)-runtime.GetInt(instr.B))
	}
}

type InstrCompareInt struct {
	A        variables.Symbol
	B        variables.Symbol
	Operator BooleanOperator
	Result   variables.Symbol
}

func (instr *InstrCompareInt) Execute(runtime *Runtime) {
	switch instr.Operator {
	case EQUALS:
		runtime.Set(instr.Result, runtime.GetInt(instr.A) == runtime.GetInt(instr.B))
	case NOTEQUALS:
		runtime.Set(instr.Result, runtime.GetInt(instr.A) != runtime.GetInt(instr.B))
	case LESS:
		runtime.Set(instr.Result, runtime.GetInt(instr.A) < runtime.GetInt(instr.B))
	case LESSOREQUAL:
		runtime.Set(instr.Result, runtime.GetInt(instr.A) <= runtime.GetInt(instr.B))
	case GREATER:
		runtime.Set(instr.Result, runtime.GetInt(instr.A) > runtime.GetInt(instr.B))
	case GREATEROREQUAL:
		fmt.Println("Comparing")
		runtime.Set(instr.Result, runtime.GetInt(instr.A) >= runtime.GetInt(instr.B))
	}
}

type InstrCompareBool struct {
	A        variables.Symbol
	B        variables.Symbol
	Operator BooleanOperator
	Result   variables.Symbol
}

func (instr *InstrCompareBool) Execute(runtime *Runtime) {
	switch instr.Operator {
	case EQUALS:
		runtime.Set(instr.Result, runtime.GetBool(instr.A) == runtime.GetBool(instr.B))
	case NOTEQUALS:
		runtime.Set(instr.Result, runtime.GetBool(instr.A) != runtime.GetBool(instr.B))
	case AND:
		runtime.Set(instr.Result, runtime.GetBool(instr.A) && runtime.GetBool(instr.B))
	case OR:
		runtime.Set(instr.Result, runtime.GetBool(instr.A) || runtime.GetBool(instr.B))
	}
}

type InstrJmp struct {
	Label string
}

func (instr *InstrJmp) Execute(runtime *Runtime) {
	runtime.Programcounter = runtime.Labels[instr.Label] - 1 // decrement, since it is autoincremented
}

type InstrLoadImmediate struct {
	Dest  variables.Symbol
	Value any
}

func (instr *InstrLoadImmediate) Execute(runtime *Runtime) {
	runtime.Set(instr.Dest, instr.Value)
}

func (instr *InstrAssign) Execute(runtime *Runtime) {
	runtime.Set(instr.Dest, runtime.Get(instr.Source))
}

type InstructionEcho struct {
	A variables.Symbol
}

func (instr *InstructionEcho) Execute(runtime *Runtime) {
	color.Println(color.Green, "Echo:", runtime.Variables[runtime.AddressFromSymbol(instr.A)])
}

type InstrCallFunction struct {
	PreludeLength int
	AddressStart  int
}

func (instr *InstrCallFunction) Execute(runtime *Runtime) {
	fmt.Println("Calling function: ", instr)
	// Account for prelude length
	runtime.PushCall(instr.PreludeLength)
	runtime.PushAddress(instr.AddressStart)
}

type InstrBeginScope struct {
	AddressStart int
}

func (instr *InstrBeginScope) Execute(runtime *Runtime) {
	fmt.Println("Calling function: ", instr)
	runtime.PushAddress(instr.AddressStart)
}

type InstrEndScope struct{}

func (instr *InstrEndScope) Execute(runtime *Runtime) {
	fmt.Println("Calling function: ", instr)
	runtime.PopAddress()
}

type InstrExitFunction struct {
	RetVal int
}

func (instr *InstrExitFunction) Execute(runtime *Runtime) {
	fmt.Println("Exiting function")
	runtime.PopCall()
	runtime.PopAddress()
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
