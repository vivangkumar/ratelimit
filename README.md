## ratelimit [![Go Reference](https://pkg.go.dev/badge/github.com/vivangkumar/ratelimit.svg)](https://pkg.go.dev/github.com/vivangkumar/ratelimit) ![CI](https://github.com/vivangkumar/ratelimit/actions/workflows/ci.yaml/badge.svg?branch=main)

Ratelimit implements a [token bucket](https://en.wikipedia.org/wiki/Token_bucket) based rate limiter.

The library presents a minimal interface to use.

## Install

```
go get github.com/vivangkumar/ratelimit
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vivangkumar/ratelimit"
)

func main() {
	// Creates a rate limiter with a maximum burst of 100 tokens
	// allowing 100 tokens per second to be acquired.
	rl := ratelimit.New(100, 100, 1*time.Second)

	// Get one token.
	if rl.Add() {
		fmt.Println("yay")
	}

	// Get multiple tokens.
	if rl.AddN(20) {
		fmt.Println("woop")
	}

	// Wait for a single token (or context cancellation)
	ctx, cancelFn := context.WithDeadline(context.Background(), 3*time.Second)
	defer cancelFn()

	err := rl.Wait(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
}
```

## Run linter and formatter

```
make fmt
```

## Run tests

```
make test
```

## Build

```
make build
```

[doc-img]: https://pkg.go.dev/badge/vivangkumar/ratelimit
[doc]: https://pkg.go.dev/vivangkumar/ratelimit
[ci-img]: https://github.com/vivangkumar/ratelimit/actions/workflows/ci.yaml/badge.svg?branch=main

## Changelog

Please add changes between release to the [Changelog](CHANGELOG.md).
