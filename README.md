# Pandora_Box

```text
██████╗  █████╗ ███╗   ██╗██████╗  ██████╗ ██████╗  █████╗         ██████╗  ██████╗ ██╗  ██╗
██╔══██╗██╔══██╗████╗  ██║██╔══██╗██╔═══██╗██╔══██╗██╔══██╗        ██╔══██╗██╔═══██╗╚██╗██╔╝
██████╔╝███████║██╔██╗ ██║██║  ██║██║   ██║██████╔╝███████║        ██████╔╝██║   ██║ ╚███╔╝ 
██╔═══╝ ██╔══██║██║╚██╗██║██║  ██║██║   ██║██╔══██╗██╔══██║        ██╔══██╗██║   ██║ ██╔██╗ 
██║     ██║  ██║██║ ╚████║██████╔╝╚██████╔╝██║  ██║██║  ██║███████╗██████╔╝╚██████╔╝██╔╝ ██╗
╚═╝     ╚═╝  ╚═╝╚═╝  ╚═══╝╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═════╝  ╚═════╝ ╚═╝  ╚═╝
```

- why?

> Hi, there! I'm learning about program analysis recently. 
> 
> However, I knew so little about the Compiler. Until I discovered the interpreter development book called the Monkey Book in the industry, I decided to start learning relevant knowledge through a small project and expand interesting functions.

> Therefore, **this is an interpreter developed by Golang**. Since I am not sure how far it will be developed in the future, I call it Pandora's Box to express my hope and alert myself...




## DevNotion


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

