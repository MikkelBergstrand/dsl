package scanner

import (
	"dsl/tokens"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

type lexer struct {
	name  string
	input string
	start int
	pos   int
	width int
	items chan tokens.Lexeme
}

const eof rune = '\x00' // necessary in 2025?

type stateFn func(*lexer) stateFn

func Lex(name string, input string) (*lexer, chan tokens.Lexeme) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan tokens.Lexeme),
	}
	go l.run()
	return l, l.items
}

func (l *lexer) run() {
	for state := lexInsideScope; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) emit(t tokens.ItemType) {
	l.items <- tokens.Lexeme{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	_rune, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += l.width
	return _rune
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	_rune := l.next()
	l.backup()
	return _rune
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) acceptRegex(valid string) {
	re := regexp.MustCompile(`^[` + valid + `]+`)

	ret := re.FindStringIndex(l.input[l.pos:])
	fmt.Println(ret)
	if ret != nil {
		l.pos += ret[1]
	}
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- tokens.Lexeme{
		tokens.ItemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func isAlphaNumeric(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9')
}

func isAlpha(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func lexInsideExpression(l *lexer) stateFn {
	for {
		r := l.next()

		if isSpace(r) {
			l.ignore()
		} else if r == '"' {
			return lexQuote
		} else if r == '+' {
			l.emit(tokens.ItemOpPlus)
			return lexInsideExpression
		} else if r == '-' {
			l.emit(tokens.ItemOpMinus)
			return lexInsideExpression
		} else if r == '*' {
			l.emit(tokens.ItemOpMult)
			return lexInsideExpression
		} else if r == '/' {
			l.emit(tokens.ItemOpDiv)
			return lexInsideExpression
		} else if '0' <= r && r <= '9' {
			l.backup()
			return lexNumber
		} else if isAlpha(r) {
			l.backup()
			return lexIdentifier
		} else if r == ';' {
			l.emit(tokens.ItemSemicolon)
			return lexInsideScope
		} else {
			return l.errorf("Non-terminated expression, missing a semicolon?")
		}
	}
}

func lexInsideScope(l *lexer) stateFn {
	for {
		r := l.next()
		if isSpace(r) {
			l.ignore()
		} else if r == eof {
			l.emit(tokens.ItemEOF)
			return nil
		} else {
			l.backup()
			return lexInsideExpression
		}
	}
}

func lexNumber(l *lexer) stateFn {
	//Optional leading sign
	l.accept("+-")
	l.acceptRun("0123456789")

	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(tokens.ItemNumber)
	return lexInsideExpression
}

func lexIdentifier(l *lexer) stateFn {
	//We know the first is alphanumeric
	fmt.Println("yo")
	l.acceptRegex("0-9a-zA-Z_")
	l.emit(tokens.ItemIdentifier)
	return lexInsideExpression
}

func lexQuote(l *lexer) stateFn {
	for {
		r := l.next()
		if r == eof || r == '\n' {
			return l.errorf("Unterminated string literal")
		} else if r == '"' {
			l.emit(tokens.ItemText)
			return lexInsideExpression
		}
	}
}
