package parser

import (
	"dsl/functions"
	"dsl/runtime"
	"dsl/storage"
	"dsl/variables"
	"fmt"
	"log"
	"strconv"
)

type condition_tree_entry struct {
	start_label string
	end         *runtime.InstructionLabelPair // Instruction at the end of a conditonal block. Can be nil, if no block follows it.
	jmp         *runtime.InstrJmpIf           // Instruction that starts the conditional block, of type InstrJmpIf. Can again be nil, for else statement.
}

type List[T any] struct {
	First  T
	Second *List[T]
}

func (list List[T]) Iterate() (ret []T) {
	ret = append(ret, list.First)

	node := list.Second
	for node != nil {
		ret = append(ret, node.First)
		node = node.Second
	}
	return ret
}

// Convert string to integer. Should not fail!
func intval(s string) int {
	intval, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return intval
}

func validateBooleanArithmetic(a variables.Symbol, b variables.Symbol, op runtime.BooleanOperator) error {
	if a.Type != b.Type || !op.IsValidFor(a.Type) {
		return fmt.Errorf("invalid type comparison of %s and %s", a.Type, b.Type)
	}
	return nil
}

func booleanArithmetic(words []any, s *storage.Storage, op runtime.BooleanOperator) variables.Symbol {
	a := words[0].(variables.Symbol)
	b := words[2].(variables.Symbol)

	err := validateBooleanArithmetic(a, b, op)
	if err != nil {
		log.Fatalln(err.Error())
	}

	newaddr := s.NewLiteral(variables.BOOL)
	if a.Type == variables.BOOL {
		s.LoadInstruction(&runtime.InstrCompareBool{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   newaddr,
			Operator: op,
		})
	} else if a.Type == variables.INT {
		s.LoadInstruction(&runtime.InstrCompareInt{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   newaddr,
			Operator: op,
		})
	}
	return newaddr
}

