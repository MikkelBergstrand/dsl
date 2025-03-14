package parser

import (
	"dsl/instructionset"
	"dsl/storage"
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
func DoActions(rule_id int, words []any, storage *storage.Storage, emitter chan<- instructionset.Instruction) any {
	switch rule_id {
	case 2:
		new_addr := storage.NewInt()
		emitter <- &instructionset.InstrArithmetic{
			A:        words[0].(int),
			B:        words[2].(int),
			Result:   new_addr,
			Operator: instructionset.ADD,
		}
		return new_addr
	case 3:
		return words[0].(int)
	case 4:
		return words[0].(int)
	case 5:
		new_addr := storage.NewInt()
		emitter <- &instructionset.InstrArithmetic{
			A:        words[0].(int),
			B:        words[2].(int),
			Result:   new_addr,
			Operator: instructionset.MULT,
		}
		return new_addr
	case 7:
		return words[0].(int)
	case 9:
		return storage.NewIntLiteral(intval(words[0].(string)))
	case 10:
		return storage.GetIntVarAddr(words[0].(string))
	case 11:
		addr := storage.NewIntVariable(words[1].(string))
		emitter <- &instructionset.InstrAssign{
			Source: words[3].(int),
			Dest:   addr,
		}
		return addr
	}

	//fmt.Println("Unhandled rule: ", rule_id)
	return nil
}
