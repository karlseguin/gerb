package core

import (
	"github.com/karlseguin/bytepool"
	"io"
)

type Context struct {
	Writer   io.Writer
	Data     map[string]interface{}
	Counters map[string]int
	Contents map[string]*bytepool.Item
}
