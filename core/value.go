package core

import (
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
	args [][]Value
}

func (v *DynamicValue) Resolve(context *Context) interface{} {
	var d interface{} = context.Data
	for i, l := 0, len(v.names); i < l; i++ {

		name := v.names[i]
		t := v.types[i]
		if t == FieldType || t == IndexedType {
			if d = r.ResolveField(d, name); d == nil {
				return nil
			}
			if t == IndexedType {
				if d = unindex(d, v.args[i], context); d == nil {
					return nil
				}
			}
		}
	}
	return r.ResolveFinal(d)
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
