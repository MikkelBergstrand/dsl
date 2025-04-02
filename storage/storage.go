package storage

import (
	"dsl/functions"
	"dsl/variables"
	"fmt"
	"log"
)

type Storage struct {
	CurrentScope *scoped_storage
	Scopes       []scoped_storage
	Functions    map[string]functions.FunctionDefinition
}

type scoped_storage struct {
	Parent            *scoped_storage
	RetVal            int
	VariableAddresses map[string]variables.SymbolTableEntry
	Offset            int
}

func newScopedStorage() scoped_storage {
	return scoped_storage{
		VariableAddresses: make(map[string]variables.SymbolTableEntry),
	}
}

func NewStorage() Storage {
	storage := Storage{}
	storage.Scopes = append(storage.Scopes, newScopedStorage())
	storage.CurrentScope = &storage.Scopes[0]
	storage.Functions = make(map[string]functions.FunctionDefinition)
	return storage
}

func (s *Storage) NewFunction(name string, definition functions.FunctionDefinition) {
	s.Functions[name] = definition
	s.NewScope()
}

func (s *Storage) NewScope() {
	fmt.Println("Creating scope")
	s.Scopes = append(s.Scopes, newScopedStorage())
	s.Scopes[len(s.Scopes)-1].Parent = s.CurrentScope
	s.CurrentScope = &s.Scopes[len(s.Scopes)-1]
}

func (s *Storage) DestroyScope() {
	fmt.Println("Destroying scope")
	s.CurrentScope = s.CurrentScope.Parent
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

func (s *Storage) GetVarAddr(name string) variables.Symbol {
	scope := s.CurrentScope
	fmt.Println(scope)
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
