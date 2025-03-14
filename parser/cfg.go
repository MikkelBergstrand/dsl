package parser

import (
	"dsl/structure"
	"dsl/tokens"
	"fmt"
	"sort"
	"strings"
)

type cfg_alternative []tokens.ItemType
type cfg_rule []cfg_alternative

type cfg_pattern struct {
	A tokens.ItemType
	B cfg_alternative
}

type ParserVariable int

type CFG struct {
	_array []cfg_pattern
	_map   map[tokens.ItemType]structure.Set[int]
}

func (cfg *CFG) Rules() []cfg_pattern {
	return cfg._array
}

func NewCFG() CFG {
	return CFG{
		_array: make([]cfg_pattern, 0),
		_map:   make(map[tokens.ItemType]structure.Set[int]),
	}
}

func (cfg *CFG) addRule(A tokens.ItemType, B cfg_alternative) {
	cfg._array = append(cfg._array, cfg_pattern{A: A, B: B})
}

func (cfg *CFG) addRules(A tokens.ItemType, Bs []cfg_alternative) {
	for _, B := range Bs {
		cfg.addRule(A, B)
	}
}

func (cfg *CFG) compile() {
	for idx, alt := range cfg._array {
		_, ok := cfg._map[alt.A]
		if !ok {
			cfg._map[alt.A] = *structure.NewSet[int]()
		}
		cfg._map[alt.A].Add(idx)
	}
}

func (cfg *CFG) updateRule(A tokens.ItemType, B cfg_rule) {
	indices := cfg._map[A].List()

	if len(indices) != len(B) {
		fmt.Printf("Cannot update a rule like this! %d %d", len(indices), len(B))
		panic("!")
	}
	for _, index := range indices {
		cfg._array[index].B = B[index]
	}

}

func (cfg CFG) GetRulesForA(token tokens.ItemType) []cfg_alternative {
	ret := make([]cfg_alternative, 0)
	for _, idx := range cfg._map[token].List() {
		ret = append(ret, cfg._array[idx].B)
	}
	return ret
}

func (cfg CFG) GetRuleIndexesForA(token tokens.ItemType) map[int]struct{} {
	return cfg._map[token].Elements()
}

func (cfg CFG) RuleByIndex(index int) cfg_pattern {
	return cfg._array[index]
}

func CreateCFG() CFG {
	cfg := NewCFG()

	cfg.addRule(tokens.NTGoal, cfg_alternative{tokens.NTStatement})
	cfg.addRule(tokens.NTStatement, cfg_alternative{tokens.NTExpr, tokens.ItemSemicolon})

	cfg.addRule(tokens.NTExpr, cfg_alternative{tokens.NTExpr, tokens.ItemOpPlus, tokens.NTTerm})
	cfg.addRule(tokens.NTExpr, cfg_alternative{tokens.NTExpr, tokens.ItemOpMinus, tokens.NTTerm})
	cfg.addRule(tokens.NTExpr, cfg_alternative{tokens.NTTerm})

	cfg.addRule(tokens.NTTerm, cfg_alternative{tokens.NTTerm, tokens.ItemOpMult, tokens.NTFactor})
	cfg.addRule(tokens.NTTerm, cfg_alternative{tokens.NTTerm, tokens.ItemOpDiv, tokens.NTFactor})
	cfg.addRule(tokens.NTTerm, cfg_alternative{tokens.NTFactor})

	cfg.addRules(tokens.NTFactor, []cfg_alternative{
		{tokens.ItemParOpen, tokens.NTExpr, tokens.ItemParClosed},
		{tokens.ItemNumber},
		{tokens.ItemIdentifier},
	})

	cfg.addRules(tokens.NTStatement, []cfg_alternative{
		{tokens.ItemKeyInt, tokens.ItemIdentifier, tokens.ItemEquals, tokens.NTExpr},
		//{tokens.ItemIdentifier, tokens.ItemEquals, tokens.NTExpr},
	})

	cfg.addRule(tokens.NTStatement, cfg_alternative{tokens.NTStatement, tokens.ItemSemicolon, tokens.NTStatement})

	cfg.compile()

	return cfg
}

