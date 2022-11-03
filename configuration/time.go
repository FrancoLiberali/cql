package configuration

import "time"

// Convert int (seconds) to [time.Duration]
func intToSecond(numberOfSeconds int) time.Duration {
	return time.Duration(numberOfSeconds) * time.Second
}
