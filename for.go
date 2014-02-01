package gerb

import (
	"errors"
	"fmt"
	"github.com/karlseguin/gerb/core"
)

func ForFactory(p *core.Parser) (core.Code, error) {
	if p.TagContains(';') {
		return ExplicitForFactory(p)
	}
	return RangedForFactory(p)
}

func ExplicitForFactory(p *core.Parser) (core.Code, error) {
	code := &ForCode{NormalContainer: new(core.NormalContainer)}
 	if p.SkipSpaces() != ';' {
 		assignment, err := p.ReadAssignment()
		if err != nil {
			return nil, err
		}
		code.init = assignment
 	}

	if p.SkipSpaces() != ';' {
		return nil, p.Error("Invalid for loop, expecting INIT; CONDITION; STEP (1)")
	}
	p.Next()

 	verifiable, err := p.ReadConditionGroup(false)
 	if err != nil {
 		return nil, err
 	}

 	code.verifiable = verifiable

	if p.SkipSpaces() != ';' {
		return nil, p.Error("Invalid for loop, expecting INIT; CONDITION; STEP (1)")
	}
	p.Next()

 	if p.SkipSpaces() != '{' {
 		value, err := p.ReadAssignment()
 		if err != nil {
 			return nil, err
 		}
 		code.step = value
 	}
	if p.SkipSpaces() != '{' {
		return nil, p.Error("Missing openening brace for if statement")
	}
	p.Next()
 	return code, nil
}

func RangedForFactory(p *core.Parser) (core.Code, error) {
	return nil, nil
}

type ForCode struct {
	*core.NormalContainer
	init *core.Assignment
	verifiable core.Verifiable
	step *core.Assignment
}

func (c *ForCode) Execute(context *core.Context) core.ExecutionState {
	if c.init != nil {
		c.init.Execute(context)
	}
	for {
		if c.verifiable.IsTrue(context) == false {
			break
		}
		c.NormalContainer.Execute(context)
		if c.step != nil {
			c.step.Execute(context)
		}
	}
	return core.NormalState
}

func (c *ForCode) IsCodeContainer() bool {
	return true
}

func (c *ForCode) IsContentContainer() bool {
	return true
}

func (c *ForCode) IsSibling() bool {
	return false
}

func (c *ForCode) AddCode(code core.Code) error {
	return errors.New(fmt.Sprintf("%v is not a valid tag as a descendant of a for loop", code))
}
