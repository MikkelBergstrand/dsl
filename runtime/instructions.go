package runtime

import (
	"dsl/color"
	"dsl/variables"
	"fmt"
	"slices"
)

const (
	ADD Operator = iota
	SUB
	MULT
	DIV
	MOD
)

type InstructionLabelPair struct {
	Instruction Instruction
	Label       string
}
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
	case MOD:
		runtime.Set(instr.Result, runtime.GetInt(instr.A)%runtime.GetInt(instr.B))
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
	runtime.Programcounter = runtime.GetLabel(instr.Label) - 1 // decrement, since it is autoincremented
}

type InstrJmpVar struct {
	FunctionPointer variables.FunctionPointer
}

func (instr *InstrJmpVar) Execute(runtime *Runtime) {
	runtime.Programcounter = runtime.GetLabel(string(instr.FunctionPointer)) - 1
}

type InstrJmpIf struct {
	Label     string
	Condition variables.Symbol
}

func (instr *InstrJmpIf) Execute(runtime *Runtime) {
	if !runtime.GetBool(instr.Condition) {
		runtime.Programcounter = runtime.GetLabel(instr.Label) - 1
	}
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
	color.Println(color.Green, runtime.Variables[runtime.AddressFromSymbol(instr.A)])
}

type InstrCallFunction struct {
	PreludeLength   int
	Arguments       []variables.Symbol
	RetVal          variables.Symbol
	FuncScopeOffset int
	SymbolicLabel   variables.Symbol
}

func (instr *InstrCallFunction) Execute(runtime *Runtime) {
	//Fetch and copy argument values
	var arg_values []any
	for i := range instr.Arguments {
		arg_values = append(arg_values, runtime.Get(instr.Arguments[i]))
	}
	//Fetch func_ptr
	func_ptr := runtime.Get(instr.SymbolicLabel).(variables.FunctionPointer)

	// Bind return value
	top_ar := runtime.CallStack.PeekRef()
	top_ar.Retval = instr.RetVal
	fmt.Println("Bound ret val to", top_ar.Retval)

	// Account for prelude length
	runtime.PushCall(instr.PreludeLength, instr.FuncScopeOffset)

	// Once "inside" the function, load argument values
	for i := range arg_values {
		runtime.Set(variables.Symbol{Offset: i, Scope: 0, Type: instr.Arguments[i].Type}, arg_values[i])
	}
	jmp_instr := InstrJmpVar{
		FunctionPointer: func_ptr,
	}
	jmp_instr.Execute(runtime)
}

type InstrBeginScope struct{}

func (instr *InstrBeginScope) Execute(runtime *Runtime) {
	runtime.CallStack.PeekRef().PushAddress()
}

type InstrEndScope struct{}

func (instr *InstrEndScope) Execute(runtime *Runtime) {
	runtime.CallStack.PeekRef().PopAddress()
}

type InstrExitFunction struct {
	RetVal variables.Symbol
}

func (instr *InstrExitFunction) Execute(runtime *Runtime) {
	src_val := runtime.Get(instr.RetVal)

	runtime.PopCall()
	top_ar := runtime.CallStack.PeekRef()

	fmt.Println("Ret val on exit", top_ar.Retval, src_val)
	runtime.Set(top_ar.Retval, src_val)
}

// Does nothing.
type InstrNOP struct{}

func (instr *InstrNOP) Execute(runtime *Runtime) {}
