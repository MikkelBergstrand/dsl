package runtime

import (
	"dsl/color"
	"dsl/structure"
	"dsl/variables"
	"fmt"
	"reflect"
)

type Runtime struct {
	Instructions   []Instruction
	Programcounter int
	Labels         map[string]int
	CallStack      structure.Stack[int]
	AddressStack   structure.Stack[int]
	Variables      []any
	Addresspointer int
}

const NIL_PROGRAM_COUNTER = -1000

type ActivationRegister struct {
	AddressStart   int
	Programcounter int
}

func New() Runtime {
	runTime := Runtime{
		Variables: make([]any, 1000),
		Labels:    map[string]int{},
	}

	runTime.CallStack.Push(0)
	runTime.AddressStack.Push(0)
	return runTime
}

func (runtime *Runtime) PushAddress(addressStart int) {
	last_adr := runtime.AddressStack.Peek()
	runtime.AddressStack.Push(last_adr + addressStart)
	fmt.Println("Pushed address ", runtime.AddressStack.Peek())
}

func (runtime *Runtime) PopAddress() {
	val := runtime.AddressStack.Pop()
	fmt.Println("Popped address ", val)
}

func (runtime *Runtime) PushCall(offset int) {
	runtime.CallStack.Push(runtime.Programcounter + offset)
	fmt.Println("Pushed", runtime.CallStack.Peek())
}

func (runtime *Runtime) PopCall() {
	val := runtime.CallStack.Pop()
	fmt.Println("Popped", val)
	runtime.Programcounter = val - 1
}

// Add the set of instructions. Return the first and last index of the inserted instructions.
func (runtime *Runtime) LoadInstructions(instructions []Instruction) (start int, end int) {
	runtime.Instructions = append(runtime.Instructions, instructions...)
	start, end = len(runtime.Instructions)-len(instructions), len(runtime.Instructions)-1

	fmt.Println("New instructions added: ", start, end)
	return start, end
}

func (runTime *Runtime) NextInstruction() int {
	return len(runTime.Instructions)
}

func (runtime *Runtime) Run() {

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
	ar := r.AddressStack[len(r.AddressStack)-1-symbol.Scope]
	return ar + symbol.Offset
}

func (s *Runtime) Get(symbol variables.Symbol) any {
	resolve := s.Variables[s.AddressFromSymbol(symbol)]
	fmt.Printf("Get() = %v from addr %v, symbol %v\n", resolve, s.AddressFromSymbol(symbol), symbol)
	return resolve
}

func (r *Runtime) GetInt(symbol variables.Symbol) int {
	return r.Get(symbol).(int)
}

func (r *Runtime) GetBool(symbol variables.Symbol) bool {
	return r.Get(symbol).(bool)
}

func (s *Runtime) Set(symbol variables.Symbol, value any) {
	addr := s.AddressFromSymbol(symbol)
	s.Variables[addr] = value
	fmt.Printf("Set %d to %d at addr %d\n", symbol, value, addr)
}
