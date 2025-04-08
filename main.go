package main

import (
	"dsl/functions"
	"dsl/parser"
	"dsl/runtime"
	"dsl/scanner"
	"dsl/storage"
	"dsl/tokens"
	"dsl/variables"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	file_contents, err := os.ReadFile("testfiles/01.txt")
	if err != nil {
		log.Fatal(err)
		return
	}

	start := time.Now()
	_, scanner_stream := scanner.Lex("test_lexer", string(file_contents))

	_exit := false
	word_stream := make([]tokens.Token, 0)
	for !_exit {
		c := <-scanner_stream
		switch c.Category {
		case tokens.ItemEOF:
			word_stream = append(word_stream, c)
			_exit = true
		case tokens.ItemError:
			log.Fatal(c)
			return
		default:
			word_stream = append(word_stream, c)
		}
	}

	fmt.Println("Scanned in ", time.Since(start))

	grammar := tokens.Grammar{
		Terminals: []tokens.ItemType{
			tokens.ItemNumber,
			tokens.ItemText,
			tokens.ItemOpPlus,
			tokens.ItemOpMinus,
			tokens.ItemOpMult,
			tokens.ItemOpDiv,
			tokens.ItemIdentifier,
			tokens.ItemParOpen,
			tokens.ItemParClosed,
			tokens.ItemSemicolon,
			tokens.ItemKeyInt,
			tokens.ItemKeyBool,
			tokens.ItemEquals,
			tokens.ItemScopeOpen,
			tokens.ItemScopeClose,
			tokens.ItemComma,
			tokens.ItemBoolEqual,
			tokens.ItemBoolNot,
			tokens.ItemBoolNotEqual,
			tokens.ItemBoolLess,
			tokens.ItemBoolLessOrEqual,
			tokens.ItemBoolGreater,
			tokens.ItemBoolGreaterOrEqual,
			tokens.ItemBoolAnd,
			tokens.ItemBoolOr,
			tokens.ItemTrue,
			tokens.ItemFalse,
			tokens.ItemFunction,
		},
		NonTerminals: []tokens.ItemType{
			tokens.NTGoal,
			tokens.NTStatement,
			tokens.NTStatementList,
			tokens.NTExpr,
			tokens.NTTerm,
			tokens.NTFactor,
			tokens.NTScopeBegin,
			tokens.NTScopeClose,
			tokens.NTFunction,
			tokens.NTArgList,
			tokens.NTArgument,
			tokens.NTNExpr,
			tokens.NTAndTerm,
			tokens.NTNotTerm,
			tokens.NTRelExpr,
			tokens.NTRels,
			tokens.NTArgumentDeclaration,
			tokens.NTArgumentDeclarationList,
			tokens.NTVarType,
			tokens.NTFunctionClose,
			tokens.NTFunctionDefinition,
			tokens.NTFunctionBody,
		},
		StartSymbol: tokens.NTGoal,
	}

	cfg := parser.CreateCFG()

	start = time.Now()
	parser := parser.CreateLRParser(grammar, cfg, parser.First(cfg, grammar))
	fmt.Println("Created parse tables in ", time.Since(start))

	words := make(chan tokens.Token)
	go func() {
		for i := range word_stream {
			words <- word_stream[i]
		}
	}()

	fmt.Println(word_stream)

	storage := storage.NewStorage()
	runtime := runtime.New()

	generateGlobalFunctions(&runtime, &storage)

	start = time.Now()
	err = parser.Parse(words, cfg, grammar, &storage, &runtime)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Parsed in ", time.Since(start))

	start = time.Now()
	runtime.Run()
	fmt.Println("Program finished in", time.Since(start))
}

func generateGlobalFunctions(rt *runtime.Runtime, storage *storage.Storage) {
	def := functions.FunctionDefinition{
		ArgumentList: []functions.Argument{
			{Type: variables.INT, Identifier: "i"},
		},
		ReturnType: variables.NONE,
	}

	storage.NewFunctionScope(def)
	storage.NewFunction("echo", def)
	storage.NewLabel("echo", rt.NextInstruction())

	storage.LoadInstruction(&runtime.InstructionEcho{
		A: storage.GetVarAddr("i"),
	})
	storage.LoadInstruction(&runtime.InstrExitFunction{})
	storage.DestroyFunctionScope(rt)
}
