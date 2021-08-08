package redis

import (
	"context"
	"errors"
	"reflect"
	"time"
)

func (c *Client) GetValue(ctx context.Context, key string, value interface{}, opts ...Option) (err error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	valValue := reflect.ValueOf(value)
	if valValue.Kind() != reflect.Ptr {
		return errors.New("value CanSet returns false")
	}
	valValue = valValue.Elem()
	switch valValue.Kind() {
	case reflect.Struct:
		if _, ok := valValue.Interface().(time.Time); ok {
			return c.getSingleValue(ctx, key, valValue, options)
		}
		return c.getStructValue(ctx, key, valValue, options)
	case reflect.Array:
		return c.getArrayValue(ctx, key, valValue, options)
	case reflect.Slice:
		return c.getSliceValue(ctx, key, valValue, options)
	case reflect.Map:
		switch val := value.(type) {
		case *map[string]string, *map[string]interface{}:
			return c.getMapValue(ctx, key, val, options)
		}
		return errors.New("value map is not map[string]string, map[string]interface{}")
	default:
		return c.getSingleValue(ctx, key, valValue, options)
	}
}

func (c *Client) GetSingleValue(ctx context.Context, key string, value interface{}, opts ...Option) (err error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	valValue := reflect.ValueOf(value)
	if valValue.Kind() == reflect.Ptr {
		valValue = valValue.Elem()
	}

	return c.getSingleValue(ctx, key, valValue, options)
}

func (c *Client) GetSliceValue(ctx context.Context, key string, value interface{}, opts ...Option) (err error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	valValue := reflect.ValueOf(value)
	if valValue.Kind() == reflect.Ptr {
		valValue = valValue.Elem()
	}

	switch valValue.Kind() {
	case reflect.Array:
		return c.getArrayValue(ctx, key, valValue, options)
	case reflect.Slice:
		return c.getSliceValue(ctx, key, valValue, options)
	default:
		return errors.New("value is not array or slice")
	}
}

func (c *Client) GetStructValue(ctx context.Context, key string, value interface{}, opts ...Option) (err error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}

	valValue := reflect.ValueOf(value)
	if valValue.Kind() != reflect.Ptr {
		return errors.New("value CanSet returns false")
	}
	valValue = valValue.Elem()

	switch valValue.Kind() {
	case reflect.Struct:
		return c.getStructValue(ctx, key, valValue, options)
	default:
		return errors.New("value is not struct")
	}
}

func (c *Client) GetMapValue(ctx context.Context, key string, value interface{}, opts ...Option) (err error) {
	options := c.options
	for _, opt := range opts {
		opt(&options)
	}
	switch val := value.(type) {
	case map[string]string, map[string]interface{}, *map[string]string, *map[string]interface{}:
		return c.getMapValue(ctx, key, val, options)
	default:
		return errors.New("value map is not map[string]string, map[string]interface{}")
	}
}

func (c *Client) getSingleValue(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	str, err := c.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return setValueByString(valValue, str)
}

func (c *Client) getSliceValue(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	if options.Stop == 0 {
		options.Stop = int64(valValue.Len() - 1)
	}
	strings, err := c.LRange(ctx, key, options.Start, options.Stop).Result()
	if err != nil {
		return err
	}
	return setSlice(strings, valValue)
}

func (c *Client) getArrayValue(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	if options.Stop == 0 {
		options.Stop = int64(valValue.Len() - 1)
	}
	strings, err := c.LRange(ctx, key, options.Start, options.Stop).Result()
	if err != nil {
		return err
	}
	return setArray(strings, valValue)
}

func (c *Client) getStructValue(ctx context.Context, key string, valValue reflect.Value, options Options) (err error) {
	fieldKeyIdxMap := make(map[string]int)
	numField := valValue.NumField()
	for i := 0; i < numField; i++ {
		key := getStructKey(valValue.Type(), i, options.Tag)
		if key == "" {
			continue
		}
		fieldKeyIdxMap[key] = i
	}
	fieldKeys := make([]string, 0, len(fieldKeyIdxMap))
	for k := range fieldKeyIdxMap {
		fieldKeys = append(fieldKeys, k)
	}
	fieldVals, err := c.HMGet(ctx, key, fieldKeys...).Result()
	if err != nil {
		return err
	}
	if len(fieldKeys) != len(fieldVals) {
		return errors.New("HMGet should have the same number of keys and vals")
	}
	for i := range fieldVals {
		key := fieldKeys[i]
		idx := fieldKeyIdxMap[key]
		field := valValue.Field(idx)
		valStr, ok := fieldVals[i].(string)
		if !ok {
			continue
		}
		err := setValueByString(field, valStr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) getMapValueSring(ctx context.Context, key string, val map[string]string, options Options) (err error) {
	stringStringMap, err := c.HGetAll(ctx, key).Result()
	if err != nil {
		return err
	}
	for k, v := range stringStringMap {
		val[k] = v
	}
	return nil
}

func (c *Client) getMapValue(ctx context.Context, key string, val interface{}, options Options) (err error) {
	switch val := val.(type) {
	case map[string]string:
		stringStringMap, err := c.HGetAll(ctx, key).Result()
		if err != nil {
			return err
		}
		for k, v := range stringStringMap {
			val[k] = v
		}
	case map[string]interface{}:
		stringStringMap, err := c.HGetAll(ctx, key).Result()
		if err != nil {
			return err
		}
		for k, v := range stringStringMap {
			val[k] = v
		}
	case *map[string]string:
		stringStringMap, err := c.HGetAll(ctx, key).Result()
		if err != nil {
			return err
		}
		for k, v := range stringStringMap {
			(*(val))[k] = v
		}
	case *map[string]interface{}:
		stringStringMap, err := c.HGetAll(ctx, key).Result()
		if err != nil {
			return err
		}
		for k, v := range stringStringMap {
			(*(val))[k] = v
		}
	default:
		return errors.New("value map is not map[string]string, map[string]interface{}")
	}
	return nil
}
