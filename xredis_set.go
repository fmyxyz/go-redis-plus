package redis

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"
)

func (c *Client) SetValue(ctx context.Context, key string, value interface{}, opts ...Option) (err error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	valValue := reflect.ValueOf(value)
	if valValue.Kind() == reflect.Ptr {
		valValue = valValue.Elem()
	}
	switch valValue.Kind() {
	case reflect.Map:
		return c.setMapValue(ctx, key, valValue, options)
	case reflect.Struct:
		if val, ok := valValue.Interface().(time.Time); ok {
			bytes := stringToBytes(val.Format(time.RFC3339Nano))
			return c.Set(ctx, key, bytes, options.Expiration).Err()
		}
		return c.setStructValue(ctx, key, valValue, options)
	case reflect.Array, reflect.Slice:
		return c.setListValue(ctx, key, valValue, options)
	default:
		if ok, bytes := dotType2Byte(value); ok {
			return c.Set(ctx, key, bytes, options.Expiration).Err()
		}
		return c.setSingleValue(ctx, key, valValue, options)
	}
}

func (c *Client) SetSingleValue(ctx context.Context, key string, value interface{}, opts ...Option) (err error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	if ok, bytes := dotType2Byte(value); ok {
		return c.Set(ctx, key, bytes, options.Expiration).Err()
	}

	valValue := reflect.ValueOf(value)
	if valValue.Kind() == reflect.Ptr {
		valValue = valValue.Elem()
	}

	return c.setSingleValue(ctx, key, valValue, options)
}

func (c *Client) SetSliceValue(ctx context.Context, key string, value interface{}, opts ...Option) (err error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	valValue := reflect.ValueOf(value)
	if valValue.Kind() == reflect.Ptr {
		valValue = valValue.Elem()
	}

	switch valValue.Kind() {
	case reflect.Array, reflect.Slice:
		switch options.SliceType {
		case List:
			return c.setListValue(ctx, key, valValue, options)
		default: // Set
			return c.setSetValue(ctx, key, valValue, options)
		}
	default:
		return errors.New("value is not array or slice")
	}
}

func (c *Client) SetStructValue(ctx context.Context, key string, value interface{}, opts ...Option) (err error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	valValue := reflect.ValueOf(value)
	if valValue.Kind() == reflect.Ptr {
		valValue = valValue.Elem()
	}

	switch valValue.Kind() {
	case reflect.Struct:
		return c.setStructValue(ctx, key, valValue, options)
	default:
		return errors.New("value is not struct")
	}
}

func (c *Client) SetMapValue(ctx context.Context, key string, value interface{}, opts ...Option) (err error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	valValue := reflect.ValueOf(value)
	if valValue.Kind() == reflect.Ptr {
		valValue = valValue.Elem()
	}

	switch valValue.Kind() {
	case reflect.Map:
		return c.setMapValue(ctx, key, valValue, options)
	default:
		return errors.New("value is not map")
	}
}

func (c *Client) setSingleValue(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	switch valValue.Kind() {
	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		bytes, err := json.Marshal(valValue.Interface())
		if err != nil {
			return err
		}
		return c.Set(ctx, key, bytes, options.Expiration).Err()
	default:
		bytes := toByte(valValue)
		return c.Set(ctx, key, bytes, options.Expiration).Err()
	}
}

func (c *Client) setListValue(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	valLen := valValue.Len()
	if valLen == 0 {
		return nil
	}
	vals := make([]interface{}, valLen)
	for i := 0; i < valLen; i++ {
		sliceVal := valValue.Index(i)
		vals[i] = toByte(sliceVal)
	}
	err = c.RPush(ctx, key, vals).Err()
	if err != nil {
		return err
	}
	err = c.expireKeyTTl(ctx, key, options.Expiration)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) setSetValue(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	valLen := valValue.Len()
	if valLen == 0 {
		return nil
	}
	vals := make([]interface{}, valLen)
	for i := 0; i < valLen; i++ {
		sliceVal := valValue.Index(i)
		vals[i] = toByte(sliceVal)
	}
	err = c.SAdd(ctx, key, vals).Err()
	if err != nil {
		return err
	}
	err = c.expireKeyTTl(ctx, key, options.Expiration)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) setStructValue(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	numField := valValue.NumField()
	m := make(map[string]interface{}, numField)
	for i := 0; i < numField; i++ {
		valType := valValue.Type()
		key := getStructKey(valType, i, options.Tag)
		m[key] = toByte(valValue.Field(i))
	}
	if len(m) == 0 {
		return nil
	}
	err = c.HSet(ctx, key, m).Err()
	if err != nil {
		return err
	}
	err = c.expireKeyTTl(ctx, key, options.Expiration)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) setMapValue(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	m := map[string]interface{}{}
	iter := valValue.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		m[toString(k)] = toByte(v)
	}
	if len(m) == 0 {
		return nil
	}
	err = c.HSet(ctx, key, m).Err()
	if err != nil {
		return err
	}
	err = c.expireKeyTTl(ctx, key, options.Expiration)
	if err != nil {
		return err
	}
	return nil
}

// ????????????????????????
func (c *Client) expireKeyTTl(ctx context.Context, key string, expiration time.Duration) error {
	if expiration > 0 {
		return c.Expire(ctx, key, expiration).Err()
	}
	return nil
}
