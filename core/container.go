package core

type Container interface {
	AddCode(Code)
}

type NormalContainer struct {
	code []Code
}

func (c *NormalContainer) AddCode(code Code) {
	c.code = append(c.code, code)
}

func (c *NormalContainer) Execute(context *Context) ExecutionState {
	for _, code := range c.code {
		if state := code.Execute(context); state != NormalState {
			return state
		}
	}
	return NormalState
}
