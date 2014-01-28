package gerb

import (
	"github.com/karlseguin/gerb/core"
)

func IfFactory(p *core.Parser) (core.Code, error) {
	code := &IfCode{NormalContainer: new(core.NormalContainer)}
	if p.TagContains(';') {
		println("a")
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
	code.verifiable = verifiable
	if p.SkipSpaces() != '{' {
		return nil, p.Error("Missing openening brace for if statement")
	}
	p.Next()
	return code, nil
}

type IfCode struct {
	*core.NormalContainer
	assignment *core.Assignment
	verifiable core.Verifiable
}

func (c *IfCode) Execute(context *core.Context) core.ExecutionState {
	//todo if assignment.definition == true, we should rollback this assignment
	if c.assignment != nil {
		if state := c.assignment.Execute(context); state != core.NormalState {
			return state
		}
	}
	if c.verifiable.IsTrue(context) {
		return c.NormalContainer.Execute(context)
	}
	return core.NormalState
}

func (c *IfCode) IsCodeContainer() bool {
	return true
}

func (c *IfCode) IsContentContainer() bool {
	return true
}
