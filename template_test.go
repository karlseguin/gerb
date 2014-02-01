package gerb

import (
	"bytes"
	"fmt"
	"github.com/karlseguin/gerb/core"
	"github.com/karlseguin/gspec"
	"testing"
)

func Test_RendersALiteral(t *testing.T) {
	assertRender(t, ` hello world `, ` hello world `)
}

func Test_RendersAnIntegerOutput(t *testing.T) {
	assertRender(t, ` <%= 9001%>`, ` 9001`)
}

func Test_RendersAnFloatOutput(t *testing.T) {
	assertRender(t, `<%= 123.45 %>`, `123.45`)
}

func Test_RendersBoolean(t *testing.T) {
	assertRender(t, `<%= true %>`, `true`)
	assertRender(t, `<%= false %>`, `false`)
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
	assertRender(t, `<%= 9000.1 + 1%>`, `9001.1`)
	assertRender(t, `<%= 9000.1 + 1.1 %>`, `9001.2`)
	assertRender(t, `<%= 9000.2 - 1 %>`, `8999.2`)
	assertRender(t, `<%= 9000.2 - 1.05 %>`, `8999.150000000001`)
	assertRender(t, `<%= 9000.3 * 2 %>`, `18000.6`)
	assertRender(t, `<%= 9000.3 * 2.2 %>`, `19800.66`)
	assertRender(t, `<%= 9000.4 / 2 %>`, `4500.2`)
	assertRender(t, `<%= 9000.4 / 2.3 %>`, `3913.217391304348`)
}

func Test_UnaryOperations(t *testing.T) {
	assertRender(t, `<%= count++ %> <%= count++ %>`, `45 46`)
	assertRender(t, `<%= count-- %> <%= count-- %><%= count-- %>`, `43 4241`)
}

func Test_UnaryXEqualOperation(t *testing.T) {
	assertRender(t, `<%= count += 5 %> <%= count += jump %>`, `49 51`)
	assertRender(t, `<%= count -= 5 %> <%= count -= jump %>`, `39 37`)
}

func Test_RendersAVariableFromAMap(t *testing.T) {
	assertRender(t, `<%= count %>`, `44`)
	assertRender(t, `<%= count * count %>`, `1936`)
	assertRender(t, `<%= count / count %>`, `1`)
}

func Test_RendersAnObjectsFields(t *testing.T) {
	assertRender(t, `<%=   user.Name%>`, `Goku`)
}

func Test_RendersAnNestedObjectsFields(t *testing.T) {
	assertRender(t, `<%= user.sensei %>`, `Roshi`)
	assertRender(t, `<%= user.sensei.name %>`, `Roshi`)
}

func Test_RendersAnObjectsMethod(t *testing.T) {
	assertRender(t, `<%= user.Analysis(9000) %>`, `it's over 9000!`)
	assertRender(t, `<%= user.Master().Analysis(1) %>`, `it's under 1`)
}

func Test_RenderSlices(t *testing.T) {
	assertRender(t, `<%= user.Name[0:3]%>`, `Gok`)
	assertRender(t, `<%= user.Name[1:3] %>`, `ok`)
	assertRender(t, ` <%=   user.Name[:2]   %>`, ` Go`)
	assertRender(t, `<%= user.Name[3:]%>`, `u`)
	assertRender(t, `<%= user.Name[2]%>`, `k`)
}

func Test_RenderSliceOfMethodReturn(t *testing.T) {
	assertRender(t, `<%= user.Analysis(count)[0:5] %>`, `it's `)
	assertRender(t, `<%= user.Analysis(count)[10:] %>`, `44!`)
}

func Test_UsesBuiltIns(t *testing.T) {
	assertRender(t, `<%= len(user.name) %>`, `4`)
}

func Test_UsesCustomBuiltIns(t *testing.T) {
	core.RegisterBuiltin("add", func(a, b int) int { return a + b })
	assertRender(t, `<%= add(user.powerlevel,10) %>`, `9011`)
}

func Test_UsesPreRegisteredPackages(t *testing.T) {
	assertRender(t, `<%= strings.ToUpper(user.name) %>`, "GOKU")
	assertRender(t, `<%= strings.IndexByte(user.name, 'O') %>`, "-1")
	assertRender(t, `<%= strings.IndexByte(strings.ToUpper(user.name), 'O') %>`, "1")
}

func Test_SingleValueAssigment(t *testing.T) {
	assertRender(t, `<% abc = "123" %><%= abc %>`, "123")
	assertRender(t, `<% c = count %><%= count %>`, "44")
	assertRender(t, `<% l = len(user.Name) %><%= l %>`, "4")
}

func Test_MultipleValueAssigment(t *testing.T) {
	assertRender(t, `<% abc,xyz = "123",987 %><%= abc %> <%= xyz %>`, "123 987")
	assertRender(t, `<% n,err = strconv.Atoi("abc") %><%= n %> - <%= err %>`, `0 - strconv.ParseInt: parsing "abc": invalid syntax`)
}

func Test_IfBool(t *testing.T) {
	assertRender(t, `<% if true { %>1<% } %>`, "1")
	assertRender(t, `<%if true { %>2<%}%>`, "2")
	assertRender(t, `<%if t { %>3<%}%>`, "3")
}

