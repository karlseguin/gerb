package gerb

import (
	"crypto/sha1"
	"fmt"
	"github.com/karlseguin/ccache"
	"io/ioutil"
	"time"
)

var cache = ccache.New(ccache.Configure().Size(1024 * 1024 * 10))

// Parse the bytes into a gerb template
func Parse(data []byte, useCache bool) (*Template, error) {
	if useCache == false {
		return newTemplate(data)
	}
	hasher := sha1.New()
	hasher.Write(data)
	key := fmt.Sprintf("%x", hasher.Sum(nil))

	t := cache.Get(key)
	if t != nil {
		return t.(*Template), nil
	}

	template, err := newTemplate(data)
	if err != nil {
		return nil, err
	}
	cache.Set(key, template, time.Hour)
	return template, nil
}

// Parse the string into a erb template
func ParseString(data string, cache bool) (*Template, error) {
	return Parse([]byte(data), cache)
}

// Turn the contents of the specified file into a gerb template
func ParseFile(path string, cache bool) (*Template, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data, cache)
}
