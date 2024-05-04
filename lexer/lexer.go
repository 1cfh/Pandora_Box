package lexer

import (
	"GolangBox_Monkey/token"
)

type Lexer struct {
	input        string
	position     int  // 所输入的字符串中的当前位置(指向当前字符串)
	readPosition int  // 所输入的字符串中的当前读取位置(指向当前字符之后的前一个字符)
	ch           byte // 当前正在查看的位置
	// only support for the ascii char
}

// create a new lexer section
func New(input string) *Lexer {
	l := &Lexer{input: input, readPosition: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // null
	} else {
		l.ch = l.input[l.readPosition] // current char
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	// 检查词法单元, 并给词法单元创建相应的token对象
	switch l.ch {
	case '=':
		// switch中的case查询不支持多字符
		var x = l.peekChar()
		// ==
		if x == '=' {
			ch := l.ch   //	存第一个字符
			l.readChar() // 移动到下一个字符
			// tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
			tok = newToken(token.EQ, string(ch)+string(l.ch))
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		var x = l.peekChar()

		// !=
		if x == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.NEQ, string(ch)+string(l.ch))
		} else {
			tok = newToken(token.EXCLAMATION, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)

	case '<':
		tok = newToken(token.LT, l.ch)

	case '>':
		tok = newToken(token.GT, l.ch)

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()          // 读取标识符
			tok.Type = token.LookupIdent(tok.Literal) // 根据关键字字典寻找对应的token类型
			return tok
		} else if isDigit(l.ch) { // 处理数字
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else { // 异常类型
			tok = newToken(token.ILLEGAL, l.ch)
		}

	}

	// 移动指针
	l.readChar()

	return tok
}

// 创建Token对象
func newToken[T byte | string](tokenType token.TokenType, ch T) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position // 记录起始位置
	for isLetter(l.ch) {   // 遍历字符串序列, 直至非letter字符
		l.readChar()
	}
	return l.input[position:l.position] // 从开始位置到非letter字符前一个字符即为 标识符
}

// 标识符只允许使用字母和_
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// 依靠指针移动, 处理特殊字符: whitespace, \t, \n, \r
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// 检查当前字节是否在数字范围内
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// 检查是否为INT, 然后截取获得token
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 返回当前查询的字符字节
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
