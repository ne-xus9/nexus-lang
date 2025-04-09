package repl

import (
	"bufio"
	"fmt"
	"io"
	"nexus/lexer"
	"nexus/token"
)

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		s := scanner.Scan()
		if !s {
			return
		}
		line := scanner.Text()
		lex := lexer.New(line)
		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Printf("%v\n", tok)
		}
	}
}