func DoActions(rule_id int, words []any, storage *storage.Storage, r *runtime.Runtime) any {
	fmt.Println(rule_id, words)
	switch rule_id {
	case 3:
		new_addr := storage.NewLiteral(variables.INT)
		storage.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.ADD,
		})
		return new_addr
	case 4:
		new_addr := storage.NewLiteral(variables.INT)
		storage.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.SUB,
		})
		return new_addr
	case 6:
		new_addr := storage.NewLiteral(variables.INT)
		storage.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.MULT,
		})
		return new_addr
	case 7:
		new_addr := storage.NewLiteral(variables.INT)
		storage.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.DIV,
		})
		return new_addr
	case 10: //New integer literal
		addr := storage.NewLiteral(variables.INT)
		storage.LoadInstruction(&runtime.InstrLoadImmediate{
			Dest:  addr,
			Value: intval(words[0].(string)),
		})
		return addr
	case 11:
		return storage.GetVarAddr(words[0].(string))
	case 12: // New integer, eg. int a = 3
		addr, err := storage.NewVariable(variables.INT, words[1].(string))
		if err != nil {
			log.Fatal(err)
		}

		if words[3].(variables.Symbol).Type != variables.INT {
			log.Fatalf("Invalid type assignment: expected int, got %s", words[3].(variables.Symbol).Type.String())
		}

		storage.LoadInstruction(&runtime.InstrAssign{
			Source: words[3].(variables.Symbol),
			Dest:   *addr,
		})
		return *addr
	case 13: // Reassignment of integer, e.g. a = 3
		addr := storage.GetVarAddr(words[0].(string))

		storage.LoadInstruction(&runtime.InstrAssign{
			Source: words[2].(variables.Symbol),
			Dest:   addr,
		})
		return addr
	case 16: // Declare scope
		storage.LoadInstruction(&runtime.InstrBeginScope{
			AddressStart: storage.CurrentScope.Offset,
		})
		storage.NewScope()
	case 17: // End scope
		storage.LoadInstruction(&runtime.InstrEndScope{})
		storage.DestroyScope()
	case 19: // call function e.g. echo ( 0 )
		arg_list := (words[2].(List[variables.Symbol])).Iterate()
		fn := storage.GetFunction(words[0].(string))

		if !fn.ValidateArgumentList(arg_list) {
			log.Fatalf("Argument list to function %s invalid\n", words[0].(string))
		}

		// Bit hacky, but from the perspective of the new function,
		// the arguments are located in the above scope. Hence, we must
		// increment the scope to account for this.
		passed_arg_list := make([]variables.Symbol, len(arg_list))
		for i := range passed_arg_list {
			passed_arg_list[i] = arg_list[i]
			passed_arg_list[i].Scope += 1
		}

		storage.LoadInstruction(&runtime.InstrCallFunction{
			AddressStart:  storage.CurrentScope.Offset,
			PreludeLength: len(fn.ArgumentList) + 2,
		})

		for i := range fn.ArgumentList {
			storage.LoadInstruction(&runtime.InstrAssign{
				Source: passed_arg_list[i],
				Dest:   variables.Symbol{Scope: 0, Offset: i}, // For simplicity, parameter #i is always stored in the scope with offset i
			})
		}

		storage.LoadInstruction(&runtime.InstrJmp{
			Label: words[0].(string),
		})

	case 20: //argument list construction, input is "symbol , List"
		second := words[2].(List[variables.Symbol])
		return List[variables.Symbol]{
			First:  words[0].(variables.Symbol),
			Second: &second,
		}
	case 21: //Initial list item in an argument list
		return List[variables.Symbol]{
			First:  words[0].(variables.Symbol),
			Second: nil}
	case 23: // a | b
		return booleanArithmetic(words, storage, runtime.OR)
	case 25: // a & b
		return booleanArithmetic(words, storage, runtime.AND)
	case 29: // a == b
		switch words[1].(string) {
		case "==":
			return booleanArithmetic(words, storage, runtime.EQUALS)
		case "!=":
			return booleanArithmetic(words, storage, runtime.NOTEQUALS)
		case "<":
			return booleanArithmetic(words, storage, runtime.LESS)
		case "<=":
			return booleanArithmetic(words, storage, runtime.LESSOREQUAL)
		case ">":
			return booleanArithmetic(words, storage, runtime.GREATER)
		case ">=":
			return booleanArithmetic(words, storage, runtime.GREATEROREQUAL)
		default:
			log.Fatalf("Undefined boolean operator.")
		}
	case 37: // false
		addr := storage.NewLiteral(variables.BOOL)
		storage.LoadInstruction(&runtime.InstrLoadImmediate{
			Dest:  addr,
			Value: false,
		})
		return addr
	case 38: // true
		addr := storage.NewLiteral(variables.BOOL)
		storage.LoadInstruction(&runtime.InstrLoadImmediate{
			Dest:  addr,
			Value: true,
		})
		return addr
	case 39: // Declaration boolean
		addr, err := storage.NewVariable(variables.BOOL, words[1].(string))
		if err != nil {
			log.Fatal(err)
		}

		if variables.BOOL != words[3].(variables.Symbol).Type {
			log.Fatalf("type mismatch in assignment of variable '%s'\n", words[1].(string))
		}

		storage.LoadInstruction(&runtime.InstrAssign{
			Source: words[3].(variables.Symbol),
			Dest:   *addr,
		})
		return *addr
	case 41: // Declare new function, format "name ( arglist ) returntype"
		arg_list := words[2].(List[functions.Argument]).Iterate()
		ret_type, err := variables.TypeFromString(words[4].(string))
		if err != nil {
			log.Fatal(err)
		}

		def := functions.FunctionDefinition{
			ArgumentList: arg_list,
			ReturnType:   ret_type,
		}

		storage.NewFunctionScope(def)
		storage.NewFunction(words[0].(string), def)
		storage.NewLabel(words[0].(string))

		return words[0].(string)
	case 42: //Function argument declaration list, second+ element
		second := words[2].(List[functions.Argument])
		return List[functions.Argument]{
			First:  words[0].(functions.Argument),
			Second: &second,
		}
	case 43: //Function argument declaration list, first element
		return List[functions.Argument]{
			First:  words[0].(functions.Argument),
			Second: nil,
		}
	case 44: //Function argument declaration
		_type, err := variables.TypeFromString(words[0].(string))
		if err != nil {
			log.Fatalln(err.Error())
		}

		return functions.Argument{
			Type:       _type,
			Identifier: words[1].(string),
		}
	case 47: // Function scope close
		storage.LoadInstruction(&runtime.InstrExitFunction{})
		storage.DestroyFunctionScope(r)
	case 49: // If statement, NTIfHeader NTLabelledScopeBegin, NTStatementList, NTLabelledScopeClose
		jmpIfInstr := words[0].(*runtime.InstrJmpIf)
		instrEnd := words[3].(*runtime.InstructionLabelPair)
		jmpIfInstr.Label = instrEnd.Label
	case 50: //NTLabelledScopeBegin
		instr := storage.LoadLabeledInstruction(&runtime.InstrBeginScope{
			AddressStart: storage.CurrentScope.Offset,
		}, storage.NewAutoLabel())
		storage.NewScope()
		return instr
	case 51: //NTLabelledScopeClose
		storage.LoadInstruction(&runtime.InstrEndScope{})
		storage.DestroyScope()

		return storage.LoadLabeledInstruction(&runtime.InstrNOP{}, storage.NewAutoLabel())

	case 52: //Open Function
		storage.LoadInstruction(&runtime.InstrNOP{})
	case 53: //NTIfHeader (if Expr)
		condition := words[1].(variables.Symbol)
		if condition.Type != variables.BOOL {
			log.Fatalln("Expected boolean statement in if clause, got", condition.Type)
		}

		instr := storage.LoadInstruction(&runtime.InstrJmpIf{
			Condition: condition,
			Label:     "", // will be set later.
		})
		jmp_instr := instr.Instruction.(*runtime.InstrJmpIf)
		return jmp_instr
	case 54: //If statement + WithElse
		jmpIfInstr := words[0].(*runtime.InstrJmpIf)
		jmpInstr := words[3].(*runtime.InstructionLabelPair)
		tree := words[4].(List[condition_tree_entry])

		// Add the first if clause to the tree.
		final_tree := List[condition_tree_entry]{
			First:  condition_tree_entry{start_label: "", jmp: jmpIfInstr, end: jmpInstr},
			Second: &tree,
		}

		tree_list := final_tree.Iterate()
		fmt.Println(len(tree_list))
		for i := range len(tree_list) - 1 {
			// Make so all JumpIfs (which begins each conditional block) jump to the next condition
			// Except the final one, which escapes the runtime
			tree_list[i].jmp.Label = tree_list[i+1].start_label
			// Make so all Jumps (which ends each condition block) jump to the end of the conditional
			tree_list[i].end.Instruction.(*runtime.InstrJmp).Label = tree_list[len(tree_list)-1].end.Label
		}
	case 55: //WithElse, else if statement, with continuation
		jmp := words[1].(*runtime.InstrJmpIf)
		end := words[4].(*runtime.InstructionLabelPair)
		list := words[5].(List[condition_tree_entry])

		return List[condition_tree_entry]{
			First: condition_tree_entry{
				start_label: words[0].(string),
				jmp:         jmp,
				end:         end,
			},
			Second: &list,
		}
	case 56: //WithElse, else if statement, no continuation
		jmp_if := words[1].(*runtime.InstrJmpIf)
		end := words[4].(*runtime.InstructionLabelPair)

		return List[condition_tree_entry]{
			First: condition_tree_entry{
				start_label: words[0].(string),
				jmp:         jmp_if,
				end:         end,
			},
			Second: nil,
		}

	case 57: //WithElse, else condition
		// Set label to first instruction, as this is not labelled (else condition has no JmpIf clause)
		start := words[1].(*runtime.InstructionLabelPair)
		end := words[3].(*runtime.InstructionLabelPair)

		return List[condition_tree_entry]{

			First: condition_tree_entry{
				start_label: start.Label,
				jmp:         nil,
				end:         end,
			},
			Second: nil,
		}
	case 58: // End conditional statement that is part of a larger conditional statement
		// At end of conditional block, we must jump to skip over the other conditionals
		return storage.LoadInstruction(&runtime.InstrJmp{})
	case 59:
		label := storage.NewAutoLabel()
		storage.NewLabel(label)
		return label
	}

	return words[0]
}
