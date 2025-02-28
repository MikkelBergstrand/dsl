package parser

import (
	"dsl/structure"
	"dsl/tokens"
)

func Start(A tokens.ItemType, B tokens.ItemType, grammar *tokens.Grammar, first *FirstSet, follow *FirstSet) structure.Set[tokens.ItemType] {
	if (*first)[B].Contains(tokens.ItemEpsilon) {
		return (*first)[B].Copy()
	}
	ret := (*first)[B].Copy()
	ret.Remove(tokens.ItemEpsilon)
	ret.Union((*follow)[A])
	return ret
}
