package junklang

type TokenType int

const (
	// Junk-specific tokens
	JUNK TokenType = iota
	BUNK
	SKUNK
	DUNK
	KLUNK
	SPUNK
	MUNK

	// Andre programmeringstokens
	IDENT
	NUM
	STR
	OPERATOR
	EOF
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
}
