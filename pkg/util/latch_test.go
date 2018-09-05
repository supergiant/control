package util

import (
	"testing"
)

func TestCountDown(t *testing.T) {
	testCases := []int{1, 5, 10, 50, 100}

	for _, n := range testCases {
		l := NewCountdownLatch(n)

		for i := 0; i < n; i++ {
			go func() {
				l.CountDown()
			}()
		}

		l.Wait()
	}
}
