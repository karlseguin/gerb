package core

import (
	t "github.com/karlseguin/expect"
	"testing"
)

type ParserTest struct{}

func Test_Parser(x *testing.T) {
	t.Expectify(new(ParserTest), x)
}

func (_ ParserTest) ReadsAIntegerValue() {
	assertStaticValue(" 123 %>", 123)
}

func (_ ParserTest) ReadsANegativeIntegerValue() {
	assertStaticValue("-9944 %>", -9944)
}

func (_ ParserTest) ReadsAFloatValue() {
	assertStaticValue(" 123.2334 %>", 123.2334)
}

func (_ ParserTest) ReadsANegativeFloatValue() {
	assertStaticValue("   -9944.991338 %>", -9944.991338)
}

func (_ ParserTest) ErrorsReadingANegativeCharValue() {
	assertErrorValue("   -'b' %>", "Don't know what to do with a negative character:    -'b' %>")
}

func (_ ParserTest) ErrorsReadingANegativeStringValue() {
	assertErrorValue(`   -"over9000" %>`, `Don't know what to do with a negative string:    -"over9000" %>`)
}

func (_ ParserTest) ReadsASimpleStringValue() {
	assertStaticValue(` "it's over 9000" %>`, `it's over 9000`)
	assertStaticValue(" `it's over 9000` %>", `it's over 9000`)
}

func (_ ParserTest) ReadsAnEscapedStringValue() {
	assertStaticValue(` "it's \n \\\"over 9000" %>`, "it's \n \\\"over 9000")
	assertStaticValue(" `ab\\n\\c` %>", "ab\\n\\c")
}

func (_ ParserTest) ErrorReadingStringWithUnknownEscapeSequence() {
	assertErrorValue(` "what\zs" %>`, `Unknown escape sequence \z:  "what\zs" %>`)
}

func assertStaticValue(data string, expected interface{}) {
	p := NewParser([]byte(data))
	value, err := p.ReadValue()
	t.Expect(err).To.Equal(nil)
	t.Expect(value.Resolve(nil)).To.Equal(expected)
}

func assertErrorValue(data string, expected string) {
	p := NewParser([]byte(data))
	value, err := p.ReadValue()
	t.Expect(err.Error()).To.Equal(expected)
	t.Expect(value).To.Equal(nil)
}
