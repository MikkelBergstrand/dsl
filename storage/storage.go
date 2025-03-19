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
	Variables    []any
	Functions    map[string]functions.FunctionDefinition
}

type scoped_storage struct {
	AddressPointerBase   int
	AddressPointerOffset int
	Parent               *scoped_storage
	RetVal               int
	VariableAddresses    map[string]int
}

func newScopedStorage(addressPointer int) scoped_storage {
	return scoped_storage{
		VariableAddresses:    make(map[string]int),
		AddressPointerBase:   addressPointer,
		AddressPointerOffset: 0,
	}
}

func NewStorage() Storage {
	storage := Storage{}
	storage.Scopes = append(storage.Scopes, newScopedStorage(0))
	storage.CurrentScope = &storage.Scopes[0]
	storage.Functions = make(map[string]functions.FunctionDefinition)
	storage.Variables = make([]any, 1000)
	return storage
}

func (s *Storage) NewScope() {
	fmt.Println("Creating scope")
	s.Scopes = append(s.Scopes, newScopedStorage(s.CurrentScope.AddressPointerBase+s.CurrentScope.Parent.AddressPointerOffset))
	s.Scopes[len(s.Scopes)-1].Parent = s.CurrentScope
	s.CurrentScope = &s.Scopes[len(s.Scopes)-1]
}

func (storage *Storage) AddressFromSymbol(symbol variables.Symbol) int {
	scope := storage.CurrentScope
	for i := 0; i < symbol.Scope; i++ {
		scope = scope.Parent
	}
	return scope.AddressPointerBase + symbol.Offset
}

func (s *Storage) DestroyScope() {
	fmt.Println("Destroying scope")
	s.CurrentScope = s.CurrentScope.Parent
}

func (s *Storage) NewIntLiteral(val int) variables.Symbol {
	s.Variables[s.CurrentScope.AddressPointerBase+s.CurrentScope.AddressPointerOffset] = val
	s.CurrentScope.AddressPointerOffset += 1
	return variables.Symbol{Scope: 0, Offset: s.CurrentScope.AddressPointerOffset}
}

func (s *Storage) NewIntVariable(name string) variables.Symbol {
	_, exists := s.CurrentScope.VariableAddresses[name]
	if exists {
		log.Fatalf("Redeclaration of variable: %s\n", name)
	}

	s.Variables[s.CurrentScope.AddressPointerBase+s.CurrentScope.AddressPointerOffset] = 0
	s.CurrentScope.VariableAddresses[name] = s.CurrentScope.AddressPointerOffset
	s.CurrentScope.AddressPointerOffset += 1
	fmt.Printf("Created int %s (addr: %d)\n", name, len(s.Variables)-1)
	return variables.Symbol{Scope: 0, Offset: s.CurrentScope.AddressPointerOffset - 1}
}

func (s *Storage) GetInt(symbol variables.Symbol) int {
	resolve, ok := s.Variables[s.AddressFromSymbol(symbol)].(int)
	if !ok {
		log.Fatalf("Not an integer")
	}
	return resolve
}

func (s *Storage) GetVarAddr(name string) variables.Symbol {
	scope := s.CurrentScope
	addr, ok := -1, false
	scopeOffset := 0
	for {
		addr, ok = (*scope).VariableAddresses[name]
		if ok {
			break
		}
		scope = (*scope).Parent
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

func (s *Storage) SetInt(symbol variables.Symbol, value int) {
	s.Variables[s.AddressFromSymbol(symbol)] = value
	fmt.Printf("Set %d to %d\n", symbol, value)
}
