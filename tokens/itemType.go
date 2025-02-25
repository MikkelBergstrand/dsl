package tokens

import (
	"fmt"
)

type ItemType int

const (
	//TERMINALS
	ItemError ItemType = iota
	ItemText
	ItemNumber
	ItemOpPlus
	ItemOpMinus
	ItemOpMult
	ItemOpDiv
	ItemIdentifier // starts with [a-zA-Z_], followed by [a-zA-Z0-9_]
	ItemSemicolon
	ItemParOpen
	ItemParClosed

	ItemA
	ItemB

	//NON-TERMINALS
	NT_BEGIN ItemType = iota + 1000
	NTExpr
	NTTerm
	NTFactor

	NT_A1
	NT_B2
	NT_GOAL

	NT_TERM

	ItemEpsilon = -1
	ItemEOF     = -2
)

var NT_Count = NT_TERM

type Grammar struct {
	Terminals    []ItemType
	NonTerminals []ItemType
}

func NewTokenID(item ItemType) ItemType {
	temp := NT_Count
	NT_Count = NT_Count + 1
	_new_str[temp] = item.String() + "'"
	return temp
}

type Lexeme struct {
	ItemType ItemType
	Value    string
}

func (l Lexeme) String() string {
	switch l.ItemType {
	case ItemEOF:
		return "EOF"
	case ItemError:
		return l.Value
	}

	if len(l.Value) > 50 {
		return fmt.Sprintf("%d %s", int(l.ItemType), l.Value[:50])
	}

	return fmt.Sprintf("%d %q", int(l.ItemType), l.Value)
}

func (l ItemType) IsTerminal() bool {
	return l >= 0 && l < 1000
}

func (l ItemType) IsNonTerminal() bool {
	return l >= 1000
}

var _new_str map[ItemType]string = make(map[ItemType]string)

func (i ItemType) String() string {

	val, ok := _new_str[i]
	if ok {
		return val
	}

	switch i {
	case NT_GOAL:
		return "Goal"
	case NT_A1:
		return "A1"
	case NT_B2:
		return "B2"
	case NTExpr:
		return "Expr"
	case NTTerm:
		return "Term"
	case NTFactor:
		return "Factor"
	case ItemOpDiv:
		return "/"
	case ItemOpMinus:
		return "-"
	case ItemOpPlus:
		return "+"
	case ItemOpMult:
		return "*"
	case ItemA:
		return "a"
	case ItemB:
		return "b"
	case ItemEpsilon:
		return "eps"
	case ItemParClosed:
		return ")"
	case ItemParOpen:
		return "("
	case ItemIdentifier:
		return "name"
	case ItemNumber:
		return "num"
	case ItemEOF:
		return "eof"
	}
	return string(int(i))
}
