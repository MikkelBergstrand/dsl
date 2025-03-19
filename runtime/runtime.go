package runtime

import (
	"dsl/functions"
	"dsl/instructions"
	"dsl/storage"
	"dsl/variables"
)

func generateGlobalFunctions(runtime *Runtime, storage *storage.Storage) {
	runtime.LoadInstruction(&instructions.InstrDeclareFunction{
		Identifier: "echo",
		ArgumentList: []functions.Argument{
			{Type: variables.INT},
		},
		ReturnType: variables.NONE,
	})

	runtime.LoadInstruction(&instructions.InstructionEcho{
		A: storage.GetVarAddr("i"),
	})
	runtime.LoadInstruction(&instructions.InstrEndDeclareFunction{})
}

type Runtime struct {
	Instructions   []instructions.Instruction
	ProgramCounter int
}

func (runTime *Runtime) LoadInstruction(instruction instructions.Instruction) {
	runTime.Instructions = append(runTime.Instructions, instruction)
}

func (runTime *Runtime) NextInstruction() int {
	return len(runTime.Instructions)
}

func (runtime *Runtime) Run(storage *storage.Storage) {
	for _, instr := range runtime.Instructions {
		instr.Execute(storage)
	}
}
