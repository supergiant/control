package util

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestWaitForCancelled(t *testing.T) {
	t.Parallel()
	fn := func() (bool, error) {
		return false, nil
	}

	d := time.Second * 1
	p := time.Millisecond * 10
	ctx, cancel := context.WithTimeout(context.Background(), d)

	// Make a few checks and cancel wait for
	go func() {
		time.Sleep(p * 3)
		cancel()
	}()

	err := WaitFor(ctx, "Test cancelled", p, fn)

	if errors.Cause(context.Canceled) != context.Canceled {
		t.Errorf("Unexpected error expected %v actual %v", context.Canceled, err)
	}
}

func TestWaitForDeadline(t *testing.T) {
	t.Parallel()
	fn := func() (bool, error) {
		return false, nil
	}

	d := time.Millisecond * 100
	p := time.Millisecond * 10

	ctx, _ := context.WithTimeout(context.Background(), d)
	err := WaitFor(ctx, "Test deadline", p, fn)

	if errors.Cause(err) != context.DeadlineExceeded {
		t.Errorf("Unexpected error expected %v actual %v", context.Canceled, err)
	}
}

func TestWaitForSucceed(t *testing.T) {
	t.Parallel()

	result := make(chan bool)

	go func() {
		result <- false
		result <- false
		result <- true
	}()

	d := time.Millisecond * 100
	p := time.Millisecond * 10

	fn := func() (bool, error) {
		return <-result, nil
	}

	ctx, _ := context.WithTimeout(context.Background(), d)
	err := WaitFor(ctx, "Test succeed", p, fn)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
