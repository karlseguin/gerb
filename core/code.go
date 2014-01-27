package core

type Executable interface {
	Execute(context *Context) ExecutionState
}

type Code interface {
	Executable
}
