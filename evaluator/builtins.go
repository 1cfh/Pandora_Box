package evaluator

import "Pandora_Box/object"

// 内建函数的映射表
var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			// 检查len的参数长度  只允许接收一个参数
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			// 检查完参数后, 直接获取第一个值作为参数
			switch arg := args[0].(type) {
			case *object.String: // 字符串对象类型则直接返整型对象
				return &object.Integer{
					Value: int64(len(arg.Value)),
				}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}

		},
	},
}
