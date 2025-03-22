package main

import (
	"dsl/parser"
	"dsl/runtime"
	"dsl/scanner"
	"dsl/storage"
	"dsl/tokens"
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

	grammar := tokens.Grammar{
		Terminals: []tokens.ItemType{
			tokens.ItemNumber,
			tokens.ItemOpPlus,
			tokens.ItemOpMinus,
			tokens.ItemOpMult,
			tokens.ItemOpDiv,
			tokens.ItemIdentifier,
			tokens.ItemParOpen,
			tokens.ItemParClosed,
			tokens.ItemSemicolon,
			tokens.ItemKeyInt,
			tokens.ItemEquals,
			tokens.ItemScopeOpen,
			tokens.ItemScopeClose,
			tokens.ItemComma,
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
		},
		StartSymbol: tokens.NTGoal,
	}

	cfg := parser.CreateCFG()

	start := time.Now()
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
	runtime := runtime.New(&storage)

	start = time.Now()
	err = parser.Parse(words, cfg, grammar, &storage, &runtime)
	fmt.Println("Parsed in ", time.Since(start))

	runtime.Run(&storage)

	if err != nil {
		log.Fatal(err)
	}

	select {}
}
