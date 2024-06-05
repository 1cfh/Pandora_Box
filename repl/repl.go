package repl

import (
	"Pandora_Box/evaluator"
	"Pandora_Box/lexer"
	"Pandora_Box/parser"
	"bufio"
	"fmt"
	"io"
)

// PROMPT prefix in each line
const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT) // PROMPT写入到标准输出流
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// 测试输入
		// fmt.Println(line)

		// 构建AST
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
		}

		evaluated := evaluator.Eval(program)

		if evaluated != nil {
			// before eval ast
			//io.WriteString(out, program.String())
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

		// 遍历所有的token
		//for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		//	fmt.Fprintf(out, "%+v\n", tok)
		//}man

	}

}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! Parser Errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
