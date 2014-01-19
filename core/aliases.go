package core

import (
	"fmt"
	"reflect"
	"strings"
)

var Aliases = make(map[string]map[string]reflect.Value)

func RegisterAlias(packageName, functionName string, f interface{}) {
	packageName = strings.ToLower(packageName)
	functionName = strings.ToLower(functionName)

	pkg, ok := Aliases[packageName]
	if ok == false {
		pkg = make(map[string]reflect.Value)
		Aliases[packageName] = pkg
	}
	pkg[functionName] = reflect.ValueOf(f)
}

func RegisterAliases(packageName string, data ...interface{}) {
	packageName = strings.ToLower(packageName)
	pkg, ok := Aliases[packageName]
	if ok == false {
		pkg = make(map[string]reflect.Value)
		Aliases[packageName] = pkg
	}
	for i := 0; i < len(data); i += 2 {
		functionName := strings.ToLower(data[i].(string))
		pkg[functionName] = reflect.ValueOf(data[i+1])
	}
}

func init() {
	RegisterAliases("strings",
		"ToUpper", strings.ToUpper,
		"ToLower", strings.ToLower,
		"Contains", strings.Contains,
		"ContainsAny", strings.ContainsAny,
		"ContainsRune", strings.ContainsRune,
		"Count", strings.Count,
		"Fields", strings.Fields,
		"HasPrefix", strings.HasPrefix,
		"HasSuffix", strings.HasSuffix,
		"Index", strings.Index,
		"IndexAny", strings.IndexAny,
		"IndexByte", strings.IndexByte,
		"IndexRune", strings.IndexRune,
		"Join", strings.Join,
		"LastIndex", strings.LastIndex,
		"LastIndexAny", strings.LastIndexAny,
		"Repeat", strings.Repeat,
		"Replace", strings.Replace,
		"Split", strings.Split,
		"SplitAfter", strings.SplitAfter,
		"SplitAfterN", strings.SplitAfterN,
		"SplitN", strings.SplitN,
		"Title", strings.Title,
		"ToLower", strings.ToLower,
		"ToLowerSpecial", strings.ToLowerSpecial,
		"ToTitle", strings.ToTitle,
		"ToTitleSpecial", strings.ToTitleSpecial,
		"ToUpper", strings.ToUpper,
		"ToUpperSpecial", strings.ToUpperSpecial,
		"Trim", strings.Trim,
		"TrimLeft", strings.TrimLeft,
		"TrimPrefix", strings.TrimPrefix,
		"TrimRight", strings.TrimRight,
		"TrimSpace", strings.TrimSpace,
		"TrimSuffix", strings.TrimSuffix,
	)

	RegisterAlias("fmt",
		"Sprintf", fmt.Sprintf)
}
