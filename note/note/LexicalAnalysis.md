# 词法分析

![img.png](../img/img.png)

词法分析器:

1. Lexer结构体


```go
package lexer

type Lexer struct {
	input        string // 输入的代码
	position     int  // 所输入的字符串中的当前位置(指向当前字符串)
	readPosition int  // 所输入的字符串中的当前读取位置(指向当前字符之后的前一个字符)
	ch           byte // 当前正在查看的位置
	// only support for the ascii char
}
```

2. 词法单元

```go
package token

type TokenType string

type Token struct {
Type    TokenType
Literal string
}

```


3. REPL

REPL指Read-Eval-Print Loop(读取-求值-打印循环), 比较通俗的说法是控制台或交互模式.



