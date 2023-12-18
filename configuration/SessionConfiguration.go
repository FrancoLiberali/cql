package configuration

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/ditrit/badaas/utils"
)

// The config keys regarding the session handling settings
const (
	SessionDurationKey     string = "session.duration"
	SessionPullIntervalKey string = "session.pullInterval"
	SessionRollIntervalKey string = "session.rollDuration"
)

// Hold the configuration values to handle the sessions
type SessionConfiguration interface {
	Holder
	GetSessionDuration() time.Duration
	GetPullInterval() time.Duration
	GetRollDuration() time.Duration
}

// Concrete implementation of the SessionConfiguration interface
type sessionConfigurationImpl struct {
	sessionDuration time.Duration
	pullInterval    time.Duration
	rollDuration    time.Duration
}

// Instantiate a new configuration holder for the session management
func NewSessionConfiguration() SessionConfiguration {
	sessionConfiguration := new(sessionConfigurationImpl)
	sessionConfiguration.Reload()

	return sessionConfiguration
}

// Return the session duration
func (sessionConfiguration *sessionConfigurationImpl) GetSessionDuration() time.Duration {
	return sessionConfiguration.sessionDuration
}

// Return the pull interval
func (sessionConfiguration *sessionConfigurationImpl) GetPullInterval() time.Duration {
	return sessionConfiguration.pullInterval
}

// Return the roll interval
func (sessionConfiguration *sessionConfigurationImpl) GetRollDuration() time.Duration {
	return sessionConfiguration.rollDuration
}

// Reload session configuration
func (sessionConfiguration *sessionConfigurationImpl) Reload() {
	sessionConfiguration.sessionDuration = utils.IntToSecond(int(viper.GetUint(SessionDurationKey)))
	sessionConfiguration.pullInterval = utils.IntToSecond(int(viper.GetUint(SessionPullIntervalKey)))
	sessionConfiguration.rollDuration = utils.IntToSecond(int(viper.GetUint(SessionRollIntervalKey)))
}

// Log the values provided by the configuration holder
func (sessionConfiguration *sessionConfigurationImpl) Log(logger *zap.Logger) {
	logger.Info("Session configuration",
		zap.Duration("sessionDuration", sessionConfiguration.sessionDuration),
		zap.Duration("pullInterval", sessionConfiguration.pullInterval),
		zap.Duration("rollDuration", sessionConfiguration.rollDuration),
	)
}
