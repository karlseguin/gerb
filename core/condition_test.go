package core

import (
	t "github.com/karlseguin/expect"
	"testing"
	"time"
)

type ConditionTest struct{}

func Test_Condition(x *testing.T) {
	t.Expectify(new(ConditionTest), x)
}

func (_ ConditionTest) EqualsConditionWithBalancedStrings() {
	assertEqualsCondition(true, true, staticValue("abc"), staticValue("abc"))
	assertEqualsCondition(true, true, staticValue(""), staticValue(""))
	assertEqualsCondition(false, true, staticValue("abc"), staticValue("123"))
}

func (_ ConditionTest) EqualsConditionWithBalancedDynamicStrings() {
	assertEqualsCondition(true, true, dynamicValue("doesnotexist"), dynamicValue("doesnotexist"))
	assertEqualsCondition(true, true, dynamicValue("string"), staticValue("astring"))
	assertEqualsCondition(false, true, dynamicValue("string"), staticValue("other"))
}

func (_ ConditionTest) EqualsConditionWithBalancedDynamicArrays() {
	assertEqualsCondition(true, true, dynamicValue("[]int"), dynamicValue("[]int"))
	assertEqualsCondition(false, true, dynamicValue("[]int"), dynamicValue("[]int2"))
}

func (_ ConditionTest) EqualsConditionWithBalancedBools() {
	assertEqualsCondition(true, true, staticValue(true), staticValue(true))
	assertEqualsCondition(false, true, staticValue(true), staticValue(false))
}

func (_ ConditionTest) EqualsConditionWithBalancedInt() {
	assertEqualsCondition(true, true, staticValue(3231), staticValue(3231))
	assertEqualsCondition(false, true, staticValue(3231), staticValue(2993))
}

func (_ ConditionTest) EqualsConditionWithBalancedFloat() {
	assertEqualsCondition(true, true, staticValue(11.33), staticValue(11.33))
	assertEqualsCondition(false, true, staticValue(11.2), staticValue(11.21))
}

func (_ ConditionTest) EqualWithUnbalancedInt() {
	assertEqualsCondition(false, true, staticValue(123), staticValue("123"))
	assertEqualsCondition(false, true, staticValue(123), staticValue("1a23"))
	assertEqualsCondition(false, true, staticValue(123), staticValue(123.0))
	assertEqualsCondition(false, true, staticValue(123), staticValue(123.1))
}

func (_ ConditionTest) EqualWithUnbalancedFloats() {
	assertEqualsCondition(false, true, staticValue(123.0), staticValue("123"))
	assertEqualsCondition(false, true, staticValue(123.0), staticValue(123))
	assertEqualsCondition(false, true, staticValue(123.0), staticValue("123.1"))
}

func (_ ConditionTest) ConditionGroupWithOneCondition() {
	assertConditionGroup(true, TrueCondition)
	assertConditionGroup(false, FalseCondition)
}

func (_ ConditionTest) ConditionGroupWithTwoOrCondition() {
	assertConditionGroup(true, TrueCondition, OR, TrueCondition)
	assertConditionGroup(true, TrueCondition, OR, FalseCondition)
	assertConditionGroup(true, FalseCondition, OR, TrueCondition)
	assertConditionGroup(false, FalseCondition, OR, FalseCondition)
}

func (_ ConditionTest) ConditionGroupWithTwoAndCondition() {
	assertConditionGroup(true, TrueCondition, AND, TrueCondition)
	assertConditionGroup(false, TrueCondition, AND, FalseCondition)
	assertConditionGroup(false, FalseCondition, AND, TrueCondition)
	assertConditionGroup(false, FalseCondition, AND, FalseCondition)
}

func (_ ConditionTest) ConditionGroupWithMultipleConditions() {
	assertConditionGroup(true, TrueCondition, OR, TrueCondition, AND, FalseCondition)
	assertConditionGroup(true, TrueCondition, AND, TrueCondition, OR, TrueCondition)
	assertConditionGroup(false, FalseCondition, OR, TrueCondition, AND, FalseCondition)
	assertConditionGroup(false, FalseCondition, OR, TrueCondition, AND, FalseCondition, OR, FalseCondition)
	assertConditionGroup(true, FalseCondition, OR, TrueCondition, AND, FalseCondition, OR, TrueCondition)
}

func assertEqualsCondition(expected bool, extra bool, left, right Value) {
	assertCondition(expected, left, Equals, right)
	assertCondition(!expected, left, NotEquals, right)
	if expected && extra {
		assertCondition(false, left, LessThan, right)
		assertCondition(false, left, GreaterThan, right)
		assertCondition(true, left, LessThanOrEqual, right)
		assertCondition(true, left, GreaterThanOrEqual, right)
	}
}

func assertCondition(expected bool, left Value, op ComparisonOperator, right Value) {
	data := map[string]interface{}{
		"[]int":         []int{1, 2, 3},
		"[]int2":        []int{2, 3, 1},
		"[]int3":        []int{},
		"[]interface{}": []interface{}{2, "a", true},
		"string":        "astring",
		"now":           time.Now(),
		"yesterday":     time.Now().Add(time.Hour * -24),
		"map[string]int": map[string]int{
			"hello": 44,
		},
	}
	c := &Condition{left, op, right}
	t.Expect(c.IsTrue(&Context{Data: data})).To.Equal(expected)
}

func assertConditionGroup(expected bool, data ...interface{}) {
	l := len(data)
	group := &ConditionGroup{
		joins:       make([]LogicalOperator, 0, l/2),
		verifiables: make([]Verifiable, 0, l-l/2),
	}
	for i := 0; i < l; i += 2 {
		group.verifiables = append(group.verifiables, data[i].(Verifiable))
		if i+1 < l {
			group.joins = append(group.joins, data[i+1].(LogicalOperator))
		}
	}

	t.Expect(group.IsTrue(nil)).To.Equal(expected)
}

func staticValue(v interface{}) Value {
	return &StaticValue{v}
}

func dynamicValue(s string) Value {
	return &DynamicValue{
		id:    "",
		names: []string{s},
		types: []DynamicFieldType{FieldType},
		args:  [][]Value{[]Value{}},
	}
}
