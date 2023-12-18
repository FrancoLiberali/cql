package configuration

const (
	defaultDatabaseRetry                  = uint(10)
	defaultDatabaseRetryDuration          = uint(5)
	defaultServerTimeout                  = 15
	defaultServerAddress                  = "0.0.0.0"
	defaultServerPort                     = 8000
	defaultServerPaginationMaxElemPerPage = uint(100)
	defaultSessionDuration                = uint(3600 * 4) // 4 hours
	defaultSessionPullInterval            = uint(30)       // 30 seconds
	defaultSessionRollInterval            = uint(3600)     // 1 hour
	defaultLoggerSlowQueryThreshold       = 200            // milliseconds
	defaultLoggerSlowTransactionThreshold = 200            // milliseconds
)
