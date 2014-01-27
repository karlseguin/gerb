package core

import (
	"io"
)

type Context struct {
	Writer   io.Writer
	Data     interface{}
	Counters map[string]int
}
