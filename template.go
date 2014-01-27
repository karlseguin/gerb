package gerb

import (
	"github.com/karlseguin/gerb/core"
	"io"
)

func newTemplate(data []byte) (*Template, error) {
	template := &Template{new(core.NormalContainer)}
	var container core.Container = template
	parser := core.NewParser(data)
	for {
		if literal := parser.ReadLiteral(); literal != nil {
			container.AddExecutable(literal)
		}
		tagType := parser.ReadTagType()
		if tagType == core.NoTag {
			return template, nil
		}
		if tagType == core.CodeTag {
			code, err := createCodeTag(parser)
			if err != nil {
				return nil, err
			}
			if code != nil {
				container.AddExecutable(code)
			}
		}

		isUnsafe := tagType == core.UnsafeTag
		if tagType == core.OutputTag || isUnsafe {
			output, err := createOutputTag(parser, isUnsafe)
			if err != nil {
				return nil, err
			}
			if output != nil {
				container.AddExecutable(output)
			}
		}
	}
	return template, nil
}

type Template struct {
	*core.NormalContainer
}

func (t *Template) Render(writer io.Writer, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	context := &core.Context{
		Writer:   writer,
		Data:     data,
		Counters: make(map[string]int),
	}
	t.Execute(context)
}
