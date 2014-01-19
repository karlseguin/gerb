package core

import (
	"fmt"
	"github.com/karlseguin/gerb/r"
	"reflect"
)

type Value interface {
	Resolve(context *Context) interface{}
}

type StaticValue struct {
	value interface{}
}

func (v *StaticValue) Resolve(context *Context) interface{} {
	return v.value
}

type DynamicFieldType int

const (
	FieldType DynamicFieldType = iota
	MethodType
	IndexedType
)

type DynamicValue struct {
	names []string
	types []DynamicFieldType
	args  [][]Value
}

func (v *DynamicValue) Resolve(context *Context) interface{} {
	var d interface{} = context.Data
	isRoot := true
	isAlias := false

	for i, l := 0, len(v.names); i < l; i++ {
		name := v.names[i]
		t := v.types[i]

		if t == FieldType {
			if d = r.ResolveField(d, name); d == nil {
				if isRoot {
					if alias, ok := FunctionAliases[name]; ok {
						d = alias
						isAlias = true
						isRoot = false
						continue
					}
				} else if pkg, ok := OtherAliases[v.names[i-1]]; ok {
					if alias, ok := pkg[name]; ok {
						return alias
					}
				}
				return v.loggedNil(i)
			}
		} else if t == IndexedType {
			if len(name) > 0 {
				if d = r.ResolveField(d, name); d == nil {
					return v.loggedNil(i)
				}
			}
			if d = unindex(d, v.args[i], context); d == nil {
				return v.loggedNil(i)
			}
		} else if t == MethodType {
			if d = run(d, name, v.args[i], isRoot, isAlias, context); d == nil {
				return v.loggedNilMethod(i)
			}
		}
		isAlias = false
		isRoot = false
	}
	return r.ResolveFinal(d)
}

func (v *DynamicValue) loggedNil(index int) interface{} {
	if index == 0 {
		Log.Error(fmt.Sprintf("%s is undefined", v.names[index]))
	} else {
		Log.Error(fmt.Sprintf("%s.%s is undefined", v.names[index-1], v.names[index]))
	}
	return nil
}

func (v *DynamicValue) loggedNilMethod(index int) interface{} {
	if index == 0 {
		Log.Error(fmt.Sprintf("%s is undefined", v.names[index]))
	} else {
		Log.Error(fmt.Sprintf("%s.%s is undefined or had undefined parameters", v.names[index-1], v.names[index]))
	}
	return nil
}

func unindex(container interface{}, params []Value, context *Context) interface{} {
	valueLength := len(params)
	if valueLength == 0 {
		return nil
	}

	value := reflect.ValueOf(container)
	kind := value.Kind()
	if kind == reflect.Array || kind == reflect.Slice || kind == reflect.String {
		length := value.Len()
		first, ok := r.ToInt(params[0].Resolve(context))
		if ok == false {
			return nil
		}
		if first < 0 {
			first = 0
		} else if first > length-1 {
			first = length
		}
		second := length
		if valueLength == 2 {
			second, ok = r.ToInt(params[1].Resolve(context))
			if ok == false {
				return nil
			}
			if second > length {
				second = length
			}
		}
		return value.Slice(first, second).Interface()
	} else if kind == reflect.Map {
		indexValue := reflect.ValueOf(params[0].Resolve(context))
		return value.MapIndex(indexValue).Interface()
	}
	return nil
}

func run(container interface{}, name string, params []Value, isRoot, isAlias bool, context *Context) interface{} {
	defer func() {
		if r := recover(); r != nil {
			Log.Error(r)
		}
	}()
	if isRoot {
		return runBuiltIn(name, params, context)
	}
	if isAlias {
		return runAlias(container.(map[string]interface{}), name, params, context)
	}

	c := reflect.ValueOf(container)
	m := r.Method(c, name)
	if m.IsValid() == false {
		return nil
	}
	v := make([]reflect.Value, len(params)+1)
	v[0] = c
	for index, param := range params {
		paramValue := reflect.ValueOf(param.Resolve(context))
		if paramValue.IsValid() == false {
			return nil
		}
		v[index+1] = paramValue
	}
	if returns := m.Call(v); len(returns) > 0 {
		return returns[0].Interface()
	}
	return nil
}

func runBuiltIn(name string, params []Value, context *Context) interface{} {
	return runFromLookup(Builtins, name, params, context)
}

func runAlias(pkg map[string]interface{}, name string, params []Value, context *Context) interface{} {
	return runFromLookup(pkg, name, params, context)
}

func runFromLookup(lookup map[string]interface{}, name string, params []Value, context *Context) interface{} {
	m, ok := lookup[name]
	if ok == false {
		return nil
	}

	switch typed := m.(type) {
	case reflect.Value:
		v := make([]reflect.Value, len(params))
		for index, param := range params {
			v[index] = reflect.ValueOf(param.Resolve(context))
		}
		if returns := typed.Call(v); len(returns) > 0 {
			return returns[0].Interface()
		}
	case reflect.Type:
		if len(params) != 1 {
			Log.Error(fmt.Sprintf("Conversion to %s should have 1 parameter", name))
			return nil
		}
		v := reflect.ValueOf(params[0].Resolve(context))
		if v.Type().ConvertibleTo(typed) == false {
			Log.Error(fmt.Sprintf("Cannot convert %s to %s", v.Type().Name(), typed.Name()))
			return nil
		}
		return v.Convert(typed).Interface()
	}
	return nil
}
