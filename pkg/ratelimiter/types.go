package ratelimiter

import "net"

type RateLimiter interface {
	Allow(conn net.Conn) bool
}
