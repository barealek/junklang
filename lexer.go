package junklang

import (
	"strings"
)

type Lexer struct {
	input   string
	tokens  []Token
	current int
	line    int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:   input,
		tokens:  make([]Token, 0),
		current: 0,
		line:    1,
	}
}

func (l *Lexer) Tokenize() []Token {
	for l.current < len(l.input) {
		l.skipWhitespace()
		if l.current >= len(l.input) {
			break
		}

		c := l.input[l.current]
		switch {
		case isLetter(c):
			l.tokenizeWord()
		case isDigit(c):
			l.tokenizeNumber()
		case isOperator(c):
			l.tokenizeOperator()
		default:
			panic("Unexpected character: " + string(c))
		}
	}

	l.tokens = append(l.tokens, Token{Type: EOF})
	return l.tokens
}

func (l *Lexer) tokenizeWord() {
	start := l.current
	for l.current < len(l.input) && (isLetter(l.input[l.current]) || isDigit(l.input[l.current])) {
		l.current++
	}
	word := l.input[start:l.current]

	var tokenType TokenType
	switch word {
	case "junk":
		tokenType = JUNK
	case "bunk":
		tokenType = BUNK
	case "skunk":
		tokenType = SKUNK
	case "dunk":
		tokenType = DUNK
	case "klunk":
		tokenType = KLUNK
	case "spunk":
		tokenType = SPUNK
	case "munk":
		tokenType = MUNK
	default:
		tokenType = IDENT
	}

	l.tokens = append(l.tokens, Token{Type: tokenType, Value: word, Line: l.line})
}

func (l *Lexer) tokenizeNumber() {
	start := l.current
	for l.current < len(l.input) && (isDigit(l.input[l.current]) || l.input[l.current] == '.') {
		l.current++
	}
	l.tokens = append(l.tokens, Token{Type: NUM, Value: l.input[start:l.current], Line: l.line})
}

func (l *Lexer) tokenizeOperator() {
	operator := string(l.input[l.current])
	l.tokens = append(l.tokens, Token{Type: OPERATOR, Value: operator, Line: l.line})
	l.current++
}

func (l *Lexer) skipWhitespace() {
	for l.current < len(l.input) {
		c := l.input[l.current]
		if c == ' ' || c == '\t' || c == '\r' {
			l.current++
		} else if c == '\n' {
			l.line++
			l.current++
		} else {
			break
		}
	}
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isOperator(c byte) bool {
	return strings.ContainsRune("+-*/(){}=,", rune(c))
}
