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

	_, ch := scanner.Lex("test_lexer", string(file_contents))

	_exit := false
	word_stream := make([]tokens.Lexeme, 0)
	for !_exit {
		c := <-ch
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
		words <- tokens.Lexeme{ItemType: tokens.ItemNumber, Value: "12"}
		words <- tokens.Lexeme{ItemType: tokens.ItemOpMult, Value: "*"}
		words <- tokens.Lexeme{ItemType: tokens.ItemNumber, Value: "19"}
		words <- tokens.Lexeme{ItemType: tokens.ItemOpPlus, Value: "+"}
		words <- tokens.Lexeme{ItemType: tokens.ItemNumber, Value: "40"}
		words <- tokens.Lexeme{ItemType: tokens.ItemOpMult, Value: "*"}
		words <- tokens.Lexeme{ItemType: tokens.ItemNumber, Value: "83"}
		words <- tokens.Lexeme{ItemType: tokens.ItemSemicolon, Value: ";"}
		words <- tokens.Lexeme{ItemType: tokens.ItemEOF, Value: ""}
	}()

	emitter := make(chan instructionset.Instruction)
	storage := storage.Storage{}

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
