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
	IDENT = "IDENT"
	INT   = "INT"

	// Operator identifiers
	ASSIGN = "="
	PLUS   = "+"
	SUBS   = "-"
	DIV    = "/"
	MULT   = "*"

	// Delimitiers
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Reserved words
	FUNCTION = "FUNCTION"
	LET      = "LET"
	RET      = "RETURN"
	CON      = "CONST"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"

	// Conditionals
	EQ  = "=="
	NOT = "!"
	NEQ = "!="
	LT  = "<"
	GT  = ">"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RET,
	"const":  CON,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
