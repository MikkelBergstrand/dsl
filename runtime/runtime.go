package runtime

import (
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

func (runtime *Runtime) GetLabel(label string) int {
	value, ok := runtime.Labels[label]
	if !ok {
		log.Fatalf("No such label %s", label)
	}
	return value
}

func (runtime *Runtime) PushAddress(addressStart int) {
	last_adr := runtime.AddressStack.Peek()
	runtime.AddressStack.Push(last_adr + addressStart)
}

func (runtime *Runtime) PopAddress() {
	runtime.AddressStack.Pop()
}

func (runtime *Runtime) PushCall(offset int) {
	runtime.CallStack.Push(runtime.Programcounter + offset)
}

func (runtime *Runtime) PopCall() {
	val := runtime.CallStack.Pop()
	runtime.Programcounter = val - 1
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
		//color.Println(color.Yellow, reflect.TypeOf(runtime.Instructions[runtime.Programcounter]), "PC = ", runtime.Programcounter)
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
}
