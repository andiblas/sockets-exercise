package ratelimiter

import (
	"fmt"
	"net"
	"time"
)

type LeakyBucket struct {
	leakyTicker    *time.Ticker
	availableSpace chan struct{}
}

func NewLeakyBucket(capacity int, fillRate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		leakyTicker:    time.NewTicker(fillRate),
		availableSpace: make(chan struct{}, capacity),
	}
}

func (lb *LeakyBucket) Start() {
	go func() {
		for range lb.leakyTicker.C {
			fmt.Println("adding new available space")
			lb.availableSpace <- struct{}{}
		}
	}()
}

func (lb *LeakyBucket) Allow(_ net.Conn) bool {
	select {
	case <-lb.availableSpace:
		fmt.Println("consuming available space")
		return true
	default:
		return false
	}
}
