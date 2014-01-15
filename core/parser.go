package core

import (
	"bytes"
	"fmt"
	"errors"
	"math"
)

type Parser struct {
	end      int
	len int
	position int
	data     []byte
}

func NewParser(data []byte) *Parser {
	p := &Parser{
		data: data,
		len: len(data),
		end:  len(data) - 1,
	}
	return p
}

func (p *Parser) ReadLiteral() *Literal {
	start := p.position
	for {
		if p.SkipUntil('%') == false {
			return nil
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
	if first == 0 {
		return nil, p.error("Unrecognized value in output tag")
	}
	if first >= '0' && first <= '9' {
		return p.ReadNumber(negate)
	}
	// if first == '\'' {
	// 	return p.ReadChar()
	// }
	// if first == '"' {
	// 	return p.ReadString()
	// }
	return nil, nil
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
			if isDecimal { break }
			target = &fraction
			partLength = 0
			isDecimal = true
			continue
		}
		if c < '0' || c > '9' {
			break
		}
		partLength++
		*target = *target * 10 + int(c - '0')
	}

	if isDecimal {
		value := float64(integer) + float64(fraction) / math.Pow10(partLength)
		if negate { value *= -1 }
		return &FloatValue{value}, nil
	}
	if negate { integer *= -1 }
	return &IntValue{integer}, nil
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
		p.position = at
		return true
	}
	p.position = len(p.data)
	return false
}

func (p *Parser) SkipSpaces() byte {
	for ; p.position < p.end; p.position++ {
		c := p.data[p.position]
		if c != ' ' && c !=  '\t' && c != '\n' && c != '\r' {
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
	return errors.New(fmt.Sprintf("%s: %q", s, string(p.data[start:end])))
}

func clone(data []byte) []byte {
	c := make([]byte, len(data))
	copy(c, data)
	return c
}
