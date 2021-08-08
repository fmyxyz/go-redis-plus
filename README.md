# go-redis-plus

provide [redis/v8](github.com/go-redis/redis) extend features

## Installation

    go get github.com/fmyxyz/go-redis-plus@latest

## Overview

type mapping:

|Action|golang type|redis type|
|-|-|-|
|get|struct|Hash|
|get|map[string]string|Hash|
|get|slice|List|
|get|array|List|
|get|other|String|
|set|struct|Hash|
|set|map|Hash|
|set|slice|List|
|set|array|List|
|set|other|String|

## Usage

See: [xredis_test.go](./xredis_test.go)

## Related pojects
- [redis/v8](github.com/go-redis/redis)
