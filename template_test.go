package gerb

import (
	"bytes"
	"fmt"
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/gerb/core"
	"testing"
)

func Test_RendersALiteral(t *testing.T) {
	assertRender(t, ` hello world `, ` hello world `)
}

func Test_RendersAnIntegerOutput(t *testing.T) {
	assertRender(t, ` <%= 9001 %>`, ` 9001`)
}

func Test_RendersAnFloatOutput(t *testing.T) {
	assertRender(t, `<%= 123.45 %>`, `123.45`)
}

func Test_RendersACharOutput(t *testing.T) {
	assertRender(t, `<%= '!' %>`, `!`)
	assertRender(t, `<%= '\'' %>`, `'`)
	assertRender(t, `<%= '\\' %>`, `\`)
}

func Test_RendersAStringOutput(t *testing.T) {
	assertRender(t, `<%= "it's over" %> 9000`, `it's over 9000`)
	assertRender(t, `<%= "it's \"over\"" %> 9000`, `it's "over" 9000`)
}

func Test_BasicIntegerOperations(t *testing.T) {
	assertRender(t, `<%= 9000 + 1 %>`, `9001`)
	assertRender(t, `<%= 9000 - 1 %>`, `8999`)
	assertRender(t, `<%= 9000 * 2 %>`, `18000`)
	assertRender(t, `<%= 9000 / 2 %>`, `4500`)
	assertRender(t, `<%= 9000 % 7 %>`, `5`)
}

func Test_BasicFloatOperations(t *testing.T) {
	assertRender(t, `<%= 9000.1 + 1 %>`, `9001.1`)
	assertRender(t, `<%= 9000.1 + 1.1 %>`, `9001.2`)
	assertRender(t, `<%= 9000.2 - 1 %>`, `8999.2`)
	assertRender(t, `<%= 9000.2 - 1.05 %>`, `8999.150000000001`)
	assertRender(t, `<%= 9000.3 * 2 %>`, `18000.6`)
	assertRender(t, `<%= 9000.3 * 2.2 %>`, `19800.66`)
	assertRender(t, `<%= 9000.4 / 2 %>`, `4500.2`)
	assertRender(t, `<%= 9000.4 / 2.3 %>`, `3913.217391304348`)
}

func Test_RendersAVariableFromAMap(t *testing.T) {
	assertRender(t, `<%= count %>`, `44`)
	assertRender(t, `<%= count * count %>`, `1936`)
	assertRender(t, `<%= count / count %>`, `1`)
}

func Test_RendersAnObjectsFields(t *testing.T) {
	assertRender(t, `<%= user.Name %>`, `Goku`)
}

func Test_RendersAnNestedObjectsFields(t *testing.T) {
	// assertRender(t, `<%= user.sensei %>`, `Roshi`)
	assertRender(t, `<%= user.sensei.name %>`, `Roshi`)
}

func Test_RendersAnObjectsMethod(t *testing.T) {
	assertRender(t, `<%= user.Analysis(9000) %>`, `it's over 9000!`)
}

func Test_RenderSlices(t *testing.T) {
	assertRender(t, `<%= user.Name[0:3] %>`, `Gok`)
	assertRender(t, `<%= user.Name[1:3] %>`, `ok`)
	assertRender(t, `<%= user.Name[:2] %>`, `Go`)
	assertRender(t, `<%= user.Name[3:] %>`, `u`)
}

func Test_RenderSliceOfMethodReturn(t *testing.T) {
	assertRender(t, `<%= user.Analysis(count)[0:5] %>`, `it's `)
	assertRender(t, `<%= user.Analysis(count)[10:] %>`, `44!`)
}

func Test_UsesBuiltIns(t *testing.T) {
	assertRender(t, `<%= len(user.name) %>`, `3`)
}

func Test_UsesCustomBuiltIns(t *testing.T) {
	core.RegisterBuiltin("add", func(a, b int) int {return a + b})
	assertRender(t, `<%= add(user.powerlevel, 10) %>`, `9010`)
}

func assertRender(t *testing.T, raw, expected string) {
	spec := gspec.New(t)
	template, err := ParseString(raw, false)
	spec.Expect(err).ToBeNil()

	data := map[string]interface{}{
		"count": 44,
		"user":  &Sayan{"Goku", 9001, &Sayan{"Roshi", 0, nil}},
	}
	buffer := new(bytes.Buffer)
	template.Render(buffer, data)
	spec.Expect(buffer.String()).ToEqual(expected)
}

type Sayan struct {
	Name       string
	PowerLevel int
	Sensei     *Sayan
}

func (s *Sayan) Analysis(cutoff int) string {
	if s.PowerLevel > cutoff {
		return fmt.Sprintf("it's over %d!", cutoff)
	}
	return fmt.Sprintf("it's under %d", cutoff)
}

func (s *Sayan) String() string {
	return s.Name
}
