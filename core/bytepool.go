package core

import (
	"github.com/karlseguin/bytepool"
)

var BytePool = bytepool.New(65536, 64)
