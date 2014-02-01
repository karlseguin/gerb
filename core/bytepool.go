package core

import (
	"github.com/karlseguin/bytepool"
)

var BytePool = bytepool.New(64, 65536)
