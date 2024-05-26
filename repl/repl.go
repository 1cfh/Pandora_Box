package repl

import (
	"Pandora_Box/lexer"
	"Pandora_Box/token"
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

		fmt.Println(line)

		l := lexer.New(line)

		// 遍历所有的token
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}

	}

}
