package core

import (
	"github.com/karlseguin/gerb/r"
	"reflect"
	"strings"
)

var Builtins = make(map[string]interface{})

func RegisterBuiltin(name string, f interface{}) {
	Builtins[strings.ToLower(name)] = reflect.ValueOf(f)
}

func init() {
	RegisterBuiltin("len", LenBuiltin)
}

func LenBuiltin(value interface{}) int {
	n, _ := r.ToLength(value)
	return n
}
