package parser

import (
	"dsl/structure"
	"dsl/tokens"
	"fmt"
)

type LLTable [][]cfg_pattern

func MakeLLTable(grammar tokens.Grammar, cfg CFG, firstSet FirstSet, followset FirstSet) LLTable {
	table := make([][]cfg_pattern, len(grammar.NonTerminals))

	terminals := make([]tokens.ItemType, len(grammar.Terminals)+1)
	copy(terminals, grammar.Terminals)
	terminals[len(terminals)-1] = tokens.ItemEOF
	fmt.Println(terminals)

	for i := range grammar.NonTerminals {
		table[i] = make([]cfg_pattern, len(terminals))
		for j := range terminals {
			table[i][j] = cfg_pattern{
				A: tokens.ItemError,
				B: cfg_alternative{},
			}
		}
	}

	for _, rule := range cfg.Rules() {
		A := rule.A
		B := rule.B
		startSet := Start(A, B[0], &grammar, firstSet, followset)
		fmt.Println(A, B, startSet)
		for _, w := range startSet.List() {
			if w.IsTerminal() {
				table[A-1001][w-1] = cfg_pattern{
					A: A,
					B: B,
				}
			}
		}

		if startSet.Contains(tokens.ItemEOF) {
			table[A-1001][len(terminals)-1] = cfg_pattern{
				A: A,
				B: B,
			}
		}
	}
	for i := range table {
		fmt.Println(tokens.ItemType(i+1001), table[i])
	}
	return table
}

type LLParser struct {
	Words   chan tokens.Lexeme
	LLTable LLTable
	Grammar tokens.Grammar
}

func NewParser(grammar tokens.Grammar, ll_table LLTable) LLParser {
	return LLParser{
		Words:   make(chan tokens.Lexeme),
		LLTable: ll_table,
		Grammar: grammar,
	}
}

func (p *LLParser) NextWord() tokens.Lexeme {
	return <-p.Words
}

func LLParse(p *LLParser) bool {
	word := p.NextWord()
	stack := structure.NewStack[tokens.ItemType]()
	stack = append(stack, tokens.ItemEOF)
	stack = append(stack, p.Grammar.StartSymbol)

	for {
		focus := stack.Peek()
		if focus == tokens.ItemEOF && word.ItemType == tokens.ItemEOF {
			return true
		} else if focus.IsTerminal() || focus == tokens.ItemEOF {
			if focus == word.ItemType {
				stack.Pop()
				word = p.NextWord()
			} else {
				return false
			}
		} else { // focus is nonterminal
			if p.LLTable[focus-1001][word.ItemType-1].A != tokens.ItemError {
				fmt.Println(stack, focus, word.ItemType)
				stack.Pop()
				B := p.LLTable[focus-1001][word.ItemType-1].B
				for i := len(B) - 1; i >= 0; i-- {
					if B[i] != tokens.ItemEpsilon {
						stack.Push(B[i])
					}
				}
			} else {
				return false
			}
		}
	}
}
