package core

type Container interface {
	AddExecutable(Executable)
}

type NormalContainer struct {
	executable []Executable
}

func (c *NormalContainer) AddExecutable(executable Executable) {
	c.executable = append(c.executable, executable)
}

func (c *NormalContainer) Execute(context *Context) ExecutionState {
	for _, executable := range c.executable {
		if state := executable.Execute(context); state != NormalState {
			return state
		}
	}
	return NormalState
}
