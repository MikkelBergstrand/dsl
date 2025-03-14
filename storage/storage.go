package storage

import "fmt"

type Storage struct {
	Ints    []int
	Strings []string

	IntVars map[string]int
}

func NewStorage() Storage {
	return Storage{
		IntVars: make(map[string]int),
	}
}

func (s *Storage) NewIntLiteral(val int) int {
	s.Ints = append(s.Ints, val)
	return len(s.Ints) - 1
}

func (s *Storage) NewIntVariable(name string) int {
	s.Ints = append(s.Ints, 0)
	s.IntVars[name] = len(s.Ints) - 1
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
	return s.IntVars[name]
}

func (s *Storage) SetInt(address int, value int) {
	s.Ints[address] = value
	fmt.Printf("Set %d to %d\n", address, value)
}
