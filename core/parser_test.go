package core

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func Test_ReadsAIntegerValue(t *testing.T) {
	assertStaticValue(t, " 123 %>", 123)
}

func Test_ReadsANegativeIntegerValue(t *testing.T) {
	assertStaticValue(t, "-9944 %>", -9944)
}

func Test_ReadsAFloatValue(t *testing.T) {
	assertStaticValue(t, " 123.2334 %>", 123.2334)
}

func Test_ReadsANegativeFloatValue(t *testing.T) {
	assertStaticValue(t, "   -9944.991338 %>", -9944.991338)
}

func Test_ErrorsReadingANegativeCharValue(t *testing.T) {
	assertErrorValue(t, "   -'b' %>", "Don't know what to do with a negative character:    -'b' %>")
}

func Test_ErrorsReadingANegativeStringValue(t *testing.T) {
	assertErrorValue(t, `   -"over9000" %>`, `Don't know what to do with a negative string:    -"over9000" %>`)
}

func Test_ReadsASimpleStringValue(t *testing.T) {
	assertStaticValue(t, ` "it's over 9000" %>`, `it's over 9000`)
}

func Test_ReadsAnEscapedStringValue(t *testing.T) {
	assertStaticValue(t, ` "it's \n \\\"over 9000" %>`, "it's \n \\\"over 9000")
}

func Test_ErrorReadingStringWithUnknownEscapeSequence(t *testing.T) {
	assertErrorValue(t, ` "what\zs" %>`, `Unknown escape sequence \z:  "what\zs" %>`)
}

func assertStaticValue(t *testing.T, data string, expected interface{}) {
	spec := gspec.New(t)
	p := NewParser([]byte(data))
	value, err := p.ReadValue()
	spec.Expect(err).ToBeNil()
	spec.Expect(value.Resolve(nil)).ToEqual(expected)
}

func assertErrorValue(t *testing.T, data string, expected string) {
	spec := gspec.New(t)
	p := NewParser([]byte(data))
	value, err := p.ReadValue()
	spec.Expect(err.Error()).ToEqual(expected)
	spec.Expect(value).ToBeNil()
}
