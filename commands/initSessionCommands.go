package commands

import (
	"github.com/ditrit/verdeter"
)

// initialize session related config keys
func initSessionCommands(cfg *verdeter.VerdeterCommand) {

	sessionDurationKey := "session.duration"
	cfg.LKey(sessionDurationKey, verdeter.IsUint, "", "The duration of a user session in seconds.")
	cfg.SetDefault(sessionDurationKey, uint(3600*4)) // 4 hours by default

	sessionPullIntervalKey := "session.pullInterval"
	cfg.LKey(sessionPullIntervalKey, verdeter.IsUint, "", "The refresh interval in seconds. Badaas refresh it's internal session cache periodically.")
	cfg.SetDefault(sessionPullIntervalKey, uint(30)) // 30 seconds by default

	sessionRollIntervalKey := "session.rollDuration"
	cfg.LKey(sessionRollIntervalKey, verdeter.IsUint, "", "The interval in which the user can renew it's session by making a request.")
	cfg.SetDefault(sessionRollIntervalKey, uint(3600)) // 1 hour by default
}
