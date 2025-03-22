package runtime

import (
	"dsl/functions"
	"dsl/storage"
	"dsl/structure"
	"dsl/variables"
	"fmt"
	"log"
	"reflect"
)

func generateGlobalFunctions(runtime *Runtime, storage *storage.Storage) {
	storage.NewFunction("echo", functions.FunctionDefinition{
		InstructionPointer: runtime.NextInstruction(),
		ArgumentList: []functions.Argument{
			{Type: variables.INT, Identifier: "i"},
		},
		ReturnType: variables.NONE,
	})

	storage.NewIntVariable("i")
	runtime.LoadInstruction(&InstructionEcho{
		A: storage.GetVarAddr("i"),
	})
	runtime.LoadInstruction(&InstrExitFunction{})
	storage.DestroyScope()

	storage.NewFunction("echo2", functions.FunctionDefinition{
		InstructionPointer: runtime.NextInstruction(),
		ArgumentList: []functions.Argument{
			{Type: variables.INT, Identifier: "i"},
			{Type: variables.INT, Identifier: "j"},
		},
		ReturnType: variables.NONE,
	})
	storage.NewIntVariable("i")
	storage.NewIntVariable("j")
	runtime.LoadInstruction(&InstructionEcho{
		A: storage.GetVarAddr("i"),
	})
	runtime.LoadInstruction(&InstructionEcho{
		A: storage.GetVarAddr("j"),
	})
	runtime.LoadInstruction(&InstrExitFunction{})
	storage.DestroyScope()

	runtime.Programcounter = runtime.NextInstruction()
}

type Runtime struct {
	Instructions   []Instruction
	Programcounter int
	ARStack        structure.Stack[ActivationRegister]
	Variables      []any
	Addresspointer int
}

type ActivationRegister struct {
	AddressStart   int
	Programcounter int
}

func New(storage *storage.Storage) Runtime {
	runTime := Runtime{
		Variables: make([]any, 1000),
	}
	generateGlobalFunctions(&runTime, storage)

	runTime.ARStack.Push(ActivationRegister{})
	fmt.Println("New AR: ", runTime.Programcounter, len(runTime.ARStack))
	return runTime
}

func (runtime *Runtime) NewAR(addressStart int, preludeLength int) {
	last_ar := runtime.ARStack.PeekRef()
	last_ar.Programcounter = runtime.Programcounter + preludeLength - 1

	runtime.ARStack.Push(ActivationRegister{AddressStart: last_ar.AddressStart + addressStart})
	fmt.Println("New AR: ", runtime.Programcounter, last_ar.AddressStart+addressStart, len(runtime.ARStack))

}

func (runtime *Runtime) PopAR() {
	runtime.ARStack.Pop()
	runtime.Programcounter = runtime.ARStack.Peek().Programcounter
	fmt.Println("Popping AR: ", runtime.Programcounter)
}

func (runTime *Runtime) LoadInstruction(instruction Instruction) {
	runTime.Instructions = append(runTime.Instructions, instruction)
}

func (runTime *Runtime) NextInstruction() int {
	return len(runTime.Instructions)
}

func (runtime *Runtime) Run(storage *storage.Storage) {
	for _, instr := range runtime.Instructions {
		fmt.Println(reflect.TypeOf(instr), instr)
	}

	for runtime.Programcounter < len(runtime.Instructions) {
		fmt.Println(reflect.TypeOf(runtime.Instructions[runtime.Programcounter]), "PC = ", runtime.Programcounter)
		runtime.Instructions[runtime.Programcounter].Execute(runtime)
		runtime.Programcounter += 1
	}
}

func (r *Runtime) AddressFromSymbol(symbol variables.Symbol) int {
	ar := r.ARStack[len(r.ARStack)-1-symbol.Scope]
	return ar.AddressStart + symbol.Offset
}

func (s *Runtime) GetInt(symbol variables.Symbol) int {
	resolve, ok := s.Variables[s.AddressFromSymbol(symbol)].(int)
	if !ok {
		log.Fatalf("Not an integer")
	}
	fmt.Printf("GetInt() = %d from addr %d, symbol %d\n", resolve, s.AddressFromSymbol(symbol), symbol)
	return resolve
}

func (s *Runtime) SetInt(symbol variables.Symbol, value int) {
	addr := s.AddressFromSymbol(symbol)
	s.Variables[addr] = value
	fmt.Printf("Set %d to %d at addr %d\n", symbol, value, addr)
}
