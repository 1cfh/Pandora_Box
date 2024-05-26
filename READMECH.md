# Pandora_Box

```text
██████╗  █████╗ ███╗   ██╗██████╗  ██████╗ ██████╗  █████╗         ██████╗  ██████╗ ██╗  ██╗
██╔══██╗██╔══██╗████╗  ██║██╔══██╗██╔═══██╗██╔══██╗██╔══██╗        ██╔══██╗██╔═══██╗╚██╗██╔╝
██████╔╝███████║██╔██╗ ██║██║  ██║██║   ██║██████╔╝███████║        ██████╔╝██║   ██║ ╚███╔╝ 
██╔═══╝ ██╔══██║██║╚██╗██║██║  ██║██║   ██║██╔══██╗██╔══██║        ██╔══██╗██║   ██║ ██╔██╗ 
██║     ██║  ██║██║ ╚████║██████╔╝╚██████╔╝██║  ██║██║  ██║███████╗██████╔╝╚██████╔╝██╔╝ ██╗
╚═╝     ╚═╝  ╚═╝╚═╝  ╚═══╝╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═════╝  ╚═════╝ ╚═╝  ╚═╝
```

- 为什么会有这个代码仓?

> 苦于程序分析中对于编译原理的浅薄理解, 特别是在看到一本业内称为猴书的解释器开发书籍后, 我决定跟着写一个自己的解释器, 并在其基础上做一些有趣的拓展。
> 
> 所以, **这是一款使用Go开发的解释器**, 由于不确定未来会开发到何种程度, 我称之为潘多拉魔盒, 以表示希冀并警醒自己...



## 开发日志

### Token

1. Token struct

```go
package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
```

### Lexical Analysis

1. lexical analysis unit

```go
package lexer

type Lexer struct {
	input        string // input your code
	position     int  // current char index
	readPosition int  // read char index (readPosition = position + 1)
	ch           byte // ch := input[position] 
}
// only support ascii
```

2. REPL

主要就是进行词法分析后, 将结果打印到标准输出即可

```go
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

```

### Syntactic Analysis

主要涉及的数据结构有


拉普拉斯语法分析器