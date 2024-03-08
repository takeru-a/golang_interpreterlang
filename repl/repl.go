package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/takeru-a/golang_interpreterlang/evaluator"
	"github.com/takeru-a/golang_interpreterlang/lexer"
	"github.com/takeru-a/golang_interpreterlang/object"
	"github.com/takeru-a/golang_interpreterlang/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Println(AQUAMARINE)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan() 
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

const AQUAMARINE = `
  ___                                                _
 / _ \                                              (_)
/ /_\ \  __ _  _   _   __ _  _ __ ___    __ _  _ __  _  _ __    ___
|  _  | / _' || | | | / _' || '_ \' _ \  / _' || '__|| || '_ \  / _ \
| | | || (_| || |_| || (_| || | | | | || (_| || |   | || | | ||  __/
\_| |_/ \__, | \__,_| \__,_||_| |_| |_| \__,_||_|   |_||_| |_|\___|
		  | |
		  |_|

`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Error in Aquamarine script")
	io.WriteString(out, " syntax errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
