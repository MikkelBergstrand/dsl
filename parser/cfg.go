package parser

import (
	"dsl/structure"
	"dsl/tokens"
	"fmt"
	"sort"
	"strings"
)

// An alternative is the right-hand side of a single production, e.g. the "B" in A -> B
type cfg_alternative []tokens.ItemType

// A set of right-hand side productions, typically grouped together by a common left-hand side symbol.
type cfg_rule []cfg_alternative

type cfg_pattern struct {
	A tokens.ItemType
	B cfg_alternative
}

// A CFG (context-free grammar) holds all valid productions (A -> B) in a grammar.
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

// Add a rule A -> B
func (cfg *CFG) addRule(A tokens.ItemType, B cfg_alternative) {
	cfg._array = append(cfg._array, cfg_pattern{A: A, B: B})
}

// Add multiple rules of form A -> B1, A -> B2, ... where B1, B2, .. are elements of Bs
func (cfg *CFG) addRules(A tokens.ItemType, Bs []cfg_alternative) {
	for _, B := range Bs {
		cfg.addRule(A, B)
	}
}

// Compiles all the rules by creating a map with the left-hand side of the production
// as the key, and a set as the value corresponding to all rules with this left-hand side
// The map significantly speeds up parsing by quickly retrieving rules from a left-hand side symbol.
// Must be called when done adding rules!
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

	cfg.addRule(tokens.NTGoal, cfg_alternative{tokens.NTStatementList})
	cfg.addRule(tokens.NTStatementList, cfg_alternative{tokens.NTStatement})
	cfg.addRule(tokens.NTStatementList, cfg_alternative{tokens.NTStatement, tokens.NTStatementList})

	cfg.addRule(tokens.NTNExpr, cfg_alternative{tokens.NTNExpr, tokens.ItemOpPlus, tokens.NTTerm})
	cfg.addRule(tokens.NTNExpr, cfg_alternative{tokens.NTNExpr, tokens.ItemOpMinus, tokens.NTTerm})
	cfg.addRule(tokens.NTNExpr, cfg_alternative{tokens.NTTerm})

	cfg.addRule(tokens.NTTerm, cfg_alternative{tokens.NTTerm, tokens.ItemOpMult, tokens.NTFactor})
	cfg.addRule(tokens.NTTerm, cfg_alternative{tokens.NTTerm, tokens.ItemOpDiv, tokens.NTFactor})
	cfg.addRule(tokens.NTTerm, cfg_alternative{tokens.NTFactor})

	cfg.addRules(tokens.NTFactor, []cfg_alternative{
		{tokens.ItemParOpen, tokens.NTExpr, tokens.ItemParClosed},
		{tokens.ItemNumber},
		{tokens.ItemIdentifier},
	})

	cfg.addRules(tokens.NTStatement, []cfg_alternative{
		{tokens.ItemKeyInt, tokens.ItemIdentifier, tokens.ItemEquals, tokens.NTNExpr, tokens.ItemSemicolon},
		{tokens.ItemIdentifier, tokens.ItemEquals, tokens.NTExpr, tokens.ItemSemicolon},
		{tokens.NTExpr, tokens.ItemSemicolon},
	})

	cfg.addRule(tokens.NTStatement, cfg_alternative{tokens.NTScopeBegin, tokens.NTStatement, tokens.NTScopeClose})

	cfg.addRule(tokens.NTScopeBegin, cfg_alternative{tokens.ItemScopeOpen})
	cfg.addRule(tokens.NTScopeClose, cfg_alternative{tokens.ItemScopeClose})

	cfg.addRule(tokens.NTFactor, cfg_alternative{tokens.NTFunction})
	cfg.addRule(tokens.NTFunction, cfg_alternative{tokens.ItemIdentifier, tokens.ItemParOpen, tokens.NTArgList, tokens.ItemParClosed})
	cfg.addRule(tokens.NTArgList, cfg_alternative{tokens.NTArgument, tokens.ItemComma, tokens.NTArgList})
	cfg.addRule(tokens.NTArgList, cfg_alternative{tokens.NTArgument})

	cfg.addRule(tokens.NTArgument, cfg_alternative{tokens.NTExpr})

	cfg.addRules(tokens.NTExpr, []cfg_alternative{
		{tokens.NTExpr, tokens.ItemBoolOr, tokens.NTAndTerm},
		{tokens.NTAndTerm},
	})
	cfg.addRules(tokens.NTAndTerm, []cfg_alternative{
		{tokens.NTAndTerm, tokens.ItemBoolAnd, tokens.NTNotTerm},
		{tokens.NTNotTerm},
	})
	cfg.addRules(tokens.NTNotTerm, []cfg_alternative{
		{tokens.ItemBoolNot, tokens.NTRelExpr},
		{tokens.NTRelExpr},
	})
	cfg.addRules(tokens.NTRelExpr, []cfg_alternative{
		{tokens.NTNExpr, tokens.NTRels, tokens.NTNExpr},
		{tokens.NTNExpr},
	})

	cfg.addRules(tokens.NTRels, []cfg_alternative{
		{tokens.ItemBoolEqual},
		{tokens.ItemBoolNotEqual},
		{tokens.ItemBoolLess},
		{tokens.ItemBoolLessOrEqual},
		{tokens.ItemBoolGreater},
		{tokens.ItemBoolGreaterOrEqual},
	})

	cfg.addRules(tokens.NTFactor, []cfg_alternative{
		{tokens.ItemFalse},
		{tokens.ItemTrue},
	})

	cfg.addRule(tokens.NTStatement,
		cfg_alternative{tokens.ItemKeyBool, tokens.ItemIdentifier, tokens.ItemEquals, tokens.NTExpr, tokens.ItemSemicolon})

	// Note that we dont use NTScopeOpen here, because other parts of the function need to create the scope for us
	// Closing is done using a special non-terminal, however.
	cfg.addRule(tokens.NTStatement,
		cfg_alternative{tokens.ItemFunction, tokens.NTFunctionDefinition, tokens.NTFunctionBody}) //40

	cfg.addRule(tokens.NTFunctionDefinition,
		cfg_alternative{tokens.ItemIdentifier, tokens.ItemParOpen, tokens.NTArgumentDeclarationList, tokens.ItemParClosed, tokens.NTVarType}) //41

	cfg.addRules(tokens.NTArgumentDeclarationList, []cfg_alternative{
		{tokens.NTArgumentDeclaration, tokens.ItemComma, tokens.NTArgumentDeclarationList}, //42
		{tokens.NTArgumentDeclaration}, //43
	})

	cfg.addRules(tokens.NTArgumentDeclaration, []cfg_alternative{
		{tokens.NTVarType, tokens.ItemIdentifier}, //44
	})

	cfg.addRules(tokens.NTVarType, []cfg_alternative{
		{tokens.ItemKeyBool}, //45
		{tokens.ItemKeyInt},  //46
	})

	cfg.addRule(tokens.NTFunctionClose, cfg_alternative{tokens.ItemScopeClose}) // 47

	cfg.addRule(tokens.NTFunctionBody, cfg_alternative{tokens.NTFunctionOpen, tokens.NTStatementList, tokens.NTFunctionClose}) // 48

	cfg.addRule(tokens.NTStatement, cfg_alternative{tokens.NTIfHeader, // 49
		tokens.NTScopeBegin,
		tokens.NTStatementList,
		tokens.NTLabelledScopeClose,
	})

	cfg.addRule(tokens.NTLabelledScopeBegin, cfg_alternative{tokens.ItemScopeOpen})  // 50
	cfg.addRule(tokens.NTLabelledScopeClose, cfg_alternative{tokens.ItemScopeClose}) // 51
	cfg.addRule(tokens.NTFunctionOpen, cfg_alternative{tokens.ItemScopeOpen})        // 52

	cfg.addRule(tokens.NTIfHeader, cfg_alternative{tokens.ItemIf, tokens.NTExpr}) // 53

	cfg.addRule(tokens.NTStatement, cfg_alternative{tokens.NTIfHeader, tokens.NTScopeBegin, tokens.NTStatementList, tokens.NTEndConditionalScope,
		tokens.NTWithElse}) // 54

	cfg.addRule(tokens.NTWithElse, cfg_alternative{tokens.NTBeginElseIf, tokens.NTIfHeader, tokens.NTScopeBegin, tokens.NTStatementList, tokens.NTEndConditionalScope,
		tokens.NTWithElse}) // 55
	cfg.addRule(tokens.NTWithElse,
		cfg_alternative{tokens.NTBeginElseIf, tokens.NTIfHeader, tokens.NTScopeBegin, tokens.NTStatementList, tokens.NTLabelledScopeClose}) // 56
	cfg.addRule(tokens.NTWithElse, cfg_alternative{tokens.ItemElse, tokens.NTLabelledScopeBegin, tokens.NTStatementList, tokens.NTLabelledScopeClose}) // 57
	cfg.addRule(tokens.NTEndConditionalScope, cfg_alternative{tokens.NTScopeClose})                                                                    // 58
	cfg.addRule(tokens.NTBeginElseIf, cfg_alternative{tokens.ItemElse})                                                                                //59
	cfg.addRule(tokens.NTTerm, cfg_alternative{tokens.NTTerm, tokens.ItemOpMod, tokens.NTFactor}) // 60
	fmt.Println("Num rules: ", len(cfg._array))
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
			new_token := grammar.NewCategoryID(keys[i])

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
