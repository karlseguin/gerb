package gerb

import (
	"fmt"
	"github.com/karlseguin/gerb/core"
)

func AssignmentFactory(p *core.Parser, name string) (core.Code, error) {
	a := &AssignmentCode{names: make([]string, 0, 2), definition: false}
	a.names = append(a.names, name)
	c := p.SkipSpaces()
	if c == ',' {
		names, err := p.ReadTokenList()
		if err != nil {
			return nil, err
		}
		if len(names) == 0 {
			return nil, p.Error("A comma (,) in assignment list should be followed by a variable names")
		}
		a.names = append(a.names, names...)
		c = p.SkipSpaces()
	}

	if c == ':' {
		a.definition = true
		c = p.Next()
	}

	if c != '=' {
		return nil, p.Error("Invalid assignment, expecting = or :=")
	}

	p.Next()
	value, err := p.ReadValue()
	if err != nil {
		return nil, err
	}

	a.nameCount = len(a.names)
	a.value = value
	return a, nil
}

type AssignmentCode struct {
	names      []string
	nameCount  int
	value      core.Value
	definition bool
}

func (c *AssignmentCode) Execute(context *core.Context) core.ExecutionState {
	values := c.value.ResolveAll(context)
	valueCount := len(values)
	if valueCount != c.nameCount {
		core.Log.Error(fmt.Sprintf("%d variable(s) on left side of assignment, but %d return valued received ", c.nameCount, valueCount))
		if c.nameCount < valueCount {
			valueCount = c.nameCount
		}
	}
	for i := 0; i < valueCount; i++ {
		context.Data[c.names[i]] = values[i]
	}
	return core.NormalState
}
