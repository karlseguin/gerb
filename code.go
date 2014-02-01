package gerb

import (
	"github.com/karlseguin/gerb/core"
)

var CodeFactories = map[string]CodeFactory{
	"if":       IfFactory,
	"content":  ContentFactory,
	"for":      ForFactory,
	"break":    BreakFactory,
	"continue": ContinueFactory,
}

var endScope = new(EndScope)

type CodeFactory func(*core.Parser) (core.Code, error)

func createCodeTag(p *core.Parser) ([]core.Code, error) {
	codes := make([]core.Code, 0, 1)
	for {
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

		var code core.Code
		if token == "}" {
			if p.SkipSpaces() == 'e' && p.ConsumeIf([]byte("else")) {
				code, err = ElseFactory(p)
			} else {
				code = endScope
			}
		} else if factory, ok := CodeFactories[token]; ok {
			code, err = factory(p)
		} else {
			p.Backwards(length)
			code, err = p.ReadAssignment()
		}
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
		if p.SkipSpaces() == '%' && p.ConsumeIf([]byte("%>")) {
			return codes, nil
		}
	}
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

func (c *EndScope) IsSibling() bool {
	return false
}

func (c *EndScope) AddExecutable(e core.Executable) {
	panic("AddExecutable called on EndScope tag")
}

func (c *EndScope) AddCode(core.Code) error {
	panic("AddCode called on EndScope tag")
}
