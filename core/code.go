package core

type Code interface {
	Execute(context *Context) ExecutionState
}