func Test_IfIntTrue(t *testing.T) {
	assertRender(t, `<% if 123 == 123 { %>3<% }%>`, "3")
	assertRender(t, `<% if 0 != 123 { %>4<% }%>`, "4")
	assertRender(t, `<% if 124 > 123 { %>5<% }%>`, "5")
	assertRender(t, `<% if 123 >= 123 { %>6<% }%>`, "6")
	assertRender(t, `<% if 125 >= 123 { %>7<% }%>`, "7")
	assertRender(t, `<% if 122 < 123 { %>8<% }%>`, "8")
	assertRender(t, `<% if 123 <= 123 { %>9<% }%>`, "9")
	assertRender(t, `<% if 122 <= 123 { %>a<% }%>`, "a")
}

func Test_IfIntFalse(t *testing.T) {
	assertRender(t, `<% if 124 == 123 { %>fail<% }%>3`, "3")
	assertRender(t, `<% if 123 != 123 { %>fail<% }%>4`, "4")
	assertRender(t, `<% if 123 > 123 { %>fail<% }%>5`, "5")
	assertRender(t, `<% if 122 >= 123 { %>fail<% }%>7`, "7")
	assertRender(t, `<% if 123 < 123 { %>fail<% }%>8`, "8")
	assertRender(t, `<% if 124 <= 123 { %>fail<% }%>9`, "9")
}

func Test_IfStringTrue(t *testing.T) {
	assertRender(t, `<% if "abc" == "abc" { %>a<% }%>`, "a")
	assertRender(t, `<% if user.name != "vegeta" { %>b<% }%>`, "b")
	assertRender(t, `<% if user.name == "Goku" { %>c<% }%>`, "c")
}

func Test_IfWithMultipleConditions(t *testing.T) {
	assertRender(t, `<% if true && false { %>fail<% }%>a`, "a")
	assertRender(t, `<% if true || false { %>b<% }%>b`, "bb")
}

func Test_IfAssignment(t *testing.T) {
	assertRender(t, `<% if count = 45; count == 45 { %> yes <% }%> <%= count %>`, " yes  45")
	assertRender(t, `<% if t := 22; t == 22 { %> yes <% }%> <%= t %>`, " yes  ")
	assertRender(t, `<% if count+=2; count == 46 { %> yes <% }%> <%= count %>`, " yes  46")
}

func Test_ElseIf(t *testing.T) {
	assertRender(t, `<% if count = 45; count == 44 { %> if <% } else if count == 45 { %> elseif <% } %>`, " elseif ")
	assertRender(t, `<% if false { %> if <% } else if false || true { %> elseif <% } %>`, " elseif ")
}

func Test_ElseIfAssignment(t *testing.T) {
	assertRender(t, `<% if count = 45; count == 44 { %> if <% } else if count = 46; count == 46 { %> elseif <% } %> <%= count %>`, " elseif  46")
	assertRender(t, `<% if count = 45; count == 44 { %> if <% } else if t := 3; true { %> elseif <% } %> <%= t %>`, " elseif  ")
}

func Test_Else(t *testing.T) {
	assertRender(t, `<% if count = 45; count == 44 { %> if <% } else if count == 43 { %> elseif <% } else {%>else<%}%>`, "else")
	assertRender(t, `<% if false { %> if <% } else { %> else <% } %>`, " else ")
}

func Test_InheritanceWithExplicitContent(t *testing.T) {
	assertRender(t, `<% content user.Name { %>contrived<% } %>`, `HEAD <%= yield("Goku") %> FOOTER`, `HEAD contrived FOOTER`)
}

func Test_InheritanceWithImplicitContent(t *testing.T) {
	assertRender(t, `<% content user.Name { %>contrived<% } %>body`, `HEAD <%= yield("Goku") %> <%= yield %> FOOTER`, `HEAD contrived body FOOTER`)
}

func Test_EndlessFor(t *testing.T) {
	assertRender(t, `<% for { %><% if count--; count == 40 { %><% break %><% } else if count == 42 { %><%continue%><% }%> <%=count%> <%}%>`, ` 43  41 `)
}

func Test_NormalFor(t *testing.T) {
	assertRender(t, `<% for i := 0; i < len(user.Name); i++ { %><%= user.Name[i]%> <%}%>`, `G o k u `)
}

func Test_RangedForOverArray(t *testing.T) {
	assertRender(t, `<% for index, score := range scores { %> <%= index %>:<%= score %><% } %>`, ` 0:3 1:10 2:25`)
}

func Test_RangedForOverMap(t *testing.T) {
	assertRender(t, `<% for color, v := range votes { %> <%= color %>:<%= v %><% } %>`, ` red:100 blue:244`)
}

func Test_GroupedTags(t *testing.T) {
	assertRender(t, `<% for {
	if count--; count == 40 {
		break
	} else if count == 42 {
		continue
	} %> <%=count%><% } %> `, ` 43 41 `)
}

func assertRender(t *testing.T, all ...string) {
	expected := all[len(all)-1]
	spec := gspec.New(t)
	template, err := ParseString(false, all[0:len(all)-1]...)
	spec.Expect(err).ToBeNil()

	data := map[string]interface{}{
		"jump":   2,
		"count":  44,
		"t":      true,
		"f":      false,
		"user":   &Sayan{"Goku", 9001, &Sayan{"Roshi", 0, nil}},
		"scores": []int{3, 10, 25},
		"votes":  map[string]int{"red": 100, "blue": 244},
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

func (s *Sayan) Master() *Sayan {
	return s.Sensei
}

func (s *Sayan) String() string {
	return s.Name
}
