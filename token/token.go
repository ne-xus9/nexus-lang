package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Literal identifiers
	INDENT = "INDENT"
	INT    = "INT"

	// Operator identifiers
	ASSIGN = "="
	PLUS   = "+"

	// Delimitiers
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	FUNCTION  = "FUNCTION" // fn
	LET       = "LET"
)
