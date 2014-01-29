package gerb

import (
	"errors"
	"github.com/karlseguin/gerb/core"
)

func IfFactory(p *core.Parser) (core.Code, error) {
	code := &IfCode{
		NormalContainer: new(core.NormalContainer),
		verifiables:     make([]core.Verifiable, 0, 3),
		codes:           make([]core.Code, 0, 3),
	}
	if p.TagContains(';') {
		assignment, err := p.ReadAssignment()
		if err != nil {
			return nil, err
		}
		code.assignment = assignment
		if p.SkipSpaces() != ';' {
			return nil, p.Error("If assignment should be followed by a semicolon")
		}
		p.Next()
	}
	verifiable, err := p.ReadConditionGroup()
	if err != nil {
		return nil, err
	}
	code.verifiables = append(code.verifiables, verifiable)
	code.codes = append(code.codes, code)
	if p.SkipSpaces() != '{' {
		return nil, p.Error("Missing openening brace for if statement")
	}
	p.Next()
	return code, nil
}

type IfCode struct {
	*core.NormalContainer
	assignment  *core.Assignment
	verifiables []core.Verifiable
	codes       []core.Code
}

func (c *IfCode) Execute(context *core.Context) core.ExecutionState {
	if c.assignment != nil {
		if state := c.assignment.Execute(context); state != core.NormalState {
			return state
		}
	}
	for index, verifiable := range c.verifiables {
		if verifiable.IsTrue(context) {
			var state core.ExecutionState
			if index == 0 {
				state = c.NormalContainer.Execute(context)
			} else {
				state = c.codes[index].Execute(context)
			}
			if c.assignment != nil {
				c.assignment.Rollback(context)
			}
			return state
		}
	}

	return core.NormalState
}

func (c *IfCode) IsCodeContainer() bool {
	return true
}

func (c *IfCode) IsContentContainer() bool {
	return true
}

func (c *IfCode) IsSibling() bool {
	return false
}

func (c *IfCode) AddCode(code core.Code) error {
	e, ok := code.(*ElseCode)
	if ok == false {
		return errors.New("If tag only accepts else if/else as a sub tag")
	}
	c.verifiables = append(c.verifiables, e.verifiable)
	c.codes = append(c.codes, e)
	return nil
}

func ElseFactory(p *core.Parser) (core.Code, error) {
	code := &ElseCode{NormalContainer: new(core.NormalContainer)}
	if p.SkipSpaces() == 'i' && p.ConsumeIf([]byte("if")) {
		verifiable, err := p.ReadConditionGroup()
		if err != nil {
			return nil, err
		}
		code.verifiable = verifiable
		if p.SkipSpaces() != '{' {
			return nil, p.Error("Missing openening brace for else if statement")
		}
	} else {
		code.verifiable = core.TrueCondition //else case
		if p.SkipSpaces() != '{' {
			return nil, p.Error("Missing openening brace for else statement")
		}
	}
	p.Next()
	return code, nil
}

type ElseCode struct {
	*core.NormalContainer
	verifiable core.Verifiable
}

func (c *ElseCode) IsCodeContainer() bool {
	return false
}

func (c *ElseCode) IsContentContainer() bool {
	return true
}

func (c *ElseCode) IsSibling() bool {
	return true
}

func (c *ElseCode) AddCode(core.Code) error {
	panic("AddCode called on ElseCode tag")
}
