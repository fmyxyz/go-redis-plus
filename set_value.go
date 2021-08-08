package redis

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"
)

func setValueByString(value reflect.Value, val string) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(val, 0, value)
	case reflect.Int8:
		return setIntField(val, 8, value)
	case reflect.Int16:
		return setIntField(val, 16, value)
	case reflect.Int32:
		return setIntField(val, 32, value)
	case reflect.Int64:
		if _, ok := value.Interface().(time.Duration); ok {
			return setTimeDuration(val, value)
		}
		return setIntField(val, 64, value)
	case reflect.Uint:
		return setUintField(val, 0, value)
	case reflect.Uint8:
		return setUintField(val, 8, value)
	case reflect.Uint16:
		return setUintField(val, 16, value)
	case reflect.Uint32:
		return setUintField(val, 32, value)
	case reflect.Uint64:
		return setUintField(val, 64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, 32, value)
	case reflect.Float64:
		return setFloatField(val, 64, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct:
		if _, ok := value.Interface().(time.Time); ok {
			return setTimeField(val, value)
		}
		return json.Unmarshal(stringToBytes(val), value.Addr().Interface())
	case reflect.Map, reflect.Array:
		return json.Unmarshal(stringToBytes(val), value.Addr().Interface())
	case reflect.Slice:
		if _, ok := value.Interface().([]byte); ok {
			return setByteSlice(val, value)
		}
		return json.Unmarshal(stringToBytes(val), value.Addr().Interface())
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, value reflect.Value) error {
	t, err := time.Parse(time.RFC3339Nano, val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(t))
	return nil
}

func setArray(vals []string, value reflect.Value) error {
	for i, s := range vals {
		err := setValueByString(value.Index(i), s)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(vals []string, value reflect.Value) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := setArray(vals, slice)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func setByteSlice(val string, value reflect.Value) error {
	value.Set(reflect.ValueOf(stringToBytes(val)))
	return nil
}

func setTimeDuration(val string, value reflect.Value) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}
