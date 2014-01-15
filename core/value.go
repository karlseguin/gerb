package core

type Value interface {
	Resolve(context *Context) interface{}
}

type StaticValue struct {
	value interface{}
}

func (v *StaticValue) Resolve(context *Context) interface{} {
	return v.value
}
