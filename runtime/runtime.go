package runtime

import (
	"dsl/color"
	"dsl/structure"
	"dsl/variables"
	"fmt"
	"log"
	"reflect"
)

type Runtime struct {
	Instructions   []Instruction
	Programcounter int
	Labels         map[string]int
	CallStack      structure.Stack[ActivationRegister]
	Variables      []any
}

type ActivationRegister struct {
	SavedPC      int
	Retval       variables.Symbol
	AddressStack structure.Stack[int]
	AddressBegin int
	StackTop     int
}

func New() Runtime {
	runTime := Runtime{
		Variables: make([]any, 1000),
		Labels:    map[string]int{},
	}

	first_ar := ActivationRegister{
		SavedPC:      0,
		AddressBegin: 0,
	}
	first_ar.AddressStack.Push(0)

	runTime.CallStack.Push(first_ar)

	return runTime
}

func (runtime *Runtime) GetLabel(label string) int {
	value, ok := runtime.Labels[label]
	if !ok {
		log.Fatalf("No such label %s", label)
	}
	return value
}

func (ar *ActivationRegister) PushAddress() {
	ar.AddressStack.Push(ar.StackTop + 1)
	fmt.Println("Adress stack pushed at ", ar)
}

func (ar *ActivationRegister) PopAddress() {
	ar.StackTop = ar.AddressStack.Pop()
}

func (runtime *Runtime) PushCall(offset int, relativeFunctionSymbol int) {
	// The address stack in the function must have an address stack equal to how it looked
	// when the function was defined.
	top_of_callstack := runtime.CallStack.Peek()
	var addr_stack structure.Stack[int]
	for i := range len(top_of_callstack.AddressStack)-relativeFunctionSymbol {
		addr_stack.Push(top_of_callstack.AddressStack[i])
	}
	// The beginning of the next address stack then begins at the next avaiable address
	addr_stack.Push(top_of_callstack.StackTop + 1)
	runtime.CallStack.Push(ActivationRegister{
		SavedPC:      runtime.Programcounter + offset,
		AddressStack: addr_stack,
		AddressBegin: top_of_callstack.StackTop + 1,
	})

	fmt.Println("PushCall with AR = ", runtime.CallStack.Peek(), relativeFunctionSymbol)
}

func (runtime *Runtime) PopCall() {
	val := runtime.CallStack.Pop()
	runtime.Programcounter = val.SavedPC - 1
	fmt.Println("PopCall, AR = ", runtime.CallStack.Peek())
}

// Add the set of instructions. Return the first and last index of the inserted instructions.
func (runtime *Runtime) LoadInstructions(instructions []InstructionLabelPair) (start int, end int) {
	for _, pair := range instructions {
		runtime.Instructions = append(runtime.Instructions, pair.Instruction)
		if pair.Label != "" {
			runtime.Labels[pair.Label] = len(runtime.Instructions) - 1
		}
	}
	start, end = len(runtime.Instructions)-len(instructions), len(runtime.Instructions)-1
	return start, end
}

func (runTime *Runtime) NextInstruction() int {
	return len(runTime.Instructions)
}

func (runtime *Runtime) Run() {

	fmt.Println(runtime.Labels)
	for i, instr := range runtime.Instructions {
		fmt.Println(i, reflect.TypeOf(instr), instr)
	}

	for runtime.Programcounter < len(runtime.Instructions) {
		color.Println(color.Yellow, reflect.TypeOf(runtime.Instructions[runtime.Programcounter]), "PC = ", runtime.Programcounter)
		runtime.Instructions[runtime.Programcounter].Execute(runtime)
		runtime.Programcounter += 1
	}
}

func (r *Runtime) AddressFromSymbol(symbol variables.Symbol) int {
	top_of_callstack := r.CallStack.PeekRef()

	fmt.Println("Resolving address symbol", symbol, top_of_callstack.AddressStack)
	ar := top_of_callstack.AddressStack[len(top_of_callstack.AddressStack)-1-symbol.Scope]
	return ar + symbol.Offset
}

func (s *Runtime) Get(symbol variables.Symbol) any {
	addr := s.AddressFromSymbol(symbol)
	resolve := s.Variables[addr]

	fmt.Println("Get", symbol, "val=", resolve, "addr=", addr)
	return resolve
}

func (r *Runtime) GetInt(symbol variables.Symbol) int {
	val := r.Get(symbol).(int)

	return val
}

func (r *Runtime) GetBool(symbol variables.Symbol) bool {
	return r.Get(symbol).(bool)
}

func (s *Runtime) Set(symbol variables.Symbol, value any) {
	addr := s.AddressFromSymbol(symbol)
	s.Variables[addr] = value
	stack_top := &s.CallStack.PeekRef().StackTop
	if addr > *stack_top {
		*stack_top = addr
	}
	fmt.Println("Set", symbol, "value=", value, "addr=", addr)
}
