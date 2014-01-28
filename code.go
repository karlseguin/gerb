package gerb

import (
	"github.com/karlseguin/gerb/core"
	"strings"
)

var CodeFactories = map[string]CodeFactory{
// "if": IfFactory,
}

type CodeFactory func(*core.Parser) (core.Code, error)

func createCodeTag(p *core.Parser) (core.Code, error) {
	token, err := p.ReadToken()
	if err != nil {
		return nil, err
	}
	if len(token) == 0 {
		return nil, nil
	}
	if factory, ok := CodeFactories[strings.ToLower(token)]; ok {
		return factory(p)
	}
	code, err := AssignmentFactory(p, token)
	if err != nil {
		return nil, err
	}
	if err = p.ReadCloseTag(); err != nil {
		return nil, err
	}
	return code, nil
}
