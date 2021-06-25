package retry

import "time"

func BackoffLinear(wait time.Duration) BackoffFunc {
	return func(attempt uint) time.Duration {
		return wait
	}
}
