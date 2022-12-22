package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
)

// initialize session related config keys
func initSessionCommands(cfg *verdeter.VerdeterCommand) {
	cfg.LKey(configuration.SessionDurationKey, verdeter.IsUint, "", "The duration of a user session in seconds.")
	cfg.SetDefault(configuration.SessionDurationKey, uint(3600*4)) // 4 hours by default

	cfg.LKey(configuration.SessionPullIntervalKey,
		verdeter.IsUint, "", "The refresh interval in seconds. Badaas refresh it's internal session cache periodically.")
	cfg.SetDefault(configuration.SessionPullIntervalKey, uint(30)) // 30 seconds by default

	cfg.LKey(configuration.SessionRollIntervalKey, verdeter.IsUint, "", "The interval in which the user can renew it's session by making a request.")
	cfg.SetDefault(configuration.SessionRollIntervalKey, uint(3600)) // 1 hour by default
}
