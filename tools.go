package redis

import (
	"encoding"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func getStructKey(valType reflect.Type, i int, tag string) (key string) {
	field := valType.Field(i)
	if tag == "" {
		return field.Name
	}
	tagName := field.Tag.Get(tag)
	if tagName == "" || tagName == "-" {
		return ""
	}
	tagName = strings.Split(tagName, ",")[0]
	return tagName
}

func dotType2Byte(val interface{}) (ok bool, bs []byte) {
	switch val := val.(type) {
	case nil:
		return ok, stringToBytes("")
	case time.Time:
		return true, stringToBytes(val.Format(time.RFC3339Nano))
	case time.Duration:
		return true, stringToBytes(val.String())
	case encoding.BinaryMarshaler:
		b, _ := val.MarshalBinary()
		return true, b
	}
	return false, nil
}

// stringToBytes converts string to byte slice without a memory allocation.
func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// bytesToString converts byte slice to string without a memory allocation.
func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
