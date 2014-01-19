package core

import (
	"github.com/karlseguin/gerb/r"
	"time"
)

type OperationFactory func(a, b Value) Value

var Operations = map[byte]OperationFactory{
	'+': AddOperation,
	'-': SubOperation,
	'*': MultiplyOperation,
	'/': DivideOperation,
	'%': ModuloOperation,
}

func AddOperation(a, b Value) Value {
	return &AdditiveValue{a, b, false}
}

func SubOperation(a, b Value) Value {
	return &AdditiveValue{a, b, true}
}

func MultiplyOperation(a, b Value) Value {
	return &MultiplicativeValue{a, b, false}
}

func DivideOperation(a, b Value) Value {
	return &MultiplicativeValue{a, b, true}
}

func ModuloOperation(a, b Value) Value {
	return &ModulatedValue{a, b}
}

type AdditiveValue struct {
	a      Value
	b      Value
	negate bool
}

func (v *AdditiveValue) Resolve(context *Context) interface{} {
	a := v.a.Resolve(context)
	b := v.b.Resolve(context)
	if na, ok := r.ToInt(a); ok {
		if nb, ok := r.ToInt(b); ok {
			if v.negate {
				nb = -nb
			}
			return na + nb
		}
		return 0
	}
	if fa, ok := r.ToFloat(a); ok {
		if fb, ok := r.ToFloat(b); ok {
			if v.negate {
				fb = -fb
			}
			return fa + fb
		}
		return 0
	}
	if ta, ok := a.(time.Duration); ok {
		if tb, ok := b.(time.Duration); ok {
			if v.negate {
				return ta - tb
			}
			return ta + tb
		}
	}
	return 0
}

type MultiplicativeValue struct {
	a      Value
	b      Value
	divide bool
}

func (v *MultiplicativeValue) Resolve(context *Context) interface{} {
	a := v.a.Resolve(context)
	b := v.b.Resolve(context)
	if na, ok := r.ToInt(a); ok {
		if nb, ok := r.ToInt(b); ok {
			if v.divide {
				return na / nb
			}
			return na * nb
		}
		return 0
	}
	if fa, ok := r.ToFloat(a); ok {
		if fb, ok := r.ToFloat(b); ok {
			if v.divide {
				return fa / fb
			}
			return fa * fb
		}
		return 0
	}
	if ta, ok := a.(time.Duration); ok {
		if tb, ok := b.(time.Duration); ok {
			if v.divide {
				return ta / tb
			}
			return ta * tb
		}
	}
	return 0
}

type ModulatedValue struct {
	a Value
	b Value
}

func (v *ModulatedValue) Resolve(context *Context) interface{} {
	a := v.a.Resolve(context)
	b := v.b.Resolve(context)
	if na, ok := r.ToInt(a); ok {
		if nb, ok := r.ToInt(b); ok {
			return na % nb
		}
	}
	return 0
}
