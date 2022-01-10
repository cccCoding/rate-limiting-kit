package fixedwindowcounter

import (
	"sync"
	"sync/atomic"
	"time"

	rate_limiting_kit "github.com/cccCoding/rate-limiting-kit"
)

var (
	once sync.Once
)

var _ rate_limiting_kit.RateLimiter = &fixedWindowCounter{}

type fixedWindowCounter struct {
	snippet         time.Duration
	currentRequests int32
	allowRequests   int32
}

func New(snippet time.Duration, allowRequests int32) *fixedWindowCounter {
	return &fixedWindowCounter{snippet: snippet, allowRequests: allowRequests}
}

func (f *fixedWindowCounter) Take() error {
	once.Do(func() {
		go func() {
			for {
				select {
				case <-time.After(f.snippet):
					atomic.StoreInt32(&f.currentRequests, 0)
				}
			}
		}()
	})

	curRequests := atomic.LoadInt32(&f.currentRequests)
	if curRequests >= f.allowRequests {
		return rate_limiting_kit.ErrExceedLimit
	}
	if !atomic.CompareAndSwapInt32(&f.currentRequests, curRequests, curRequests+1) {
		return rate_limiting_kit.ErrExceedLimit
	}
	return nil
}
