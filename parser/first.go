package parser

import (
	"dsl/structure"
	"dsl/tokens"
	"fmt"
)

type FirstSet map[tokens.ItemType]structure.Set[tokens.ItemType]

func (set FirstSet) String() string {
	s := ""
	for k, v := range set {
		s += fmt.Sprintf("%s: ", k)

		for _, val := range v.List() {
			s += fmt.Sprintf("%s ", val)
		}
		s += "\n"
	}
	return s
}
func First(cfg CFG, grammar tokens.Grammar) FirstSet {
	firstSet := make(FirstSet)

	for _, terminal := range grammar.Terminals {
		firstSet[terminal] = *structure.NewSet[tokens.ItemType]()
		firstSet[terminal].Add(terminal)
	}

	firstSet[tokens.ItemEOF] = *structure.NewSet[tokens.ItemType]()
	firstSet[tokens.ItemEOF].Add(tokens.ItemEOF)

	firstSet[tokens.ItemEpsilon] = *structure.NewSet[tokens.ItemType]()
	firstSet[tokens.ItemEpsilon].Add(tokens.ItemEpsilon)

	for _, terminal := range grammar.NonTerminals {
		firstSet[terminal] = *structure.NewSet[tokens.ItemType]()
	}

	changing := true
	for changing {
		changing = false
		for _, rule := range cfg.Rules() {
			lhs := rule.A
			alt := rule.B

			rhs := firstSet[alt[0]].Copy().Remove(tokens.ItemEpsilon)
			trailing := true

			for i := 0; i < len(alt)-1; i++ {
				if firstSet[alt[i]].Contains(tokens.ItemEpsilon) {
					epsilonAlreadyPresent := rhs.Contains(tokens.ItemEpsilon)
					rhs.Union(firstSet[alt[i+1]])
					if !epsilonAlreadyPresent {
						rhs.Remove(tokens.ItemEpsilon)
					}
				} else {
					trailing = false
					break
				}
			}

			if trailing && firstSet[alt[len(alt)-1]].Contains(tokens.ItemEpsilon) {
				rhs.Add(tokens.ItemEpsilon)
			}

			prevCount := firstSet[lhs].Size()
			firstSet[lhs].Union(rhs)

			// Compare size of set before and after union, determine if FIRST is still changing
			if !changing && prevCount != firstSet[lhs].Size() {
				changing = true
			}
		}
	}

	fmt.Println(firstSet)
	return firstSet
}
