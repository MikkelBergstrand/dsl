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
			tokens.NTExpr,
			tokens.NTTerm,
			tokens.NTFactor,
		},
		StartSymbol: tokens.NTGoal,
	}

	/*
		grammar = tokens.Grammar{
			Terminals: []tokens.ItemType{
				tokens.TItemParOpen,
				tokens.TItemParClosed,
			},
			NonTerminals: []tokens.ItemType{
				tokens.NTTGoal,
				tokens.NTTList,
				tokens.NTTPair,
			},
			StartSymbol: tokens.NTTGoal,
		}*/
	cfg := parser.CreateCFG()
	fmt.Println(cfg)
	//cfg = parser.EliminateLeftRecursion(cfg, &grammar)

	first := parser.First(cfg, grammar)

	fmt.Println("FIRST")
	fmt.Println(first)

	follow := parser.Follow(cfg, grammar, first)

	fmt.Println("FOLLOW")
	fmt.Println(follow)

	action, _goto := parser.CreateLRTable(grammar, cfg, first)
	//fmt.Println(closure.String(cfg))

	/*ll_table := parser.MakeLLTable(grammar, cfg, first, follow)
	//(ll_parser := parser.NewParser(grammar, ll_table)

	*/
	words := make(chan tokens.Lexeme)
	go func() {
		words <- tokens.Lexeme{ItemType: tokens.ItemNumber, Value: "3"}
		words <- tokens.Lexeme{ItemType: tokens.ItemSemicolon, Value: ";"}
		words <- tokens.Lexeme{ItemType: tokens.ItemEOF, Value: ""}
	}()

	parser.LRParser(action, _goto, words, cfg, grammar)
}
