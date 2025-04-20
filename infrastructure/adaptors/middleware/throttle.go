package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"kv/utils"
)

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

type bucket struct {
	tokens uint
	time   time.Time
	sync.RWMutex
}

func DefaultIdGetter(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

type RemoteIdGetter func(r *http.Request) string

func Throttle(get_id RemoteIdGetter, max_tokens uint, tokens_per_sec uint) MiddlewareFunc {
	var lock sync.RWMutex
	buckets := map[string]*bucket{}

	go func() {
		for range time.Tick(5 * time.Second) {
			clearOldBuckets(&lock, buckets, time.Duration(10*time.Second))
		}
	}()

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			uid := get_id(r)
			b := getOrCreateBucket(uid, &lock, max_tokens, buckets)
			if !isAllowerRequest(b, tokens_per_sec, max_tokens) {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			next(w, r)
		}
	}
}

func getOrCreateBucket(uid string, lock *sync.RWMutex, max_tokens uint, buckets map[string]*bucket) *bucket {
	lock.RLock()
	b, exists := buckets[uid]
	lock.RUnlock()
	// fmt.Println(b, uid)
	if !exists {
		b = &bucket{tokens: max_tokens, time: time.Now()}
		lock.Lock()
		buckets[uid] = b
		lock.Unlock()
	}
	return b
}

func isAllowerRequest(b *bucket, tokens_per_sec, max_tokens uint) bool {
	delta_time := time.Since(b.time)
	delta_token := uint(delta_time.Microseconds() * int64(tokens_per_sec) / 1_000_000)
	b.Lock()
	defer b.Unlock()
	b.tokens = utils.Min2(max_tokens-1, b.tokens+delta_token-1)
	b.time = time.Now()
	return b.tokens > 0
}

func clearOldBuckets(lock *sync.RWMutex, buckets map[string]*bucket, maxAge time.Duration) {
	lock.Lock()
	defer lock.Unlock()
	for uid, b := range buckets {
		if time.Since(b.time) > maxAge {
			log.Printf("delete: %s with tokens %d", uid, b.tokens)
			delete(buckets, uid)

		}
	}
}
