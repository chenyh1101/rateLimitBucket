package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	Rate         int64 //the rate of generating token
	Cap          int64 //the capacity of bucket
	Tokens       int64 //current tokens
	LastTokenSec int64 //the time of last generated tokens
	lock         *sync.Mutex
}

func NewBucket(rate, cap int64) *TokenBucket {
	if cap < 1 {
		panic("error cap")
	} else {
		return &TokenBucket{
			Cap:    cap,
			Rate:   rate,
			Tokens: cap,
			lock:   new(sync.Mutex),
		}
	}
}

func (tb *TokenBucket) Consume() bool {
	tb.lock.Lock()
	defer tb.lock.Unlock()
	now := time.Now().Unix()
	tb.Tokens += (now - tb.LastTokenSec) * tb.Rate
	if tb.Cap < tb.Tokens {
		tb.Tokens = tb.Cap
	}
	tb.LastTokenSec = now
	if tb.Tokens > 0 {
		tb.Tokens--
		return true
	} else {
		return false
	}

}
func main() {
	tb := NewBucket(10, 100)
	for i := 0; i < 100000; i++ {
		if tb.Consume() {
			fmt.Printf("tokens:=%d\n", tb.Tokens)
			fmt.Printf("%+v s,i=%d\n", i/100, i)
		} else {
			fmt.Printf("%+v s\n", i/100)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
