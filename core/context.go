package core

import (
	"io"
)

type Context struct {
	Writer   io.Writer
	Data     map[string]interface{}
	Counters map[string]int
}
