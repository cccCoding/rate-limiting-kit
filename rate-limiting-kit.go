package rate_limiting_kit

import "errors"

var (
	ErrExceedLimit = errors.New("Too many requests, exceeded the limit.")
)

type RateLimiter interface {
	Take() error
}
