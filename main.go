package main

import (
	"dsl/parser"
	"dsl/scanner"
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
	for !_exit {
		c := <-ch
		switch c.ItemType {
		case tokens.ItemEOF:
			_exit = true
			break
		case tokens.ItemError:
			log.Fatal(c)
			return
		default:
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
			tokens.ItemParClosed,
			tokens.ItemParOpen,
		},
		NonTerminals: []tokens.ItemType{
			tokens.NT_GOAL,
			tokens.NTExpr,
			tokens.NTTerm,
			tokens.NTFactor,
		},
	}
	cfg := parser.CreateCFG()
	cfg = parser.EliminateLeftRecursion(cfg, &grammar)

	first := parser.First(cfg, grammar)
	follow := parser.Follow(cfg, grammar, first)

	fmt.Println("FOLLOW")
	fmt.Println(follow)

}
