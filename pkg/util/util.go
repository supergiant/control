package util

import (
	"fmt"
	"math/rand"
	"time"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// func UniqueStrings(in []string) (out []string) {
// 	tab := make(map[string]struct{})
// 	for _, str := range in {
// 		if _, ok := tab[str]; !ok {
// 			tab[str] = struct{}{}
// 			out = append(out, str)
// 		}
// 	}
// 	return out
// }

func WaitFor(desc string, d time.Duration, i time.Duration, fn func() (bool, error)) error {
	started := time.Now()
	for {
		if done, err := fn(); done {
			return nil
		} else if err != nil {
			return err
		}
		elapsed := time.Since(started)
		if elapsed > d {
			return fmt.Errorf("Timed out waiting for %s", desc)
		}
		time.Sleep(i)
	}
}
