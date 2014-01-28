package gerb

import (
	"github.com/karlseguin/gerb/core"
	"strings"
)

var CodeFactories = map[string]CodeFactory{
	"if": IfFactory,
}

var endScope = new(EndScope)

type CodeFactory func(*core.Parser) (core.Code, error)

func createCodeTag(p *core.Parser) (code core.Code, err error) {
	token, err := p.ReadToken()
	if err != nil {
		return nil, err
	}
	length := len(token)
	if length == 0 {
		if err := p.ReadCloseTag(); err != nil {
			return nil, err
		}
		return nil, nil
	}

	if factory, ok := CodeFactories[strings.ToLower(token)]; ok {
		code, err = factory(p)
	} else if token == "}" {
		code = endScope
	} else {
		p.Backwards(length)
		code, err = p.ReadAssignment()
	}

	if err != nil {
		return nil, err
	}
	if err = p.ReadCloseTag(); err != nil {
		return nil, err
	}
	return code, nil
}

type EndScope struct{}

func (c *EndScope) Execute(context *core.Context) core.ExecutionState {
	panic("Execute called on EndScope tag")
}

func (c *EndScope) IsCodeContainer() bool {
	panic("IsCodeContainer called on EndScope tag")
}

func (c *EndScope) IsContentContainer() bool {
	panic("IsContentContainer called on EndScope tag")
}
func (c *EndScope) AddExecutable(e core.Executable) {
	panic("AddExecutable called on EndScope tag")
}
