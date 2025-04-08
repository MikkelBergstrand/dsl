package storage

import (
	"dsl/functions"
	"dsl/runtime"
	"dsl/variables"
	"fmt"
	"log"
)

type Storage struct {
	CurrentScope *scoped_storage
	Scopes       []scoped_storage
}

type scoped_storage struct {
	Parent            *scoped_storage
	RetVal            int
	VariableAddresses map[string]variables.SymbolTableEntry
	Functions         map[string]*functions.FunctionDefinition
	Offset            int
	Instructions      []runtime.Instruction //Instructions from statements/expressions in the local scope.
	Labels            map[string]int
}

func newScopedStorage() scoped_storage {
	return scoped_storage{
		Functions:         make(map[string]*functions.FunctionDefinition),
		VariableAddresses: make(map[string]variables.SymbolTableEntry),
		Labels:            make(map[string]int),
	}
}

func NewStorage() Storage {
	storage := Storage{}
	storage.Scopes = append(storage.Scopes, newScopedStorage())
	storage.CurrentScope = &storage.Scopes[0]
	return storage
}

func (s *Storage) NewFunction(name string, definition functions.FunctionDefinition) {
	// Make the function visible in the function's parent scope
	s.CurrentScope.Parent.Functions[name] = &definition

	fmt.Println("Declared new function", name, definition)
}

func (s *Storage) NewFunctionScope(definition functions.FunctionDefinition) {
	s.NewScope()

	// Create variable entries for the argument. They are placed first in the function's symbol table
	for _, arg := range definition.ArgumentList {
		s.NewVariable(arg.Type, arg.Identifier)
	}
}

func (s *Storage) NewScope() *scoped_storage {
	fmt.Println("Creating scope")
	s.Scopes = append(s.Scopes, newScopedStorage())
	s.Scopes[len(s.Scopes)-1].Parent = s.CurrentScope
	s.CurrentScope = &s.Scopes[len(s.Scopes)-1]
	return s.CurrentScope
}

func (s *Storage) DestroyFunctionScope(runTime *runtime.Runtime) (int, int) {
	fmt.Println("Destroying function scope")
	start, end := runTime.LoadInstructions(s.CurrentScope.Instructions)

	// Load labels
	for label, value := range s.CurrentScope.Labels {
		runTime.Labels[label] = value
	}

	s.CurrentScope = s.CurrentScope.Parent

	if s.CurrentScope != nil {
		// Since DestroyScope emits e-s +1 instructions, all labels in the current scope must be modified to reflect this.
		for label := range s.CurrentScope.Labels {
			s.CurrentScope.Labels[label] += end - start + 1
			runTime.Labels[label] += end - start - 1
			fmt.Printf("Updated label %s = %d\n", label, s.CurrentScope.Labels[label])
		}
	}
	return start, end
}

func (s *Storage) DestroyScope() {
	fmt.Println("Destroying regular scope")
	instructions := s.CurrentScope.Instructions
	s.CurrentScope = s.CurrentScope.Parent
	s.CurrentScope.Instructions = append(s.CurrentScope.Instructions, instructions...)
}

func (s *Storage) NewLiteral(vartype variables.Type) variables.Symbol {
	s.CurrentScope.Offset += 1
	return variables.Symbol{Scope: 0, Offset: s.CurrentScope.Offset - 1, Type: vartype}
}

func (s *Storage) NewVariable(vartype variables.Type, name string) (*variables.Symbol, error) {
	_, exists := s.CurrentScope.VariableAddresses[name]
	if exists {
		log.Fatalf("Redeclaration of variable: %s\n", name)
		return nil, fmt.Errorf("redeclaration of variable: %s\n", name)
	}

	addr := s.CurrentScope.Offset
	s.CurrentScope.VariableAddresses[name] = variables.SymbolTableEntry{
		Type:   vartype,
		Offset: addr,
	}
	fmt.Printf("Created variable %s of type %s (addr: %d)\n", vartype.String(), name, addr)

	s.CurrentScope.Offset += 1

	return &variables.Symbol{Scope: 0, Offset: s.CurrentScope.Offset - 1, Type: vartype}, nil
}

func (s *Storage) GetFunction(name string) functions.FunctionDefinition {
	scope := s.CurrentScope
	def, ok := functions.FunctionDefinition{}, false
	scopeOffset := 0
	for {
		var func_ptr *functions.FunctionDefinition
		func_ptr, ok = scope.Functions[name]

		if ok {
			def = *func_ptr
			break
		}
		scope = (*scope).Parent
		if scope == nil {
			break
		}
		scopeOffset += 1
	}

	if !ok {
		log.Fatalf("Could not resolve function name: %s", name)
	}
	return def
}

func (s *Storage) GetVarAddr(name string) variables.Symbol {
	scope := s.CurrentScope
	symbol, ok := variables.SymbolTableEntry{}, false
	scopeOffset := 0
	for {
		symbol, ok = (*scope).VariableAddresses[name]
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
		log.Fatalf("Could not resolve variable name: %s", name)
	}
	return variables.Symbol{
		Scope:  scopeOffset,
		Offset: symbol.Offset,
		Type:   symbol.Type,
	}
}

func (s *Storage) LoadInstruction(instruction runtime.Instruction) {
	s.CurrentScope.Instructions = append(s.CurrentScope.Instructions, instruction)
}

func (s *Storage) NewLabel(label string, address int) {
	s.CurrentScope.Labels[label] = address
	fmt.Printf("Added label %s = %d\n", label, address)
}
