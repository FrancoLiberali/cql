package utils

import "time"

// Convert int (seconds) to [time.Duration]
func IntToSecond(numberOfSeconds int) time.Duration {
	return time.Duration(numberOfSeconds) * time.Second
}