func EliminateLeftRecursion(cfg CFG, grammar *tokens.Grammar) CFG {
	keys := make([]tokens.ItemType, len(cfg._map))
	i := 0
	for a := range cfg._map {
		keys[i] = a
		i += 1
	}
	sort.Slice(keys, func(i, j int) bool { return int(keys[i]) < int(keys[j]) })

	new_cfg := NewCFG()
	new_cfg._array = make([]cfg_pattern, len(cfg._array))
	copy(new_cfg._array, cfg._array)

	for i := range keys {
		new_cfg._map[keys[i]] = cfg._map[keys[i]].Copy()
		for s := 0; s < i; s++ {
			for _, rule := range cfg.GetRulesForA(keys[i]) {
				if rule.startsWith(keys[s]) { // Find forms A_i -> A_s T
					//Replace them with forms A_i -> d_i T, where d_i is rules of the form A_i -> d_i
					new_cfg.updateRule(keys[i], rule_append(cfg.GetRulesForA(keys[s]), rule[1:]))
				}
			}
		}

		//Eliminate direct A_i recursion
		var A []cfg_alternative // Rules that follow the recursed variable, e.g. E -> E a
		var B []cfg_alternative // Rules that do not follow the recursed variable, e.g. E -> B
		var A_indices []int
		var B_indices []int
		for _, index := range new_cfg._map[keys[i]].List() {
			alt := new_cfg._array[index]
			if alt.B.startsWith(keys[i]) {
				A = append(A, alt.B[1:])
				A_indices = append(A_indices, index)
			} else {
				B = append(B, alt.B)
				B_indices = append(B_indices, index)
			}
		}

		// We found an option of the form E -> E a. Must recurse.
		if len(A) > 0 {
			// Create a new token from E, E' (E tilde)
			new_token := grammar.NewTokenID(keys[i])

			// Set up so that E' -> (a E' | epsilon)
			A = rule_append(A, []tokens.ItemType{new_token})

			new_cfg._map[new_token] = *structure.NewSet[int]()
			for j, A_index := range A_indices {
				new_cfg._array[A_index] = cfg_pattern{
					A: new_token,
					B: A[j],
				}
				new_cfg._map[keys[i]].Remove(A_index)
				new_cfg._map[new_token].Add(A_index)
			}
			// add the epsilon rule
			new_cfg._array = append(new_cfg._array, cfg_pattern{
				A: new_token,
				B: cfg_alternative{tokens.ItemEpsilon},
			})
			new_cfg._map[new_token].Add(len(new_cfg._array) - 1)

			// Set up such that E -> B E'
			B = rule_append(B, []tokens.ItemType{new_token})
			for j, B_index := range B_indices {
				new_cfg._array[B_index].B = B[j]
			}
		}

	}
	return new_cfg
}

func (rule cfg_rule) String() string {
	var out []string
	for _, alt := range rule {
		out = append(out, alt.String())
	}
	return strings.Join(out, "| ")
}
func (alt cfg_alternative) String() string {
	s := ""
	for _, item := range alt {
		s += fmt.Sprint(item, " ")
	}
	return s
}

func (cfg CFG) String() string {
	s := ""
	for a, b := range cfg._map {
		for _, i := range b.List() {
			s += fmt.Sprintf("%s\t> %s\n", a, cfg._array[i].B)
		}
	}
	return s
}

func (alt cfg_alternative) startsWith(i tokens.ItemType) bool {
	return alt[0] == i
}

func rule_append(rule cfg_rule, to_append []tokens.ItemType) cfg_rule {
	var new_rules cfg_rule
	for i := range rule {
		var new_alt cfg_alternative
		new_alt = append(new_alt, rule[i]...)
		new_alt = append(new_alt, to_append...)
		new_rules = append(new_rules, new_alt)
	}
	return new_rules
}
