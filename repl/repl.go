package repl

import (
	"bufio"
	"fmt"
	"io"
	"nexus/evaluator"
	"nexus/lexer"
	"nexus/parser"
)

const PROMPT = ">>> "

func printParseErrors(w io.Writer, err []error) {
	for _, msg := range err {
		io.WriteString(w, "\t"+msg.Error()+"\n")
	}
}

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
		par := parser.New(lex)
		prog := par.ParseProgram()

		if err := par.Errors(); len(err) > 0 {
			printParseErrors(out, err)
			continue
		}

		ev := evaluator.Eval(prog)

		if ev != nil {
			io.WriteString(out, ev.Inspect())
			io.WriteString(out, "\n")
		}

	}
}
