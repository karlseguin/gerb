package gerb

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
)

// Parse the bytes into a gerb template
func Parse(data []byte, useCache bool) (*Template, error) {
	if useCache == false {
		return newTemplate(data)
	}
	hasher := sha1.New()
	hasher.Write(data)
	key := fmt.Sprintf("%x", hasher.Sum(nil))

	template := cache.get(key)
	if template == nil {
		var err error
		template, err = newTemplate(data)
		if err != nil {
			return nil, err
		}
		cache.set(key, template)
	}
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
