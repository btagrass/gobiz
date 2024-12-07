package utl

import (
	"github.com/expr-lang/expr"
)

func Eval(expression string, args ...map[string]any) (any, error) {
	program, err := expr.Compile(expression)
	if err != nil {
		return nil, err
	}
	var env any
	if len(args) > 0 {
		env = args[0]
	}
	return expr.Run(program, env)
}
