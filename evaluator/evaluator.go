package evaluator

import (
	"nexus/ast"
	"nexus/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefix(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfix(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.IfExpression:
		return evalIf(node)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func nativeBooleanObject(i bool) *object.Boolean {
	if i {
		return TRUE
	}
	return FALSE
}

func evalPrefix(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalNot(right)
	case "-":
		return evalMinus(right)
	default:
		return NULL
	}
}

func evalNot(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinus(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return NULL
	}
	val := right.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func evalInfix(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntInfix(op, left, right)
	case op == "==":
		return nativeBooleanObject(left == right)
	case op == "!=":
		return nativeBooleanObject(left != right)
	default:
		return NULL
	}
}

func evalIntInfix(op string, left, right object.Object) object.Object {
	lval := left.(*object.Integer).Value
	rval := right.(*object.Integer).Value
	result := &object.Integer{}

	switch op {
	case "+":
		result.Value = lval + rval
	case "-":
		result.Value = lval - rval
	case "*":
		result.Value = lval * rval
	case "/":
		result.Value = lval / rval
	case "<":
		return nativeBooleanObject(lval < rval)
	case ">":
		return nativeBooleanObject(lval > rval)
	case "==":
		return nativeBooleanObject(lval == rval)
	case "!=":
		return nativeBooleanObject(lval != rval)
	default:
		return NULL
	}
	return result
}

func evalIf(ie *ast.IfExpression) object.Object {
	cond := Eval(ie.Condition)
	if isTruthy(cond) {
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
