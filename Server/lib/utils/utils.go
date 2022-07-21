package utils

import (
	"fmt"
	"time"
)

var (
	Unlimited = -1
)

// Retry execute f() n times before returning error
func Retry(n int, sleep time.Duration, f func() error) (err error) {
	for i := 0; i < n || n == Unlimited; i++ {
		if i > 0 {
			time.Sleep(sleep)
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("after %d n, last error: %s", n, err)
}
