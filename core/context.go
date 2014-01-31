package core

import (
	"io"
	"github.com/karlseguin/bytepool"
)

type Context struct {
	Writer   io.Writer
	Data     map[string]interface{}
	Counters map[string]int
	Contents map[string]*bytepool.Item
}
