package parser

import (
	"dsl/tokens"
	"fmt"
)

func LLTable(grammar tokens.Grammar, cfg CFG, firstSet FirstSet, followset FirstSet) {
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

	for A, Bs := range cfg {
		for _, B := range Bs {
			startSet := Start(A, B[0], &grammar, &firstSet, &followset)
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
	}
	for i := range table {
		fmt.Println(tokens.ItemType(i+1001), table[i])
	}
}
