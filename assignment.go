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
		p.Next()
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
	values, err := p.ReadValueList()
	if err != nil {
		return nil, err
	}
	a.values = values
	return a, nil
}

type AssignmentCode struct {
	names      []string
	values     []core.Value
	definition bool
}

func (c *AssignmentCode) Execute(context *core.Context) core.ExecutionState {
	index := 0
	hasNew := false
	remaining := len(c.names)
	for _, value := range c.values {
		values := value.ResolveAll(context)
		valueCount := len(values)
		if valueCount > remaining {
			core.Log.Error(fmt.Sprintf("%d more return value than there are variables", valueCount-remaining))
			valueCount = remaining
		}
		for i := 0; i < valueCount; i++ {
			name := c.names[index]
			index++
			remaining--
			if _, exists := context.Data[name]; !exists {
				hasNew = true
			}
			context.Data[name] = values[i]
		}
	}

	if remaining > 0 {
		core.Log.Error(fmt.Sprintf("Expected %d variable(s) but only got %d", len(c.names), len(c.names)-remaining))
	}

	if hasNew && !c.definition {
		core.Log.Error(fmt.Sprintf("Assigning to %v, which are undefined, using =", c.names))
	} else if !hasNew && c.definition {
		core.Log.Error(fmt.Sprintf("Assigning to %v, which are already defined, using :=", c.names))
	}
	return core.NormalState
}
