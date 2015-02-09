package gerb

import (
	"bytes"
	"fmt"
	. "github.com/karlseguin/expect"
	"github.com/karlseguin/gerb/core"
	"testing"
)

type TemplateTest struct{}

func Test_Template(t *testing.T) {
	Expectify(new(TemplateTest), t)
}

func (_ TemplateTest) RendersALiteral() {
	assertRender(` hello world `, ` hello world `)
}

func (_ TemplateTest) RendersAnIntegerOutput() {
	assertRender(` <%= 9001%>`, ` 9001`)
}

func (_ TemplateTest) RendersAnFloatOutput() {
	assertRender(`<%= 123.45 %>`, `123.45`)
}

func (_ TemplateTest) RendersBoolean() {
	assertRender(`<%= true %>`, `true`)
	assertRender(`<%= false %>`, `false`)
}

func (_ TemplateTest) RendersACharOutput() {
	assertRender(`<%= '!' %>`, `!`)
	assertRender(`<%= '\'' %>`, `'`)
	assertRender(`<%= '\\' %>`, `\`)
}

func (_ TemplateTest) RendersAStringOutput() {
	assertRender(`<%= "it's over" %> 9000`, `it's over 9000`)
	assertRender(`<%= "it's \"over\"" %> 9000`, `it's "over" 9000`)
	assertRender("<%= `it's over` %> 9000", "it's over 9000")
}

func (_ TemplateTest) BasicIntegerOperations() {
	assertRender(`<%= 9000 + 1 %>`, `9001`)
	assertRender(`<%= 9000 - 1 %>`, `8999`)
	assertRender(`<%= 9000 * 2 %>`, `18000`)
	assertRender(`<%= 9000 / 2 %>`, `4500`)
	assertRender(`<%= 9000 % 7 %>`, `5`)
}

func (_ TemplateTest) BasicFloatOperations() {
	assertRender(`<%= 9000.1 + 1%>`, `9001.1`)
	assertRender(`<%= 9000.1 + 1.1 %>`, `9001.2`)
	assertRender(`<%= 9000.2 - 1 %>`, `8999.2`)
	assertRender(`<%= 9000.2 - 1.05 %>`, `8999.150000000001`)
	assertRender(`<%= 9000.3 * 2 %>`, `18000.6`)
	assertRender(`<%= 9000.3 * 2.2 %>`, `19800.66`)
	assertRender(`<%= 9000.4 / 2 %>`, `4500.2`)
	assertRender(`<%= 9000.4 / 2.3 %>`, `3913.217391304348`)
}

func (_ TemplateTest) UnaryOperations() {
	assertRender(`<%= count++ %> <%= count++ %>`, `45 46`)
	assertRender(`<%= count-- %> <%= count-- %><%= count-- %>`, `43 4241`)
}

func (_ TemplateTest) UnaryXEqualOperation() {
	assertRender(`<%= count += 5 %> <%= count += jump %>`, `49 51`)
	assertRender(`<%= count -= 5 %> <%= count -= jump %>`, `39 37`)
}

func (_ TemplateTest) RendersAVariableFromAMap() {
	assertRender(`<%= count %>`, `44`)
	assertRender(`<%= count * count %>`, `1936`)
	assertRender(`<%= count / count %>`, `1`)
}

func (_ TemplateTest) RendersAnObjectsFields() {
	assertRender(`<%=   user.Name%>`, `Goku`)
}

func (_ TemplateTest) RendersAnNestedObjectsFields() {
	assertRender(`<%= user.sensei %>`, `Roshi`)
	assertRender(`<%= user.sensei.name %>`, `Roshi`)
}

func (_ TemplateTest) RendersAnObjectsMethod() {
	assertRender(`<%= user.Analysis(9000) %>`, `it's over 9000!`)
	assertRender(`<%= user.Master().Analysis(1) %>`, `it's under 1`)
}

func (_ TemplateTest) RenderSlices() {
	assertRender(`<%= user.Name[0:3]%>`, `Gok`)
	assertRender(`<%= user.Name[1:3] %>`, `ok`)
	assertRender(` <%=   user.Name[:2]   %>`, ` Go`)
	assertRender(`<%= user.Name[3:]%>`, `u`)
	assertRender(`<%= user.Name[2]%>`, `k`)
}

func (_ TemplateTest) RendersAMapValue() {
	assertRender(`<%= other["a"]%>`, `1`)
	assertRender("<%= other[`b`]%>", `2`)
}

func (_ TemplateTest) RenderSliceOfMethodReturn() {
	assertRender(`<%= user.Analysis(count)[0:5] %>`, `it's `)
	assertRender(`<%= user.Analysis(count)[10:] %>`, `44!`)
}

func (_ TemplateTest) UsesBuiltIns() {
	assertRender(`<%= len(user.name) %>`, `4`)
}

func (_ TemplateTest) UsesCustomBuiltIns() {
	core.RegisterBuiltin("add", func(a, b int) int { return a + b })
	assertRender(`<%= add(user.powerlevel,10) %>`, `9011`)
}

func (_ TemplateTest) UsesPreRegisteredPackages() {
	assertRender(`<%= strings.ToUpper(user.name) %>`, "GOKU")
	assertRender(`<%= strings.IndexByte(user.name, 'O') %>`, "-1")
	assertRender(`<%= strings.IndexByte(strings.ToUpper(user.name), 'O') %>`, "1")
}

func (_ TemplateTest) SingleValueAssigment() {
	assertRender(`<% abc := "123" %><%= abc %>`, "123")
	assertRender(`<% c := count %><%= count %>`, "44")
	assertRender(`<% l := len(user.Name) %><%= l %>`, "4")
}

func (_ TemplateTest) MultipleValueAssigment() {
	assertRender(`<% abc,xyz := "123",987 %><%= abc %> <%= xyz %>`, "123 987")
	assertRender(`<% n,err := strconv.Atoi("abc") %><%= n %> - <%= err %>`, `0 - strconv.ParseInt: parsing "abc": invalid syntax`)
}

func (_ TemplateTest) IfBool() {
	assertRender(`<% if true { %>1<% } %>`, "1")
	assertRender(`<%if true { %>2<%}%>`, "2")
	assertRender(`<%if t { %>3<%}%>`, "3")
}

func (_ TemplateTest) IfIntTrue() {
	assertRender(`<% if 123 == 123 { %>3<% }%>`, "3")
	assertRender(`<% if 0 != 123 { %>4<% }%>`, "4")
	assertRender(`<% if 124 > 123 { %>5<% }%>`, "5")
	assertRender(`<% if 123 >= 123 { %>6<% }%>`, "6")
	assertRender(`<% if 125 >= 123 { %>7<% }%>`, "7")
	assertRender(`<% if 122 < 123 { %>8<% }%>`, "8")
	assertRender(`<% if 123 <= 123 { %>9<% }%>`, "9")
	assertRender(`<% if 122 <= 123 { %>a<% }%>`, "a")
}

func (_ TemplateTest) IfIntFalse() {
	assertRender(`<% if 124 == 123 { %>fail<% }%>3`, "3")
	assertRender(`<% if 123 != 123 { %>fail<% }%>4`, "4")
	assertRender(`<% if 123 > 123 { %>fail<% }%>5`, "5")
	assertRender(`<% if 122 >= 123 { %>fail<% }%>7`, "7")
	assertRender(`<% if 123 < 123 { %>fail<% }%>8`, "8")
	assertRender(`<% if 124 <= 123 { %>fail<% }%>9`, "9")
}

func (_ TemplateTest) IfStringTrue() {
	assertRender(`<% if "abc" == "abc" { %>a<% }%>`, "a")
	assertRender(`<% if user.name != "vegeta" { %>b<% }%>`, "b")
	assertRender(`<% if user.name == "Goku" { %>c<% }%>`, "c")
}

func (_ TemplateTest) IfWithMultipleConditions() {
	assertRender(`<% if true && false { %>fail<% }%>a`, "a")
	assertRender(`<% if true || false { %>b<% }%>b`, "bb")
}

func (_ TemplateTest) IfAssignment() {
	assertRender(`<% if count = 45; count == 45 { %> yes <% }%> <%= count %>`, " yes  45")
	assertRender(`<% if ttt := 22; ttt == 22 { %> yes <% }%> <%= ttt %>`, " yes  ")
	assertRender(`<% if count+=2; count == 46 { %> yes <% }%> <%= count %>`, " yes  46")
}

func (_ TemplateTest) ElseIf() {
	assertRender(`<% if count = 45; count == 44 { %> if <% } else if count == 45 { %> elseif <% } %>`, " elseif ")
	assertRender(`<% if false { %> if <% } else if false || true { %> elseif <% } %>`, " elseif ")
}

func (_ TemplateTest) ElseIfAssignment() {
	assertRender(`<% if count = 45; count == 44 { %> if <% } else if count = 46; count == 46 { %> elseif <% } %> <%= count %>`, " elseif  46")
	assertRender(`<% if count = 45; count == 44 { %> if <% } else if ttt := 3; true { %> elseif <% } %> `, " elseif  ")
}

func (_ TemplateTest) Else() {
	assertRender(`<% if count = 45; count == 44 { %> if <% } else if count == 43 { %> elseif <% } else {%>else<%}%>`, "else")
	assertRender(`<% if false { %> if <% } else { %> else <% } %>`, " else ")
}

func (_ TemplateTest) InheritanceWithExplicitContent() {
	assertRender(`<% content user.Name { %>contrived<% } %>`, `HEAD <%= yield("Goku") %> FOOTER`, `HEAD contrived FOOTER`)
}

func (_ TemplateTest) InheritanceWithImplicitContent() {
	assertRender(`<% content user.Name { %>contrived<% } %>body`, `HEAD <%= yield("Goku") %> <%= yield %> FOOTER`, `HEAD contrived body FOOTER`)
}

func (_ TemplateTest) EndlessFor() {
	assertRender(`<% for { %><% if count--; count == 40 { %><% break %><% } else if count == 42 { %><%continue%><% }%> <%=count%> <%}%>`, ` 43  41 `)
}

func (_ TemplateTest) NormalFor() {
	assertRender(`<% for i := 0; i < len(user.Name); i++ { %><%= user.Name[i]%> <%}%>`, `G o k u `)
}

func (_ TemplateTest) RangedForOverArray() {
	assertRender(`<% for index, score := range scores { %> <%= index %>:<%= score %><% } %>`, ` 0:3 1:10 2:25`)
}

func (_ TemplateTest) DurationHandling() {
	assertRender(`<%= time.Unix(1410569706, 0).Sub(time.Unix(1410569700, 0)).Seconds() %>`, `6`)
}

func (_ TemplateTest) StripNewlines() {
	input := `
<%% if true { %%>
	abc
<%% } %%>
`
	assertRender(input, `	abc`)
}

func (_ TemplateTest) StripLeadNewlines() {
	input := `<%% if true { %%>
	abc
<%% } %%>
`
	assertRender(input, `	abc`)
}

func (_ TemplateTest) RangedForOverMap() {
	template, err := ParseString(false, `<% for color, v := range votes { %> <%= color %>:<%= v %><% } %>`)
	Expect(err).To.Equal(nil)

	buffer := new(bytes.Buffer)
	template.Render(buffer, map[string]interface{}{
		"votes": map[string]int{"red": 100, "blue": 244},
	})
	if buffer.String() != ` red:100 blue:244` && buffer.String() != ` blue:244 red:100` {
		Fail("expecting output to be ' red:100 blue:244'")
	}
}

func (_ TemplateTest) GroupedTags() {
	assertRender(`<% for {
	if count--; count == 40 {
		break
	} else if count == 42 {
		continue
	} %> <%=count%><% } %> `, ` 43 41 `)
}

func (_ TemplateTest) Comments() {
	assertRender(`<%# comment 1 %>
not comment 1
<%# comment 2 %>
not comment 2`, "\nnot comment 1\n\nnot comment 2")
}

func (_ TemplateTest) StrippedComments() {
	assertRender(`<%%# comment 1 %%>
not comment 1
<%%# comment 2 %>
not comment 2`, "not comment 1\nnot comment 2")
}

func (_ TemplateTest) PercentInComment() {
	assertRender(`<%# 4 % 5 %>cm`, "cm")
}

func assertRender(all ...string) {
	expected := all[len(all)-1]
	template, err := ParseString(false, all[0:len(all)-1]...)
	Expect(err).To.Equal(nil)

	data := map[string]interface{}{
		"jump":   2,
		"count":  44,
		"t":      true,
		"f":      false,
		"user":   &Sayan{"Goku", 9001, &Sayan{"Roshi", 0, nil}},
		"scores": []int{3, 10, 25},
		"other":  map[string]int{"a": 1, "b": 2},
	}
	buffer := new(bytes.Buffer)
	template.Render(buffer, data)
	Expect(buffer.String()).To.Equal(expected)
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
