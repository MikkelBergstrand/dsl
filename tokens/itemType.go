package tokens

import (
	"fmt"
	"strconv"
)

type ItemType int

// Constant types across all grammars
const (
	ItemError   ItemType = 10000
	ItemEpsilon ItemType = 10001
	ItemEOF     ItemType = 10002
)

const (
	//TERMINALS
	ItemText ItemType = iota + 1
	ItemNumber
	ItemOpPlus
	ItemOpMinus
	ItemOpMult
	ItemOpDiv
	ItemIdentifier // starts with [a-zA-Z_], followed by [a-zA-Z0-9_]
	ItemParOpen
	ItemParClosed
	ItemSemicolon
)

const (
	//NON-Terminals
	NTGoal ItemType = iota + 1001
	NTExpr
	NTTerm
	NTFactor
)

type Grammar struct {
	Terminals    []ItemType
	NonTerminals []ItemType
	StartSymbol  ItemType
}

func (grammar *Grammar) NewTokenID(item ItemType) ItemType {
	temp := (ItemType)(grammar.NonTerminals[len(grammar.NonTerminals)-1] + 1)
	fmt.Println(temp)
	_new_str[temp] = item.String() + "'"
	grammar.NonTerminals = append(grammar.NonTerminals, temp)
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
	return l >= 1000 && l < 10000
}

var _new_str map[ItemType]string = make(map[ItemType]string)

func (i ItemType) String() string {

	val, ok := _new_str[i]
	if ok {
		return val
	}

	switch i {
	case NTGoal:
		return "Goal"
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
	case ItemEpsilon:
		return "eps"
	case ItemIdentifier:
		return "name"
	case ItemNumber:
		return "num"
	case ItemEOF:
		return "eof"
	case ItemParOpen:
		return "("
	case ItemParClosed:
		return ")"
	case ItemError:
		return "err"
	}
	return strconv.Itoa(int(i))
}
