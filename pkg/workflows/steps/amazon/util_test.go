package amazon

import (
	"context"
	"testing"
	"time"

	"github.com/supergiant/control/pkg/sgerrors"
)

func TestFindOutboundIPCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := FindOutboundIP(ctx, func() (string, error) {
		return "", sgerrors.ErrNotFound
	})

	if err != context.Canceled {
		t.Errorf("Expected error %v actual %v", context.Canceled, err)
	}
}

func TestFindOutboundIPDeadline(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Nanosecond*1)
	time.Sleep(time.Nanosecond * 2)

	_, err := FindOutboundIP(ctx, func() (string, error) {
		return "", sgerrors.ErrNotFound
	})

	if err != context.DeadlineExceeded {
		t.Errorf("Expected error %v actual %v", context.DeadlineExceeded, err)
	}
}

func TestFindOutboundIPSuccess(t *testing.T) {
	expectedAddr := "10.20.30.40"

	addr, _ := FindOutboundIP(context.Background(), func() (string, error) {
		return expectedAddr, nil
	})

	if addr != expectedAddr {
		t.Errorf("Expected addr %v actual %v", expectedAddr, addr)
	}
}
