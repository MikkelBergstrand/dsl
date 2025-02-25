package parser

import (
	"dsl/structure"
	"dsl/tokens"
)

func Follow(cfg CFG, grammar tokens.Grammar, first FirstSet) FirstSet {
	follow := make(FirstSet)
	for _, term := range grammar.NonTerminals {
		follow[term] = *structure.NewSet[tokens.ItemType]()
	}

	changing := true
	for changing {
		for lhs, rule := range cfg {
			for _, alt := range rule {
				trailer := follow[lhs].Copy()
				for i := len(alt) - 1; i >= 0; i-- {
					if alt[i].IsNonTerminal() {
						follow[alt[i]].Union(trailer)
						if first[alt[i]].Contains(tokens.ItemEpsilon) {
							alreadyContainsEpsilon := trailer.Contains(tokens.ItemEpsilon)
							trailer.Union(first[alt[i]])
							if !alreadyContainsEpsilon {
								trailer.Remove(tokens.ItemEpsilon)
							}
						} else {
							trailer.Union(first[alt[i]])
						}
					} else if alt[i].IsTerminal() {
						trailer = first[alt[i]].Copy()
					}
				}
			}
		}
	}
	return follow
}
