package r

import (
	"fmt"
	"strconv"
)

// Convert arbitrary data to []byte
func ToBytes(data interface{}) []byte {
	switch typed := data.(type) {
	case byte:
		return []byte{typed}
	case []byte:
		return typed
	case string:
		return []byte(typed)
	case bool:
		return []byte(strconv.FormatBool(typed))
	case float64:
		return []byte(strconv.FormatFloat(typed, 'g', -1, 64))
	case uint64:
		return []byte(strconv.FormatUint(typed, 10))
	case uint:
		return []byte(strconv.FormatUint(uint64(typed), 10))
	case int:
		return []byte(strconv.Itoa(typed))
	case fmt.Stringer:
		return []byte(typed.String())
	}
	return []byte(fmt.Sprintf("%v", data))
}
