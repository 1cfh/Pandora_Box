package evaluator

import (
	"Pandora_Box/ast"
	"Pandora_Box/object"
	"fmt"
)

// 以下字面量可以直接穷举, 所以直接创建对象
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Env) object.Object {
	// 根据AST上的节点对应的类型来确定对应的解析函数

	switch _node := node.(type) {
	// 表达式
	// 整型字面值序列
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: _node.Value,
		}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(_node.Value)

	// 语句
	case *ast.Program:
		return evalProgram(_node, env)
		// return evalStatements(_node.Statements)

	case *ast.ExpressionStatement:
		return Eval(_node.Expression, env)

	// 表达式
	case *ast.PrefixExpression:
		right := Eval(_node.Right, env)
		return evalPrefixExpression(_node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(_node.Left, env)
		right := Eval(_node.Right, env)
		return evalInfixExpression(_node.Operator, left, right)

	// 块
	case *ast.BlockStatement:
		return evalBlockStatement(_node, env)

	case *ast.IfExpression:
		return evalIfExpression(_node, env)

	case *ast.ReturnStatement:
		val := Eval(_node.ReturnValue, env)
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(_node.Value, env)
		if isError(val) {
			return val
		}
		// let的声明语句会产生环境的变化
		env.Set(_node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(_node, env)

	case *ast.FunctionLiteral:
		params := _node.Parameters
		body := _node.Body
		return &object.Function{
			Parameters: params,
			Env:        env,
			Body:       body,
		}

	case *ast.CallExpression:
		// 相当于获取函数指针
		function := Eval(_node.Function, env)
		if isError(function) {
			return function
		}
		// 解析参数列表
		args := evalExpressions(_node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		//
		return evalFunction(function, args)

	}
	return nil
}

func evalStatements(stmts []ast.Statement, env *object.Env) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)
		if retVal, ok := result.(*object.ReturnValue); ok {
			return retVal.Value
		}
	}
	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalExclamationOpExpression(right)
	case "-":
		return evalMinusPrefixOpExpression(right)
	default:
		return newError("unknown operator: %s%s", op, right.Type())
	}
}

func evalExclamationOpExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default: // 对于0,  !0 = TRUE
		return FALSE
	}
}

func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	// 根据中缀表达式的运算符进行switch
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(op, left, right)
	case op == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}

}

func evalIntegerInfixExpression(op string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{
			Value: leftVal + rightVal,
		}
	case "-":
		return &object.Integer{
			Value: leftVal - rightVal,
		}
	case "*":
		return &object.Integer{
			Value: leftVal * rightVal,
		}
	case "/":
		return &object.Integer{
			Value: leftVal / rightVal,
		}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}

}

func evalMinusPrefixOpExpression(right object.Object) object.Object {
	// 检查负号后面的对象类型是否为整型对象
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{
		Value: -value,
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Env) object.Object {
	condition := Eval(ie.Condition, env)

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalProgram(program *ast.Program, env *object.Env) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		// 检测result是否为object.Error, 如果是直接返回而不继续执行
		if retVal, ok := result.(*object.Error); ok {
			return retVal
		} else if retVal, ok := result.(*object.ReturnValue); ok {
			return retVal.Value
		}
	}

	// 返回空
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Env) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}

		}
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(obj object.Object) bool {
	_, ok := obj.(*object.Error)
	if ok {
		return true // 是Error Object
	} else {
		return false // 不是Error Object
	}
}

// evalIdentifier 在交互模式中, 对于标识符直接打印
func evalIdentifier(node *ast.Identifier, env *object.Env) object.Object {
	val, ok := env.Get(node.Value)
	if !ok { // 没找到标识符
		return newError("identifier not found: " + node.Value)
	}
	return val
}

// evalExpressions
func evalExpressions(exps []ast.Expression, env *object.Env) []object.Object {
	var result []object.Object

	// 遍历执行每一条expressions => 参数列表是从左到右进行执行的
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) { // 遇到错误直接返回
			return []object.Object{
				evaluated,
			}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}
	// 获得函数内部的一个新环境, 避免污染外部环境
	extendedEnv := extendFunctionEnv(function, args)
	// 执行函数体
	evaluated := Eval(function.Body, extendedEnv)
	// 如果是返回值类型, 剥出其中的Value字段
	return unwrapRetVal(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Env {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		// 以key-value形式设置值
		env.Set(param.Value, args[paramIdx])
	}

	// 设置env
	return env
}

func unwrapRetVal(obj object.Object) object.Object {
	if retVal, ok := obj.(*object.ReturnValue); ok {
		return retVal.Value
	}
	return obj
}
