package gerb

import (
	"sync"
)

var cache = &Cache{lookup: make(map[string]*Template)}

type Cache struct {
	sync.RWMutex
	lookup map[string]*Template
}

func (c *Cache) get(key string) *Template {
	c.RLock()
	defer c.RUnlock()
	return c.lookup[key]
}

func (c *Cache) set(key string, template *Template) {
	c.Lock()
	defer c.Unlock()
	c.lookup[key] = template
}
