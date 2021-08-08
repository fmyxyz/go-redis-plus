package redis

import (
	"encoding/json"
	"reflect"
	"strconv"
)

func toString(value reflect.Value) (bs string) {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.Complex64, reflect.Complex128:
		return strconv.FormatComplex(value.Complex(), 'f', -1, 128)
	case reflect.String:
		return value.String()
	default:
		bytes := toByte(value)
		if bytes != nil {
			return bytesToString(bytes)
		}
		return ""
	}
}

func toByte(value reflect.Value) (bs []byte) {
	v := value.Interface()
	if ok, bytes := dotType2Byte(v); ok {
		return bytes
	}
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Map, reflect.Struct, reflect.Slice, reflect.Array:
		bytes, _ := json.Marshal(v)
		return bytes
	case reflect.Interface:
		return toByte(reflect.ValueOf(v))
	default:
		return stringToBytes(toString(value))
	}
}
