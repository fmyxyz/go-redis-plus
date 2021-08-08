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
		return c.setMap2Redis(ctx, key, valValue, options)
	case reflect.Struct:
		if val, ok := valValue.Interface().(time.Time); ok {
			bytes := stringToBytes(val.Format(time.RFC3339Nano))
			return c.Set(ctx, key, bytes, options.Expiration).Err()
		}
		return c.setStruct2Redis(ctx, key, valValue, options)
	case reflect.Array, reflect.Slice:
		return c.setList2Redis(ctx, key, valValue, options)
	default:
		if ok, bytes := dotType2Byte(value); ok {
			return c.Set(ctx, key, bytes, options.Expiration).Err()
		}
		return c.setSingle2Redis(ctx, key, valValue, options)
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

	return c.setSingle2Redis(ctx, key, valValue, options)
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
		return c.setList2Redis(ctx, key, valValue, options)
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
		return c.setStruct2Redis(ctx, key, valValue, options)
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
		return c.setMap2Redis(ctx, key, valValue, options)
	default:
		return errors.New("value is not map")
	}
}

func (c *Client) setSingle2Redis(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	switch valValue.Kind() {
	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		bytes, err := json.Marshal(valValue.Interface())
		if err != nil {
			return err
		}
		return c.Set(ctx, key, bytes, options.Expiration).Err()
	default:
		return c.Set(ctx, key, valValue.Interface(), options.Expiration).Err()
	}
}

func (c *Client) setList2Redis(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	valLen := valValue.Len()
	vals := make([]interface{}, valLen)
	for i := 0; i < valLen; i++ {
		sliceVal := valValue.Index(i)
		vals[i] = toByte(sliceVal)
	}
	err = c.LPush(ctx, key, vals).Err()
	if err != nil {
		return err
	}
	if options.Expiration != -1 {
		return c.Expire(ctx, key, options.Expiration).Err()
	}
	return nil
}

func (c *Client) setStruct2Redis(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	numField := valValue.NumField()
	m := make(map[string]interface{}, numField)
	for i := 0; i < numField; i++ {
		key := getStructKey(valValue.Type(), i, options.Tag)
		m[key] = toByte(valValue.Field(i))
	}
	err = c.HSet(ctx, key, m).Err()
	if err != nil {
		return err
	}
	if options.Expiration != -1 {
		return c.Expire(ctx, key, options.Expiration).Err()
	}
	return nil
}

func (c *Client) setMap2Redis(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	m := map[string]interface{}{}
	iter := valValue.MapRange()
	if iter.Next() {
		k := iter.Key()
		v := iter.Value()
		m[toString(k)] = toByte(v)
	}
	err = c.HSet(ctx, key, m).Err()
	if err != nil {
		return err
	}
	if options.Expiration != -1 {
		return c.Expire(ctx, key, options.Expiration).Err()
	}
	return nil
}
