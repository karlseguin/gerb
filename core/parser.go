package core

import (
	"bytes"
	"errors"
	"fmt"
	"math"
)

type Parser struct {
	end      int
	len      int
	position int
	data     []byte
}

func NewParser(data []byte) *Parser {
	p := &Parser{
		data: data,
		len:  len(data),
		end:  len(data) - 1,
	}
	return p
}

func (p *Parser) ReadLiteral() *Literal {
	start := p.position
	for {
		if p.SkipUntil('%') == false {
			return &Literal{clone(p.data[start:p.len])}
		}
		if p.Prev() == '<' {
			p.position++ //move past the %
			return &Literal{clone(p.data[start : p.position-2])}
		}
	}
}

func (p *Parser) ReadValue() (Value, error) {
	first := p.SkipSpaces()
	negate := false
	if first == '-' {
		negate = true
		p.position++
		first = p.SkipSpaces()
	}
	var value Value
	var err error
	if first == 0 {
		return nil, p.error("Expected value, got nothing")
	}
	if first >= '0' && first <= '9' {
		value, err = p.ReadNumber(negate)
	} else if first == '\'' {
		value, err = p.ReadChar(negate)
	} else if first == '"' {
		value, err = p.ReadString(negate)
	} else {
		value, err = p.ReadDynamic(negate)
	}
	if err != nil {
		return nil, err
	}
	c1 := p.SkipSpaces()
	if c1 == '%' && p.data[p.position+1] == '>' {
		return value, nil
	}
	factory, ok := Operations[c1]
	if ok == false {
		return value, nil
	}
	p.position++
	right, err := p.ReadValue()
	if err != nil {
		return nil, err
	}
	return factory(value, right), nil
}

func (p *Parser) ReadNumber(negate bool) (Value, error) {
	integer := 0
	fraction := 0
	target := &integer
	partLength := 0
	isDecimal := false
	for ; p.position < p.end; p.position++ {
		c := p.data[p.position]
		if c == '.' {
			if isDecimal {
				break
			}
			target = &fraction
			partLength = 0
			isDecimal = true
			continue
		}
		if c < '0' || c > '9' {
			break
		}
		partLength++
		*target = *target*10 + int(c-'0')
	}

	if isDecimal {
		value := float64(integer) + float64(fraction)/math.Pow10(partLength)
		if negate {
			value *= -1
		}
		return &StaticValue{value}, nil
	}
	if negate {
		integer *= -1
	}
	return &StaticValue{integer}, nil
}

func (p *Parser) ReadChar(negate bool) (Value, error) {
	if negate {
		return nil, p.error("Don't know what to do with a negative character")
	}
	c := p.Next()
	if c == '\\' {
		c = p.Next()
	}
	if p.Next() != '\'' {
		return nil, p.error("Invalid character")
	}
	p.position++
	return &StaticValue{c}, nil
}

func (p *Parser) ReadString(negate bool) (Value, error) {
	if negate {
		return nil, p.error("Don't know what to do with a negative string")
	}
	p.position++
	start := p.position
	escaped := 0

	for ; p.position < p.end; p.position++ {
		c := p.data[p.position]
		if c == '\\' {
			escaped++
			p.position++
			continue
		}
		if c == '"' {
			break
		}
	}

	var data []byte
	var err error
	if escaped > 0 {
		data, err = p.unescape(p.data[start:p.position], escaped)
		if err != nil {
			return nil, err
		}
	} else {
		data = p.data[start:p.position]
	}
	p.position++ //consume the "
	return &StaticValue{string(data)}, nil
}

func (p *Parser) ReadDynamic(negate bool) (Value, error) {
	start := p.position
	fields := make([]string, 0, 5)
	types := make([]DynamicFieldType, 0, 5)
	args := make([][]Value, 0, 5)
	for p.position < p.end {
		c := p.data[p.position]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' {
			p.position++
			continue
		}
		// if c == ' ' {
		// 	start++
		// 	p.position++
		// 	continue
		// }
		if start == p.position {
			break
		}
		field := string(bytes.ToLower(p.data[start:p.position]))
		isEnd := c != '.' && c != '(' && c != '['
		if c == '.' || isEnd {
			fields = append(fields, field)
			types = append(types, FieldType)
			args = append(args, nil)
			if isEnd {
				break
			}
			p.position++
		} else if c == '[' {
			fields = append(fields, field)
			types = append(types, IndexedType)
			p.position++
			arg, err := p.ReadIndexing()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		} else if c == '(' {
			fields = append(fields, field)
			types = append(types, MethodType)
			p.position++
			arg, err := p.ReadArgs()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
		start = p.position
	}
	return &DynamicValue{fields, types, args}, nil
}

