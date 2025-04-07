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

	// Delimitiers
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	FUNCTION  = "FUNCTION" // fn
	LET       = "LET"
	RET       = "RETURN"
	CON       = "CONST"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RET,
	"const":  CON,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
