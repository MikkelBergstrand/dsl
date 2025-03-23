package parser

import (
	"dsl/runtime"
	"dsl/storage"
	"dsl/variables"
	"fmt"
	"log"
	"strconv"
)

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

func DoActions(rule_id int, words []any, storage *storage.Storage, r *runtime.Runtime) (any, bool) {
	fmt.Println(rule_id, words)
	switch rule_id {
	case 3:
		new_addr := storage.NewLiteral(variables.INT)
		r.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.ADD,
		})
		return new_addr, false
	case 4:
		new_addr := storage.NewLiteral(variables.INT)
		r.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.SUB,
		})
		return new_addr, false
	case 5:
		return words[0], false
	case 6:
		new_addr := storage.NewLiteral(variables.INT)
		r.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.MULT,
		})
		return new_addr, false
	case 7:
		new_addr := storage.NewLiteral(variables.INT)
		r.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.DIV,
		})
		return new_addr, false
	case 8:
		return words[0], false
	case 9:
		return words[0], false
	case 10: //New integer literal
		addr := storage.NewLiteral(variables.INT)
		r.LoadInstruction(&runtime.InstrLoadImmediate{
			Dest:  addr,
			Value: intval(words[0].(string)),
		})
		return addr, false
	case 11:
		return storage.GetVarAddr(words[0].(string)), false
	case 12: // New integer, eg. int a = 3
		addr := storage.NewVariable(variables.INT, words[1].(string))
		r.LoadInstruction(&runtime.InstrAssign{
			Source: words[3].(variables.Symbol),
			Dest:   addr,
		})
		return addr, false
	case 13: // Reassignment of integer, e.g. a = 3
		addr := storage.GetVarAddr(words[0].(string))
		r.LoadInstruction(&runtime.InstrAssign{
			Source: words[2].(variables.Symbol),
			Dest:   addr,
		})
		return addr, false
	case 16: // Declare scope
		storage.NewScope()
	case 17: // End scope
		storage.DestroyScope()
	case 19: // call function e.g. echo ( 0 )
		arg_list := (words[2].(List[variables.Symbol])).Iterate()
		fn := storage.Functions[words[0].(string)]

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

		r.LoadInstruction(&runtime.InstrCallFunction{
			ArgumentList: arg_list,
			Func:         fn,
			AddressStart: storage.CurrentScope.Offset,
		})

		for i := range fn.ArgumentList {
			r.LoadInstruction(&runtime.InstrAssign{
				Source: passed_arg_list[i],
				Dest:   variables.Symbol{Scope: 0, Offset: i}, // For simplicity, parameter i is always stored in the scope with offset i
			})
		}

		r.LoadInstruction(&runtime.InstrJmp{
			NewPC: fn.InstructionPointer,
		})

	case 20: //argument list construction, input is "symbol , List"
		second := words[2].(List[variables.Symbol])
		return List[variables.Symbol]{
			First:  words[0].(variables.Symbol),
			Second: &second,
		}, false
	case 21:
		return List[variables.Symbol]{
			First:  words[0].(variables.Symbol),
			Second: nil}, false
	case 22:
		return words[0], false
	case 24:
		return words[0], false
	case 26:
		return words[0], false
	case 28:
		return words[0], false
	case 30: // Evaluate
		return words[0], false
	case 37:
		addr := storage.NewLiteral(variables.BOOL)
		r.LoadInstruction(&runtime.InstrLoadImmediate{
			Dest:  addr,
			Value: false,
		})
		return addr, false
	case 38:
		addr := storage.NewLiteral(variables.BOOL)
		r.LoadInstruction(&runtime.InstrLoadImmediate{
			Dest:  addr,
			Value: true,
		})
		return addr, false
	case 39: // Declaration boolean
		addr := storage.NewVariable(variables.BOOL, words[1].(string))
		r.LoadInstruction(&runtime.InstrAssign{
			Source: words[3].(variables.Symbol),
			Dest:   addr,
		})
		return addr, false
	}

	return nil, false
}
