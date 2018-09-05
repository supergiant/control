package util

import (
	"context"
	"testing"
	"time"
)

func TestCountDown(t *testing.T) {
	testCases := []int{1, 5, 10, 50, 100}

	for _, n := range testCases {
		l := NewCountdownLatch(context.Background(), n)

		for i := 0; i < n; i++ {
			go func() {
				l.CountDown()
			}()
		}

		l.Wait()
	}
}

func TestWaitTimeout(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Nanosecond*0)
	l := NewCountdownLatch(ctx, 1)
	l.Wait()
}
