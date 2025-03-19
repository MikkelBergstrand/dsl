package parser

import (
	"dsl/instructions"
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
func DoActions(rule_id int, words []any, storage *storage.Storage, runtime *runtime.Runtime) any {
	fmt.Println(rule_id, words)
	switch rule_id {
	case 3:
		new_addr := storage.NewIntLiteral(0)
		runtime.LoadInstruction(&instructions.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: instructions.ADD,
		})
		return new_addr
	case 4:
		new_addr := storage.NewIntLiteral(0)
		runtime.LoadInstruction(&instructions.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: instructions.SUB,
		})
		return new_addr
	case 5:
		return words[0]
	case 6:
		new_addr := storage.NewIntLiteral(0)
		runtime.LoadInstruction(&instructions.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: instructions.MULT,
		})
		return new_addr
	case 7:
		new_addr := storage.NewIntLiteral(0)
		runtime.LoadInstruction(&instructions.InstrArithmetic{
			A:        words[0].(variables.Symbol),
			B:        words[2].(variables.Symbol),
			Result:   new_addr,
			Operator: instructions.DIV,
		})
		return new_addr
	case 8:
		return words[0]
	case 9:
		return words[0]
	case 10:
		return storage.NewIntLiteral(intval(words[0].(string)))
	case 11:
		return storage.GetVarAddr(words[0].(string))
	case 12: // New integer, eg. int a = 3
		addr := storage.NewIntVariable(words[1].(string))
		runtime.LoadInstruction(&instructions.InstrAssign{
			Source: words[3].(variables.Symbol),
			Dest:   addr,
		})
		return addr
	case 13: // Reassignment of integer, e.g. a = 3
		addr := storage.GetVarAddr(words[0].(string))
		runtime.LoadInstruction(&instructions.InstrAssign{
			Source: words[2].(variables.Symbol),
			Dest:   addr,
		})
		return addr
	case 16: // Declare scope
		storage.NewScope()
	case 17: // End scope
		storage.DestroyScope()
	case 20:
		return words
	case 21:
		return words[0]
	case 22:
		return words[0]
	}

	return nil
}
