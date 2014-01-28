package core

import (
	"testing"
	"time"
)

func Test_EqualsConditionWithBalancedStrings(t *testing.T) {
	assertEqualsCondition(t, true, true, staticValue("abc"), staticValue("abc"))
	assertEqualsCondition(t, true, true, staticValue(""), staticValue(""))
	assertEqualsCondition(t, false, true, staticValue("abc"), staticValue("123"))
}

func Test_EqualsConditionWithBalancedDynamicStrings(t *testing.T) {
	assertEqualsCondition(t, true, true, dynamicValue("doesnotexist"), dynamicValue("doesnotexist"))
	assertEqualsCondition(t, true, true, dynamicValue("string"), staticValue("astring"))
	assertEqualsCondition(t, false, true, dynamicValue("string"), staticValue("other"))
}

func Test_EqualsConditionWithBalancedDynamicArrays(t *testing.T) {
	assertEqualsCondition(t, true, true, dynamicValue("[]int"), dynamicValue("[]int"))
	assertEqualsCondition(t, false, true, dynamicValue("[]int"), dynamicValue("[]int2"))
}

func Test_EqualsConditionWithBalancedBools(t *testing.T) {
	assertEqualsCondition(t, true, true, staticValue(true), staticValue(true))
	assertEqualsCondition(t, false, true, staticValue(true), staticValue(false))
}

func Test_EqualsConditionWithBalancedInt(t *testing.T) {
	assertEqualsCondition(t, true, true, staticValue(3231), staticValue(3231))
	assertEqualsCondition(t, false, true, staticValue(3231), staticValue(2993))
}

func Test_EqualsConditionWithBalancedFloat(t *testing.T) {
	assertEqualsCondition(t, true, true, staticValue(11.33), staticValue(11.33))
	assertEqualsCondition(t, false, true, staticValue(11.2), staticValue(11.21))
}

func Test_EqualWithUnbalancedInt(t *testing.T) {
	assertEqualsCondition(t, false, true, staticValue(123), staticValue("123"))
	assertEqualsCondition(t, false, true, staticValue(123), staticValue("1a23"))
	assertEqualsCondition(t, false, true, staticValue(123), staticValue(123.0))
	assertEqualsCondition(t, false, true, staticValue(123), staticValue(123.1))
}

func Test_EqualWithUnbalancedFloats(t *testing.T) {
	assertEqualsCondition(t, false, true, staticValue(123.0), staticValue("123"))
	assertEqualsCondition(t, false, true, staticValue(123.0), staticValue(123))
	assertEqualsCondition(t, false, true, staticValue(123.0), staticValue("123.1"))
}

func Test_ConditionGroupWithOneCondition(t *testing.T) {
	assertConditionGroup(t, true, trueCondition)
	assertConditionGroup(t, false, falseCondition)
}

func Test_ConditionGroupWithTwoOrCondition(t *testing.T) {
	assertConditionGroup(t, true, trueCondition, OR, trueCondition)
	assertConditionGroup(t, true, trueCondition, OR, falseCondition)
	assertConditionGroup(t, true, falseCondition, OR, trueCondition)
	assertConditionGroup(t, false, falseCondition, OR, falseCondition)
}

func Test_ConditionGroupWithTwoAndCondition(t *testing.T) {
	assertConditionGroup(t, true, trueCondition, AND, trueCondition)
	assertConditionGroup(t, false, trueCondition, AND, falseCondition)
	assertConditionGroup(t, false, falseCondition, AND, trueCondition)
	assertConditionGroup(t, false, falseCondition, AND, falseCondition)
}

func Test_ConditionGroupWithMultipleConditions(t *testing.T) {
	assertConditionGroup(t, true, trueCondition, OR, trueCondition, AND, falseCondition)
	assertConditionGroup(t, true, trueCondition, AND, trueCondition, OR, trueCondition)
	assertConditionGroup(t, false, falseCondition, OR, trueCondition, AND, falseCondition)
	assertConditionGroup(t, false, falseCondition, OR, trueCondition, AND, falseCondition, OR, falseCondition)
	assertConditionGroup(t, true, falseCondition, OR, trueCondition, AND, falseCondition, OR, trueCondition)
}

func assertEqualsCondition(t *testing.T, expected bool, extra bool, left, right Value) {
	assertCondition(t, expected, left, Equals, right)
	assertCondition(t, !expected, left, NotEquals, right)
	if expected && extra {
		assertCondition(t, false, left, LessThan, right)
		assertCondition(t, false, left, GreaterThan, right)
		assertCondition(t, true, left, LessThanOrEqual, right)
		assertCondition(t, true, left, GreaterThanOrEqual, right)
	}
}

func assertCondition(t *testing.T, expected bool, left Value, op ComparisonOperator, right Value) {
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
	actual := c.IsTrue(&Context{Data: data})
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func assertConditionGroup(t *testing.T, expected bool, data ...interface{}) {
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

	actual := group.IsTrue(nil)
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
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
