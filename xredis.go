package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// NewClient returns a client to the Redis Server specified by Options.
func NewClient(opt *Options) *Client {
	if opt.Tag == "" {
		opt.Tag = "json"
	}
	if opt.Expiration == 0 {
		opt.Expiration = -1
	}
	client := redis.NewClient(&opt.Options)
	return &Client{Client: client, options: *opt}
}

type Client struct {
	*redis.Client
	options Options
}

// Options keeps the settings to setup redis connection.
type Options struct {
	redis.Options
	Tag         string
	Expiration  time.Duration
	Start, Stop int64
	SliceType
}

type SliceType uint8

const (
	List SliceType = iota
	Set
)

type Option func(opt *Options)

func Tag(tag string) Option {
	return func(opt *Options) {
		opt.Tag = tag
	}
}

func Expiration(expiration time.Duration) Option {
	return func(opt *Options) {
		opt.Expiration = expiration
	}
}

func Range(start, stop int64) Option {
	return func(opt *Options) {
		opt.Start = start
		opt.Stop = stop
	}
}

func RedisTypeList() Option {
	return func(opt *Options) {
		opt.SliceType = List
	}
}

func RedisTypeSet() Option {
	return func(opt *Options) {
		opt.SliceType = Set
	}
}
