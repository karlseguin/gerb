package core

import (
	"bytes"
	"fmt"
)

type Parser struct {
	end      int
	position int
	data     []byte
}

func NewParser(data []byte) *Parser {
	p := &Parser{
		data: data,
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
		p.Dump()
		if p.Prev() == '<' {
			return &Literal{clone(p.data[start : p.position-1])}
		}
	}
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

func (p *Parser) SkipUntil(b byte) bool {
	if at := bytes.IndexByte(p.data[p.position:], b); at != -1 {
		p.position = at
		return true
	}
	p.position = len(p.data)
	return false
}

func (p *Parser) Consume() byte {
	if p.position == p.end {
		return 0
	}
	c := p.data[p.position]
	p.position++
	return c
}

func (p *Parser) Prev() byte {
	return p.data[p.position-1]
}

func (p *Parser) Dump() {
	fmt.Println(string(p.data[p.position:]))
}

func clone(data []byte) []byte {
	c := make([]byte, len(data))
	copy(c, data)
	return c
}
