package parser

import (
	"Pandora_Box/ast"
	"Pandora_Box/lexer"
	"Pandora_Box/token"
	"fmt"
	"strconv"
	"testing"
)

// 优先级
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

// 词法单元到优先级的映射
var precedences = map[token.TokenType]int{
	// 等于和不等于
	token.EQ:  EQUALS,
	token.NEQ: EQUALS,

	// 小于和大于
	token.LT: LESSGREATER,
	token.GT: LESSGREATER,

	// 加法和减法
	token.PLUS:  SUM,
	token.MINUS: SUM,

	// 乘法和除法的词法单元
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,

	token.LPAREN: CALL,

	token.LBRACKET: INDEX,
}

/*
Parser a
*/
type Parser struct {
	l         *lexer.Lexer // 词法分析器的指针
	curToken  token.Token  // 当前词法单元
	peekToken token.Token  // 当前词法单元的下一个词法单元
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New 创建一个新的语法分析器
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// make: 初始化数据结构
	// map[token.TokenType]prefixParseFn ==> 从token.TokenType到prefixParseFn的映射

	/* 将所有前缀表达式的解析函数注册到语法分析器中 便于构建AST节点 */
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// 解析标识符
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	// 解析整数序列
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	// 解析感叹号
	p.registerPrefix(token.EXCLAMATION, p.parsePrefixExpression)
	// 解析负号
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	// 解析TRUE
	p.registerPrefix(token.TRUE, p.parseBoolean)
	// 解析FALSE
	p.registerPrefix(token.FALSE, p.parseBoolean)
	// 解析LPAREN
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	// 解析IF
	p.registerPrefix(token.IF, p.parseIfExpression)
	// 解析FUNCTION
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	// 解析STRING
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	// 解析ArrayLiteral
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	/* 为中缀表达式注册一个中缀解析函数 */
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	// 解析 +
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	// 解析 -
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	// 解析 /
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	// 解析 *
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	// 解析
	// p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	// 解析 ==
	p.registerInfix(token.EQ, p.parseInfixExpression)
	// 解析 !=
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	// 解析 <
	p.registerInfix(token.LT, p.parseInfixExpression)
	// 解析 >
	p.registerInfix(token.GT, p.parseInfixExpression)
	// 解析
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	//
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// 读取两个词法单元, 以设置curToken和peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken 同时将curToken和peekToken指针偏移至下一个词法单元
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram **
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // 创建一个指向ast.Program的指针  (program为AST的根节点)
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()

	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type] // 选取前缀表达式解析函数

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix() // 执行解析函数

	// 普拉特语法分析器核心
	/*
		普拉特解析算法是一种自顶向下的解析算法
	*/
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() { // 查看当前是否遍历到
		// 取出中缀表达式的解析函数
		infix := p.infixParseFns[p.peekToken.Type]

		// 无中缀表达式, 则直接返回
		if infix == nil {
			return leftExp
		}

		// 解析下一个token
		p.nextToken()
		// 使用中缀表达式进行解析
		leftExp = infix(leftExp)
		// 表达式: ((1+3) > 2;)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

/*
	具体的parse对象
*/

// 检查是否满足 let identifier = ... 这种格式
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}

	// 检查标识符
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	//
	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// 检查赋值号
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: 检查赋值号和分号之间的东西
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// 检查最后是否有分号
	//for !p.curTokenIs(token.SEMICOLON) {
	//	p.nextToken()
	//}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()

	// TODO: 检查Expression
	stmt.ReturnValue = p.parseExpression(LOWEST)

	// 检查是否有分号
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer
	stmt := &ast.ExpressionStatement{
		Token: p.curToken,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token: p.curToken,
	}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// checkParserErrors 利用Parser结构体中的errors数组来确定是否发生了语法解析错误
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

/*
	前缀解析函数
	中缀解析函数
*/

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(expression ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	// 设定优先级
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// 查询peekToken的优先级
func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence // 找到则返回
	}

	return LOWEST // 否则则返回最小优先级
}

// 查询curToken的优先级
func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.curToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()                  //
	p.nextToken()                                    // 词法单元前移
	expression.Right = p.parseExpression(precedence) // 填充AST的右节点

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseIfExpression 检测if和if-else语句
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.curToken,
	}

	// 使用expectPeek检测是否满足if的语法规则

	// 检测 (
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// 解析if括号中的条件表达式
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	// 检测 )
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// 检测 {
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	// 检测 else
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.curToken,
	}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseFunctionLiteral 解析函数序列
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{
		Token: p.curToken,
	}

	// 检测 (
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// 检测 {
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// 函数体解析
	lit.Body = p.parseBlockStatement()

	return lit
}

// parseFunctionParameters 解析参数列表
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// 当参数列表为空时, 后移一个词法单元, 然后返回
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	// 不为空时, cur和peek同时向后移动一个单元的token
	p.nextToken()

	// 第一个参数
	ident := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// 合并到总的形参列表中
	identifiers = append(identifiers, ident)

	// 遍历形参列表, 此处使用的是peekTokenIs检查是否还有逗号, 从而判断是否还有没有append的参数
	// 不为空时, 解析参数列表, 以逗号为
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{
		Token: p.curToken,
	}

	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	var list []ast.Expression

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}
