package utils

import "time"

// RunWithRetry run a closure and retry for times when it fails
func RunWithRetry(f func() error, interval time.Duration, maxRetries int) (err error) {
	var times int
	for times < maxRetries {
		if err = f(); err != nil {
			times++

			time.Sleep(interval)
		} else {
			break
		}
	}
	return
}
