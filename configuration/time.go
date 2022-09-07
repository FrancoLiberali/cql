package configuration

import "time"

func intToSecond(numberOfSeconds int) time.Duration {
	return time.Duration(numberOfSeconds) * time.Second
}
