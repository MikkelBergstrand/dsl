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
	Ints         []int
	Strings      []string
	Functions    map[string]functions.Function
}

func generateGlobalFunctions() (ret map[string]functions.Function) {
	ret = make(map[string]functions.Function)

	ret["echo"] = functions.Function{
		ArgumentList: []functions.Argument{
			{Type: variables.INT},
		},
	}

	return ret
}

type scoped_storage struct {
	Parent  *scoped_storage
	IntVars map[string]int
}

func newScopedStorage() scoped_storage {
	return scoped_storage{
		IntVars: make(map[string]int),
	}
}

func NewStorage() Storage {
	storage := Storage{}
	storage.Scopes = append(storage.Scopes, newScopedStorage())
	storage.CurrentScope = &storage.Scopes[0]
	storage.Functions = generateGlobalFunctions()
	return storage
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

func (s *Storage) NewIntLiteral(val int) int {
	s.Ints = append(s.Ints, val)
	return len(s.Ints) - 1
}

func (s *Storage) NewIntVariable(name string) int {
	_, exists := s.CurrentScope.IntVars[name]
	if exists {
		log.Fatalf("Redeclaration of variable: %s\n", name)
	}

	s.Ints = append(s.Ints, 0)
	s.CurrentScope.IntVars[name] = len(s.Ints) - 1
	fmt.Printf("Created int %s (addr: %d)\n", name, len(s.Ints)-1)
	return len(s.Ints) - 1
}

func (s *Storage) NewInt() int {
	s.Ints = append(s.Ints, 0)
	return len(s.Ints) - 1
}

func (s *Storage) GetInt(address int) int {
	return s.Ints[address]
}

func (s *Storage) GetIntVarAddr(name string) int {
	scope := s.CurrentScope
	addr, ok := -1, false
	for !ok && scope != nil {
		addr, ok = (*scope).IntVars[name]
		scope = (*scope).Parent
	}

	if !ok {
		log.Fatalf("Could not resolve integer variable name: %s", name)
	}
	return addr
}

func (s *Storage) SetInt(address int, value int) {
	s.Ints[address] = value
	fmt.Printf("Set %d to %d\n", address, value)
}
