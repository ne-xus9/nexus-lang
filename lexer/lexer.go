package lexer

type Lexer struct {
	input   string
	pos     uint
	readPos uint
	ch      byte
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) readChar() {
	if l.readPos >= uint(len(l.input)) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}
