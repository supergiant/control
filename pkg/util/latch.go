package util

import (
	"sync"
)

type CountdownLatch struct {
	zeroChan chan bool

	once  sync.Once
	m     sync.Mutex
	count int32
}

func NewCountdownLatch(count int) *CountdownLatch {
	return &CountdownLatch{
		make(chan bool, 1),
		sync.Once{},
		sync.Mutex{},
		int32(count),
	}
}

// CountDown decrements the counter
func (c *CountdownLatch) CountDown() {
	c.m.Lock()
	defer c.m.Unlock()

	c.count -= 1
	if c.count <= 0 {
		c.once.Do(func() {
			close(c.zeroChan)
		})
		c.count = 0
	}
}

// Wait until count down to zero
func (c *CountdownLatch) Wait() {
	<-c.zeroChan
}
