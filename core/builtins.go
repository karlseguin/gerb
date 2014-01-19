package core

import (
	"reflect"
	"strings"
	"github.com/karlseguin/gerb/r"
)

var Builtins = make(map[string]reflect.Value)

func RegisterBuiltin(name string, f interface{}) {
	Builtins[strings.ToLower(name)] = reflect.ValueOf(f)
}


func init() {
	RegisterBuiltin("len", LenBuiltin)
}

func LenBuiltin(value interface{}) int {
	n, _  := r.ToLength(value)
	return n
}
