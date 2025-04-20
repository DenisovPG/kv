package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Effector func(context.Context) (string, error)

type Throttled func(string) bool

type bucket struct {
	tokens uint
	time   time.Time
	sync.RWMutex
}

func Throttle(max_tokens uint, tokens_per_sec uint) Throttled {
	var lock sync.RWMutex
	buckets := map[string]*bucket{}

	return func(uid string) bool {
		lock.RLock()
		b, exists := buckets[uid]
		lock.Unlock()
		fmt.Println(b, uid)
		if !exists {
			lock.Lock()
			buckets[uid] = &bucket{tokens: tokens_per_sec - 1, time: time.Now()}
			lock.Unlock()
			return true
		}
		delta_time := time.Since(b.time)
		delta_token := uint(delta_time.Microseconds() * int64(tokens_per_sec) / 1_000_000)
		b.Lock()
		defer b.Unlock()
		b.tokens = min(max_tokens-1, b.tokens+delta_token-1)
		b.time = time.Now()
		fmt.Println(b, delta_time, delta_token)
		return b.tokens > 0
	}
}

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}
