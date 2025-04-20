package tokens

import (
	"fmt"
	"log"
	"strconv"
)

type ItemType int

// Constant types across all grammars
const (
	ItemError   ItemType = 10000
	ItemEpsilon ItemType = 10001
	ItemEOF     ItemType = 10002
)

const TERMINAL_START = 1
const NONTERMINAL_START = 1001

const (
	//TERMINALS
	ItemNumber ItemType = iota + TERMINAL_START
	ItemOpPlus
	ItemOpMinus
	ItemOpMult
	ItemOpDiv
	ItemOpMod
	ItemIdentifier // starts with [a-zA-Z_], followed by [a-zA-Z0-9_]
	ItemParOpen
	ItemParClosed
	ItemSemicolon
	ItemText
	ItemKeyInt
	ItemKeyBool
	ItemFalse
	ItemTrue
	ItemEquals
	ItemScopeOpen
	ItemScopeClose
	ItemComma
	ItemBoolAnd
	ItemBoolOr
	ItemBoolNot
	ItemBoolLess
	ItemBoolLessOrEqual
	ItemBoolEqual
	ItemBoolGreaterOrEqual
	ItemBoolGreater
	ItemBoolNotEqual
	ItemFunction
	ItemIf
	ItemElse
	ItemReturn
	TERMINALS_LENGTH
)

const (
	//NON-Terminals
	NTGoal ItemType = iota + NONTERMINAL_START
	NTStatement
	NTStatementList
	NTExpr
	NTTerm
	NTFactor
	NTScopeBegin
	NTScopeClose
	NTFunctionCall
	NTArgument
	NTArgList
	NTNExpr
	NTAndTerm
	NTNotTerm
	NTRelExpr
	NTRels
	NTArgumentDeclaration
	NTArgumentDeclarationList
	NTVarType
	NTFunctionOpen
	NTFunctionClose
	NTFunctionDefinition
	NTFunctionBody
	NTIfStatement
	NTLabelledScopeBegin
	NTLabelledScopeClose
	NTIfHeader
	NTWithElse
	NTEndConditionalScope
	NTBeginElseIf
	NTTypeList
	NTImplicitFunctionDefinition
	NONTERMINALS_LENGTH
)

type Grammar struct {
	Terminals    []ItemType
	NonTerminals []ItemType
	StartSymbol  ItemType
}

func NewGrammar(startSymbol ItemType) Grammar {
	grammar := Grammar{
		StartSymbol: startSymbol,
	}

	for i := TERMINAL_START; i < int(TERMINALS_LENGTH); i++ {
		grammar.Terminals = append(grammar.Terminals, ItemType(i))
	}
	for i := NONTERMINAL_START; i < int(NONTERMINALS_LENGTH); i++ {
		grammar.NonTerminals = append(grammar.NonTerminals, ItemType(i))
	}
	return grammar
}

func (grammar *Grammar) MapToArrayindex(item ItemType) int {
	if item.IsTerminal() {
		return int(item) - 1
	} else if item.IsNonTerminal() {
		return int(item) - 1002
	} else if item == ItemEOF {
		return len(grammar.Terminals)
	}

	log.Fatalf("Attempted to store %s in an array!", item.String())
	panic("See log")
}

func (grammar *Grammar) NewCategoryID(item ItemType) ItemType {
	temp := (ItemType)(grammar.NonTerminals[len(grammar.NonTerminals)-1] + 1)
	fmt.Println(temp)
	_new_str[temp] = item.String() + "'"
	grammar.NonTerminals = append(grammar.NonTerminals, temp)
	return temp
}

type Token struct {
	Category ItemType
	Lexeme   string
}

func (l Token) String() string {
	switch l.Category {
	case ItemEOF:
		return "EOF"
	case ItemError:
		return l.Lexeme
	}

	if len(l.Lexeme) > 50 {
		return fmt.Sprintf("%d %s", int(l.Category), l.Lexeme[:50])
	}

	return fmt.Sprintf("%d %q", int(l.Category), l.Lexeme)
}

func (l ItemType) IsTerminal() bool {
	return l >= 0 && l < 1000
}

func (l ItemType) IsNonTerminal() bool {
	return l >= 1000 && l < 10000
}

var _new_str map[ItemType]string = make(map[ItemType]string)

func (i ItemType) String() string {
	return strconv.Itoa(int(i))
}
