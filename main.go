package main

import (
	"dsl/instructionset"
	"dsl/parser"
	"dsl/scanner"
	"dsl/storage"
	"dsl/tokens"
	"fmt"
	"log"
	"os"
)

func main() {
	file_contents, err := os.ReadFile("testfiles/01.txt")
	if err != nil {
		log.Fatal(err)
		return
	}

	_, scanner_stream := scanner.Lex("test_lexer", string(file_contents))

	_exit := false
	word_stream := make([]tokens.Lexeme, 0)
	for !_exit {
		c := <-scanner_stream
		switch c.ItemType {
		case tokens.ItemEOF:
			word_stream = append(word_stream, c)
			_exit = true
			break
		case tokens.ItemError:
			log.Fatal(c)
			return
		default:
			word_stream = append(word_stream, c)
			fmt.Println(c)
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
		},
		NonTerminals: []tokens.ItemType{
			tokens.NTGoal,
			tokens.NTStatement,
			tokens.NTExpr,
			tokens.NTTerm,
			tokens.NTFactor,
		},
		StartSymbol: tokens.NTGoal,
	}

	cfg := parser.CreateCFG()
	fmt.Println(cfg)
	//cfg = parser.EliminateLeftRecursion(cfg, &grammar)

	first := parser.First(cfg, grammar)
	//follow := parser.Follow(cfg, grammar, first)

	action, _goto := parser.CreateLRTable(grammar, cfg, first)

	words := make(chan tokens.Lexeme)
	go func() {
		for i := range word_stream {
			words <- word_stream[i]
		}
	}()

	emitter := make(chan instructionset.Instruction)
	storage := storage.NewStorage()

	go func() {
		for {
			emitted := <-emitter
			fmt.Println("Emitted: ", emitted)
			emitted.Execute(&storage)
		}
	}()

	err = parser.LRParser(action, _goto, words, cfg, grammar, emitter, &storage)
	if err != nil {
		log.Fatal(err)
	}

}
