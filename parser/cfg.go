package parser

import (
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

type CFG map[tokens.ItemType]cfg_rule

func CreateCFG() CFG {
	cfg := make(CFG)
	cfg[tokens.NTGoal] = cfg_rule{
		cfg_alternative{tokens.NTExpr},
	}
	cfg[tokens.NTExpr] = cfg_rule{
		cfg_alternative{tokens.NTExpr, tokens.ItemOpPlus, tokens.NTTerm},
		cfg_alternative{tokens.NTExpr, tokens.ItemOpMinus, tokens.NTTerm},
		cfg_alternative{tokens.NTTerm},
	}
	cfg[tokens.NTTerm] = cfg_rule{
		cfg_alternative{tokens.NTTerm, tokens.ItemOpMult, tokens.NTFactor},
		cfg_alternative{tokens.NTTerm, tokens.ItemOpDiv, tokens.NTFactor},
		cfg_alternative{tokens.NTFactor},
	}

	cfg[tokens.NTFactor] = cfg_rule{
		cfg_alternative{tokens.ItemParOpen, tokens.NTExpr, tokens.ItemParClosed},
		cfg_alternative{tokens.ItemNumber},
		cfg_alternative{tokens.ItemIdentifier},
	}

	return cfg
}

func EliminateLeftRecursion(cfg CFG, grammar *tokens.Grammar) CFG {
	keys := make([]tokens.ItemType, len(cfg))
	i := 0
	for k := range cfg {
		keys[i] = k
		i += 1
	}
	sort.Slice(keys, func(i, j int) bool { return int(keys[i]) < int(keys[j]) })

	new_cfg := make(CFG)
	for k, v := range cfg {
		new_cfg[k] = v[:]
	}

	for i := range keys {
		new_cfg[keys[i]] = cfg[keys[i]]
		for s := 0; s < i; s++ {
			rule := cfg[keys[i]] // Get the A_i -> ... alternatives
			for _, alt := range rule {
				if alt.startsWith(keys[s]) { // Find forms A_i -> A_s T
					//Replace them with forms A_i -> d_i T, where d_i is rules of the form A_i -> d_i
					new_cfg[keys[i]] = rule_append(cfg[keys[s]], alt[1:])
				} else {
					new_cfg[keys[i]] = cfg[keys[i]]
				}
			}
		}

		//Eliminate direct A_i recursion
		var A []cfg_alternative // Rules that follow the recursed variable, e.g. E -> E a
		var B []cfg_alternative // Rules that do not follow the recursed variable, e.g. E -> B
		for _, alt := range new_cfg[keys[i]] {
			if alt.startsWith(keys[i]) {
				A = append(A, alt[1:])
			} else {
				B = append(B, alt)
			}
		}

		// We found an option of the form E -> E a. Must recurse.
		if len(A) > 0 {
			// Create a new token from E, E' (E tilde)
			new_token := grammar.NewTokenID(keys[i])

			// Set up so that E' -> (a E' | epsilon)
			A = rule_append(A, []tokens.ItemType{new_token})
			A = append(A, cfg_alternative{tokens.ItemEpsilon})
			new_cfg[new_token] = A

			// Set up such that E -> B E'
			B = rule_append(B, []tokens.ItemType{new_token})
			new_cfg[keys[i]] = B
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
	for a, b := range cfg {
		s += fmt.Sprintf("%s\t> %s\n", a, b)
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
