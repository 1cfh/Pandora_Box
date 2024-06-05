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

#### 接口 和 数据结构

​	首先声明三种接口，这是实现更复杂的数据结构的基础

``` go
type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Expression interface{
    Node
    expressionNode()
}

```
​	然后是一些通用的数据结构
``` go
// 程序
type Program struct {
	Statements []Statement
}


// 标识符结构体
type Identifier struct{
    Token token.Token
    Value string
}


// 词法解析器
type Parser struct{
    l *lexer.Lexer
    
    errors []string
    curToken token.Token
    peekToken token.Token
    
    prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
    
}	// Parser.parseStatment方法中汇聚了解析let,return,expression等函数
```



#### 解析let语句

``` go
type LetStatement struct{
    Token token.Token
    Name *Identifier
    Value Expression
}
```





#### 解析return语句

``` go
type ReturnStatement struct{
    Token token.Token
    ReturnValue Expression
}
```





#### 解析expression

​	PrefixExpression

``` go
type PrefixExpression struct {
	Token    token.Token // 前缀词法单元
	Operator string      // 前缀表达式中的运算符
	Right    Expression  // 表达式
}
```

​	InfixExpression

```go
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}
```



#### Pratt Parsers

参考我的博客链接：https://ch31sbest.github.io/2024/05/27/Pratt-Parsers-WorkingMode/







### 求值

> 一些术语：

- JIT：Just in time
  - 一种将热点代码的字节码编译成本地机器码的编译技术
- 解释AST

  - 词法分析：

    - 将源代码分解成一系列Token，Token是最小的语法单元。

    - Token会组成一个Token序列

  - 语法分析

    - 根据语法规则将Token序列填充入语法树AST

    - 每个AST的节点为一个Token序列=>Expression

  - 语义分析

    - 检查抽象语法树是否符合语言的语义规则，例如类型检查、变量声明

    - 解释执行

    - 遍历抽象语法树，根据节点的类型执行相应的操作，解释器根据抽象语法树的结构逐步执行源代码对应的操作

- expression、statement和program这三者
  - expression（表达式）：
    - 代码片段，一般是常量、变量、操作符或者函数调用的组合，它可以被求值并返回一个值
    - 例如：`2+3`，`x*y`，`foo(5)`
  - statement（语句）：
    - 独立的语句，可以改变程序状态（如赋值语句）、控制程序流（如条件语句和循环语句）、或者进行输入或输出操作
    - 语句通常不返回值
    - 例如：`x=5`，`if(x>0){y=x;}`
  - Program（程序）：
    - 一个 `Program` 是由多个语句、表达式、函数定义、类定义等组成的完整代码块。
    - 程序是可以被解释器或编译器执行的完整单元，表示一个完整的任务或应用。
    - 例如：一个包含多行代码的代码块



#### 树遍历解释器

基本的思路就是遍历解析AST



#### 面向对象

简单的来说，就是将原生的数据类型都以对象的形式进行封装

- 对象类型接口

``` go
package object
type ObjectType string

type Object interface{
    Type() ObjectType
    Inspect() string
}
```



- 整数

``` go
type Integer struct{
    Value int64
}
```



- 布尔值

``` go
type Boolean struct{
    Value bool
}
```



- 空值

``` go
type Null struct{}
```



