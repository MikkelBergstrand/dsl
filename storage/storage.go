package storage

import (
	"dsl/runtime"
	"dsl/variables"
	"fmt"
	"log"
	"strconv"
)

type Storage struct {
	CurrentScope *scoped_storage
	Scopes       []scoped_storage
	LabelIndex   int //Used for auto-generated labels. They must be unique across scopes.
	NextLabel    string
}

type scoped_storage struct {
	Parent       *scoped_storage
	Variables    map[string]variables.SymbolTableEntry
	Offset       int
	Instructions []runtime.InstructionLabelPair //Instructions and associated label from statements/expressions in the local scope.
}

func newScopedStorage() scoped_storage {
	return scoped_storage{
		Variables: make(map[string]variables.SymbolTableEntry),
	}
}

func NewStorage() Storage {
	storage := Storage{}
	storage.Scopes = append(storage.Scopes, newScopedStorage())
	storage.CurrentScope = &storage.Scopes[0]
	return storage
}

func (s *Storage) NewFunction(name string, definition variables.TypeDefinition) {
	func_symbol, err := s.NewVariable(definition, name)
	if err != nil {
		log.Fatal(err)
	}

	label := s.NewAutoLabel()
	s.LoadInstruction(&runtime.InstrLoadFunction{
		Symbol: *func_symbol,
		Label:  label,
	})

	s.newFunctionScope(definition)
	s.NewLabel(label)
}

func (s *Storage) newFunctionScope(definition variables.TypeDefinition) {
	s.NewScope()

	// Create variable entries for the arguments. They are placed first in the function's symbol table
	for _, arg := range definition.ArgumentList {
		_, err := s.NewVariable(arg.Definition, arg.Identifier)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *Storage) NewScope() *scoped_storage {
	s.Scopes = append(s.Scopes, newScopedStorage())
	s.Scopes[len(s.Scopes)-1].Parent = s.CurrentScope
	s.CurrentScope = &s.Scopes[len(s.Scopes)-1]
	return s.CurrentScope
}

func (s *Storage) DestroyFunctionScope(runTime *runtime.Runtime) (int, int) {
	start, end := runTime.LoadInstructions(s.CurrentScope.Instructions)

	s.CurrentScope = s.CurrentScope.Parent

	return start, end
}

func (s *Storage) DestroyScope() {
	instructions := s.CurrentScope.Instructions

	s.CurrentScope = s.CurrentScope.Parent
	s.CurrentScope.Instructions = append(s.CurrentScope.Instructions, instructions...)

}

func (s *Storage) NewLiteral(vartype variables.TypeDefinition) variables.Symbol {
	s.CurrentScope.Offset += 1
	sym := variables.Symbol{Scope: 0, Offset: s.CurrentScope.Offset - 1, Type: vartype}
	fmt.Println("New literal", sym)
	return sym
}

func (s *Storage) NewVariable(vartype variables.TypeDefinition, name string) (*variables.Symbol, error) {
	_, exists := s.CurrentScope.Variables[name]
	if exists {
		log.Fatalf("Redeclaration of variable: %s\n", name)
		return nil, fmt.Errorf("redeclaration of variable: %s\n", name)
	}

	addr := s.CurrentScope.Offset
	s.CurrentScope.Variables[name] = variables.SymbolTableEntry{
		Type:   vartype,
		Offset: addr,
	}
	s.CurrentScope.Offset += 1

	fmt.Println("New variable", s.CurrentScope.Variables[name], name)

	return &variables.Symbol{Scope: 0, Offset: s.CurrentScope.Offset - 1, Type: vartype}, nil
}

func (s *Storage) GetVarAddr(name string) (variables.Symbol, error) {
	scope := s.CurrentScope
	symbol, ok := variables.SymbolTableEntry{}, false
	scopeOffset := 0
	for {
		symbol, ok = (*scope).Variables[name]
		if ok {
			break
		}
		scope = (*scope).Parent
		if scope == nil {
			break
		}
		scopeOffset += 1
	}

	if !ok {
		return variables.Symbol{}, fmt.Errorf("could not resolve variable name: %s", name)
	}
	return variables.Symbol{
		Scope:  scopeOffset,
		Offset: symbol.Offset,
		Type:   symbol.Type,
	}, nil
}

func (s *Storage) LoadLabeledInstruction(instruction runtime.Instruction, label string) *runtime.InstructionLabelPair {
	instr := s.LoadInstruction(instruction)
	instr.Label = label
	return instr
}
func (s *Storage) LoadInstruction(instruction runtime.Instruction) *runtime.InstructionLabelPair {
	s.CurrentScope.Instructions = append(s.CurrentScope.Instructions, runtime.InstructionLabelPair{
		Instruction: instruction,
		Label:       s.NextLabel,
	})
	s.NextLabel = ""
	return &s.CurrentScope.Instructions[len(s.CurrentScope.Instructions)-1]
}

func (s *Storage) InsertInstructionAt(instruction runtime.Instruction, label string, offset int) {

}

func (s *Storage) NewLabel(label string) {
	if s.NextLabel != "" {
		log.Fatalf("Label %s overridden by %s!", s.NextLabel, label)
	}
	s.NextLabel = label
}

func (s *Storage) NewAutoLabel() (label string) {
	s.LabelIndex += 1
	return strconv.Itoa(s.LabelIndex)
}
