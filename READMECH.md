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
> 所以, **这是一款使用go开发的解释器**, 由于不确定未来会开发到何种程度, 我称之为潘多拉魔盒, 以表示希冀并警醒自己...



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


