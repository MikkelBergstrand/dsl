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
	VariableAddresses map[string]int
	Offset            int
}

func newScopedStorage() scoped_storage {
	return scoped_storage{
		VariableAddresses: make(map[string]int),
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

func (s *Storage) NewIntLiteral() variables.Symbol {
	s.CurrentScope.Offset += 1
	return variables.Symbol{Scope: 0, Offset: s.CurrentScope.Offset - 1}
}

func (s *Storage) NewIntVariable(name string) variables.Symbol {
	_, exists := s.CurrentScope.VariableAddresses[name]
	if exists {
		log.Fatalf("Redeclaration of variable: %s\n", name)
	}

	addr := s.CurrentScope.Offset
	s.CurrentScope.VariableAddresses[name] = addr
	fmt.Printf("Created int %s (rel.adr.: %d, addr: %d)\n", name, s.CurrentScope.Offset-1, addr)

	s.CurrentScope.Offset += 1
	return variables.Symbol{Scope: 0, Offset: s.CurrentScope.Offset - 1}
}

func (s *Storage) GetVarAddr(name string) variables.Symbol {
	scope := s.CurrentScope
	fmt.Println(scope)
	addr, ok := -1, false
	scopeOffset := 0
	for {
		addr, ok = (*scope).VariableAddresses[name]
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
		log.Fatalf("Could not resolve integer variable name: %s", name)
	}
	return variables.Symbol{
		Scope:  scopeOffset,
		Offset: addr,
	}
}
