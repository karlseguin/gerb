package gerb

import (
	"github.com/karlseguin/gerb/core"
)

type Configuration struct{}

func Configure() *Configuration {
	return new(Configuration)
}

func (c *Configuration) Logger(logger core.Logger) {
	if logger == nil {
		core.Log = new(core.SilentLogger)
	} else {
		core.Log = logger
	}
}

func (c *Configuration) MaxContentSize(size int) {
	core.BytePool.SetCapacity(size)
}

func (c *Configuration) MinContentPoolSize(count int) {
	core.BytePool.SetCount(count)
}
