package evaluator

import (
	"Pandora_Box/ast"
	"Pandora_Box/object"
)

// 以下字面量可以直接穷举, 所以直接创建对象
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
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
		return evalProgram(_node)
		// return evalStatements(_node.Statements)

	case *ast.ExpressionStatement:
		return Eval(_node.Expression)

	// 表达式
	case *ast.PrefixExpression:
		right := Eval(_node.Right)
		return evalPrefixExpression(_node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(_node.Left)
		right := Eval(_node.Right)
		return evalInfixExpression(_node.Operator, left, right)

	// 块
	case *ast.BlockStatement:
		return evalBlockStatement(_node)

	case *ast.IfExpression:
		return evalIfExpression(_node)

	case *ast.ReturnStatement:
		val := Eval(_node.ReturnValue)
		return &object.ReturnValue{Value: val}

	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)
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
		return NULL
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
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(op, left, right)
	case op == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
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
		return NULL
	}

}

func evalMinusPrefixOpExpression(right object.Object) object.Object {
	// 检查负号后面的对象类型是否为整型对象
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{
		Value: -value,
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
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

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

		if retVal, ok := result.(*object.ReturnValue); ok {
			return retVal.Value
		}

	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}

	}
	return result
}
