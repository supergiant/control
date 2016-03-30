package client

import (
	"fmt"
	"time"
)

func waitFor(desc string, d time.Duration, i time.Duration, fn func() (bool, error)) error {
	started := time.Now()
	for {
		elapsed := time.Since(started)
		if elapsed > d {
			return fmt.Errorf("Timed out waiting for %s", desc)
		}
		if ok, err := fn(); ok {
			return nil
		} else if err != nil {
			return err
		}
		time.Sleep(i)
	}
}
