package parser

import (
	"dsl/structure"
	"dsl/tokens"
	"fmt"
)

func Start(A tokens.ItemType, B tokens.ItemType, grammar *tokens.Grammar, first FirstSet, follow FirstSet) structure.Set[tokens.ItemType] {
	if !first[B].Contains(tokens.ItemEpsilon) {
		return first[B].Copy()
	}
	ret := first[B].Copy().Remove(tokens.ItemEpsilon).Union(follow[A])

	fmt.Println(first[B], follow[A])
	return ret
}
