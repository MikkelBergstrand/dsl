package parser

import (
	"dsl/runtime"
	"dsl/storage"
	"dsl/variables"
	"fmt"
	"strconv"
)

// Convert string to integer. Should not fail!
func intval(s string) int {
	intval, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return intval
}
func DoActions(rule_id int, words []any, storage *storage.Storage, r *runtime.Runtime) any {
	fmt.Println(rule_id, words)
	switch rule_id {
	case 3:
		new_addr := storage.NewIntLiteral()
		r.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.ADD,
		})
		return new_addr
	case 4:
		new_addr := storage.NewIntLiteral()
		r.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.SUB,
		})
		return new_addr
	case 5:
		return words[0]
	case 6:
		new_addr := storage.NewIntLiteral()
		r.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.MULT,
		})
		return new_addr
	case 7:
		new_addr := storage.NewIntLiteral()
		r.LoadInstruction(&runtime.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: runtime.DIV,
		})
		return new_addr
	case 8:
		return words[0]
	case 9:
		return words[0]
	case 10:
		addr := storage.NewIntLiteral()
		r.LoadInstruction(&runtime.InstrLoadImmediate{
			Dest:  addr,
			Value: intval(words[0].(string)),
		})
		return addr
	case 11:
		return storage.GetVarAddr(words[0].(string))
	case 12: // New integer, eg. int a = 3
		addr := storage.NewIntVariable(words[1].(string))
		r.LoadInstruction(&runtime.InstrAssign{
			Source: words[3].(variables.Symbol),
			Dest:   addr,
		})
		return addr
	case 13: // Reassignment of integer, e.g. a = 3
		addr := storage.GetVarAddr(words[0].(string))
		r.LoadInstruction(&runtime.InstrAssign{
			Source: words[2].(variables.Symbol),
			Dest:   addr,
		})
		return addr
	case 16: // Declare scope
		storage.NewScope()
	case 17: // End scope
		storage.DestroyScope()
	case 19:
		arg_list := []variables.Symbol{words[2].(variables.Symbol)}

		passed_arg_list := make([]variables.Symbol, len(arg_list))
		for i := range passed_arg_list {
			passed_arg_list[i] = arg_list[i]
			passed_arg_list[i].Scope += 1
		}

		fn := storage.Functions[words[0].(string)]
		r.LoadInstruction(&runtime.InstrCallFunction{
			ArgumentList: arg_list,
			Func:         fn,
			AddressStart: storage.CurrentScope.Offset,
		})

		for i := range fn.ArgumentList {
			r.LoadInstruction(&runtime.InstrAssign{
				Source: passed_arg_list[i],
				Dest:   variables.Symbol{Scope: 0, Offset: i},
			})
		}

		r.LoadInstruction(&runtime.InstrJmp{
			NewPC: fn.InstructionPointer,
		})

	case 20:
		return words
	case 21:
		return words[0]
	case 22:
		return words[0]
	}

	return nil
}
