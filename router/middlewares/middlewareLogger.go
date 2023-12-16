package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noirbizarre/gonja"
	"github.com/noirbizarre/gonja/exec"
	"go.uber.org/zap"

	"github.com/ditrit/badaas/configuration"
)

func AddLoggerMiddleware(router *mux.Router, middlewareLogger MiddlewareLogger) {
	router.Use(middlewareLogger.Handle)
}

// Log the requests data
type MiddlewareLogger interface {
	// [github.com/gorilla/mux] compatible middleware function
	Handle(next http.Handler) http.Handler
}

// check interface compliance
var _ MiddlewareLogger = (*middlewareLoggerImpl)(nil)

// MiddlewareLogger implementation
type middlewareLoggerImpl struct {
	template *exec.Template
	logger   *zap.Logger
}

// MiddlewareLogger constructor
func NewMiddlewareLogger(
	logger *zap.Logger,
	loggerConfiguration configuration.LoggerConfiguration,
) (MiddlewareLogger, error) {
	requestLogTemplate, err := gonja.FromString(
		loggerConfiguration.GetRequestTemplate(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build jinja template from configuration %w", err)
	}
	return &middlewareLoggerImpl{
		logger:   logger,
		template: requestLogTemplate,
	}, nil
}

// The constructor of the logger middleware
//
// The goal of this middleware is only to print the method used and the API endpoint hit
func (middlewareLogger *middlewareLoggerImpl) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middlewareLogger.logger.Debug(getLogMessage(middlewareLogger.template, r), zap.String("userAgent", r.UserAgent()))
		next.ServeHTTP(w, r)
	})
}

// Create a log message for the logging middleware
func getLogMessage(template *exec.Template, r *http.Request) string {
	result, _ := template.Execute(
		gonja.Context{
			"protocol": r.Proto,
			"method":   r.Method,
			"url":      r.URL.Path,
		})
	return result
}
