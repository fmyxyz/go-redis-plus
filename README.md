# go-redis-plus

provide [go-redis/redis](github.com/go-redis/redis) extend features

## Installation

    go get github.com/fmyxyz/go-redis-plus@latest

## Overview

type mapping:

|Action|golang type|redis type|
|-|-|-|
|get|struct|Hash|
|get|map[string]string|Hash|
|get|map[string]interface{}|Hash|
|get|slice|List/Set|
|get|array|List/Set|
|get|other|String|
|set|struct|Hash|
|set|map|Hash|
|set|slice|List/Set|
|set|array|List/Set|
|set|other|String|

## Usage

See: [xredis_test.go](./xredis_test.go)

## Related pojects
- [go-redis/redis](github.com/go-redis/redis)
