package runtime

import (
	"dsl/instructions"
	"dsl/storage"
)

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
