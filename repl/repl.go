package repl

import (
	"bufio"
	"io"

	"github.com/jeremi-traverse/monkey/lexer"
	"github.com/jeremi-traverse/monkey/parser"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {

		// fmt.Fprintf(out, PROMT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		lexer := lexer.New(line)
		/*
			i := 0
				for tok := lexer.NextToken(); i < 6; tok = lexer.NextToken() {
					fmt.Fprintf(out, "%+v\n", tok)
					i++
				}
		*/
		p := parser.New(lexer)
		p.ParseProgram()

	}
}
