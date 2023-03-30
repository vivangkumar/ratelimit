package ratelimiter_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	ratelimiter "github.com/vivangkumar/ratelimit"
)

func TestRateLimiter_Add_TokensAvailable(t *testing.T) {
	r, err := ratelimiter.New(100, 100, 1*time.Second)
	assert.Nil(t, err)

	assert.True(t, r.Add())
	assert.True(t, r.Add())
	assert.True(t, r.Add())
	assert.True(t, r.Add())
}

func TestRateLimiter_NewError(t *testing.T) {
	_, err := ratelimiter.New(100, 0, 1*time.Second)
	assert.Error(t, err)
}

func TestRateLimiter_AddN(t *testing.T) {
	r, err := ratelimiter.New(10, 10, 1*time.Second)
	assert.Nil(t, err)

	assert.True(t, r.AddN(10))
}

func TestRateLimiter_Add_NoTokens(t *testing.T) {
	r, err := ratelimiter.New(1, 1, 1*time.Second)
	assert.Nil(t, err)

	assert.True(t, r.Add())
	assert.False(t, r.Add())
}

func TestRateLimiter_RefreshToken(t *testing.T) {
	r, err := ratelimiter.New(10, 10, 1*time.Second)
	assert.Nil(t, err)

	assert.True(t, r.AddN(10))

	// This should be false since we have a refresh of 100ms.
	<-time.After(10 * time.Millisecond)
	assert.False(t, r.Add())

	// This should now be true since we're above the 100ms
	// required for the token to be refilled.
	<-time.After(90 * time.Millisecond)
	assert.True(t, r.Add())
}

func TestRateLimiter_Wait(t *testing.T) {
	r, err := ratelimiter.New(100, 10, 1*time.Second)
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Take all tokens.
	assert.True(t, r.AddN(100))

	err = r.Wait(ctx)
	assert.Nil(t, err)
}

func TestRateLimiter_WaitN(t *testing.T) {
	r, err := ratelimiter.New(100, 10, 1*time.Second)
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Take all tokens.
	assert.True(t, r.AddN(100))

	// Wait for 20 tokens.
	err = r.WaitN(ctx, 20)
	assert.Nil(t, err)
}

func TestRateLimiter_WaitN_ExceedsMaxTokens(t *testing.T) {
	r, err := ratelimiter.New(100, 10, 1*time.Second)
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Take all tokens.
	assert.True(t, r.AddN(100))

	// Ask for 200 tokens which exceeds the max.
	err = r.WaitN(ctx, 200)
	assert.Error(t, err)
}

func TestRateLimiter_WaitCancel(t *testing.T) {
	r, err := ratelimiter.New(100, 10, 1*time.Minute)
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Take all tokens.
	assert.True(t, r.AddN(100))

	err = r.Wait(ctx)
	assert.Error(t, err)
}

func ExampleNew() {
	// Create a new rate limiter instance.
	_, err := ratelimiter.New(100, 10, time.Second)
	if err != nil {
		fmt.Printf("context cancelled waiting for token: %s\n", err.Error())
	}
}

func ExampleRateLimiter_Add() {
	r, err := ratelimiter.New(100, 10, time.Second)
	if err != nil {
		fmt.Printf("context cancelled waiting for token: %s\n", err.Error())
	}

	// Use rate limiter.
	if !r.Add() {
		fmt.Println("Oops! We're not allowed!")
	}

	fmt.Println("Yay! We're allowed")
}

func ExampleRateLimiter_AddN() {
	r, err := ratelimiter.New(100, 10, time.Second)
	if err != nil {
		fmt.Printf("error creating rate limiter: %s\n", err.Error())
	}

	// Use rate limiter.
	if !r.AddN(80) {
		fmt.Println("Oops! We're not allowed!")
	}

	fmt.Println("Yay! We're allowed")
}

func ExampleRateLimiter_Wait() {
	r, err := ratelimiter.New(100, 10, time.Second)
	if err != nil {
		fmt.Printf("error creating rate limiter: %s\n", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = r.Wait(ctx)
	if err != nil {
		fmt.Printf("context cancelled waiting for token: %s\n", err.Error())
	}

	fmt.Println("Yay! We're allowed")
}

func ExampleRateLimiter_WaitN() {
	r, err := ratelimiter.New(100, 10, time.Second)
	if err != nil {
		fmt.Printf("error creating rate limiter: %s\n", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = r.WaitN(ctx, 100)
	if err != nil {
		fmt.Printf("context cancelled waiting for token: %s\n", err.Error())
	}

	fmt.Println("Yay! We're allowed")
}

func ExampleRateLimiter() {
	r, err := ratelimiter.New(100, 10, 1*time.Second)
	if err != nil {
		fmt.Printf("error creating rate limiter: %s\n", err.Error())
	}

	if r.AddN(100) {
		fmt.Println("got all tokens")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	for i := 1; i <= 10; i++ {
		err := r.Wait(ctx)
		if err != nil {
			fmt.Println("error waiting for token")
			return
		}

		now := time.Now()
		// ~100ms
		fmt.Println(i, now.Sub(start))
		start = now
	}
}
