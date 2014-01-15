package core

type Value interface {
	Resolve(context *Context) interface{}
}

type IntValue struct {
	value int
}

func (v *IntValue) Resolve(context *Context) interface{} {
	return v.value
}

type FloatValue struct {
	value float64
}

func (v *FloatValue) Resolve(context *Context) interface{} {
	return v.value
}

type CharValue struct {
	value byte
}

func (v *CharValue) Resolve(context *Context) interface{} {
	return v.value
}