func (p *Parser) ReadIndexing() ([]Value, error) {
	implicitStart := false
	if p.SkipSpaces() == ':' {
		implicitStart = true
		p.position++
	}
	first, err := p.ReadValue()
	if err != nil {
		return nil, err
	}
	if implicitStart {
		p.position++
		return []Value{&StaticValue{0}, first}, nil
	}

	c := p.SkipSpaces()
	if c == ']' {
		p.position++
		return []Value{first}, nil
	}
	if c != ':' {
		return nil, p.error("Unrecognized array/map index")
	}

	p.position++
	if p.SkipSpaces() == ']' {
		p.position++
		return []Value{first}, nil
	}
	second, err := p.ReadValue()
	if err != nil {
		return nil, err
	}

	if c = p.SkipSpaces(); c != ']' {
		return nil, p.error("Expected closing array/map bracket")
	}
	p.position++
	return []Value{first, second}, nil
}

func (p *Parser) ReadArgs() ([]Value, error) {
	if p.data[p.position] == ')' {
		p.position++
		return nil, nil
	}

	values := make([]Value, 0, 3)
	for {
		value, err := p.ReadValue()
		if err != nil {
			return nil, err
		}
		values = append(values, value)
		c := p.SkipSpaces()
		if c == ')' {
			p.position++
			break
		}
		if c != ',' {
			return nil, p.error("Invalid argument list given to function")
		}
	}
	return values, nil
}

func (p *Parser) ReadTagType() TagType {
	switch p.Consume() {
	case 0:
		return NoTag
	case '=':
		return OutputTag
	case '!':
		return UnsafeTag
	default:
		return NoTag //todo CodeTag
	}
}

func (p *Parser) ReadCloseTag() error {
	if p.SkipSpaces() != '%' || p.Next() != '>' {
		return p.error("Expected closing tag")
	}
	p.position++
	return nil
}

func (p *Parser) SkipUntil(b byte) bool {
	if at := bytes.IndexByte(p.data[p.position:], b); at != -1 {
		p.position = p.position + at
		return true
	}
	p.position = len(p.data)
	return false
}

func (p *Parser) SkipSpaces() byte {
	for ; p.position < p.end; p.position++ {
		c := p.data[p.position]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			return c
		}
	}
	return 0
}

func (p *Parser) Consume() byte {
	if p.position > p.end {
		return 0
	}
	c := p.data[p.position]
	p.position++
	return c
}

func (p *Parser) Next() byte {
	p.position++
	if p.position > p.end {
		return 0
	}
	return p.data[p.position]
}

func (p *Parser) Prev() byte {
	return p.data[p.position-1]
}

func (p *Parser) Dump() {
	fmt.Println(string(p.data[p.position:]))
}

func (p *Parser) error(s string) error {
	end := p.position
	for ; end < p.end; end++ {
		if p.data[end] == '%' && p.data[end+1] == '>' {
			break
		}
	}
	end += 2 //consume the > + this is exclusive
	if end > p.len {
		end = p.len
	}
	start := p.position
	for ; start > 0; start-- {
		if p.data[start] == '%' && p.data[start-1] == '<' {
			start--
			break
		}
	}
	return errors.New(fmt.Sprintf("%s: %v", s, string(p.data[start:end])))
}

func (p *Parser) unescape(data []byte, escaped int) ([]byte, error) {
	value := make([]byte, len(data)-escaped)
	at := 0
	for {
		index := bytes.IndexByte(data, '\\')
		if index == -1 {
			copy(value[at:], data)
			break
		}
		at += copy(value[at:], data[:index])
		switch data[index+1] {
		case 'n':
			value[at] = '\n'
		case 'r':
			value[at] = '\r'
		case 't':
			value[at] = '\t'
		case '"':
			value[at] = '"'
		case '\\':
			value[at] = '\\'
		default:
			return nil, p.error(fmt.Sprintf("Unknown escape sequence \\%s", string(data[index+1])))
		}
		at++
		data = data[index+2:]
	}
	return value, nil
}

func clone(data []byte) []byte {
	c := make([]byte, len(data))
	copy(c, data)
	return c
}
