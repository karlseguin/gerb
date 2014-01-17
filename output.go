package gerb

import (
	"github.com/karlseguin/gerb/core"
	"github.com/karlseguin/gerb/r"
)

type OutputTag struct {
	value core.Value
}

func (o *OutputTag) Execute(context *core.Context) core.ExecutionState {
	context.Writer.Write(r.ToBytes(o.value.Resolve(context)))
	return core.NormalState
}

func createOutputTag(p *core.Parser, isUnsafe bool) (core.Code, error) {
	value, err := p.ReadValue()
	if err != nil {
		return nil, err
	}
	if err = p.ReadCloseTag(); err != nil {
		return nil, err
	}
	return &OutputTag{value}, nil
}
