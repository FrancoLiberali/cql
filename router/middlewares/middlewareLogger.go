package middlewares

import (
	"net/http"

	"github.com/ditrit/badaas/configuration"
	"github.com/noirbizarre/gonja"
	"github.com/noirbizarre/gonja/exec"
	"go.uber.org/zap"
)

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

// The constructor of the logger middleware
//
// The goal of this middleware is only to print the method used and the API endpoint hit
func CreateLoggerMiddleware() func(next http.Handler) http.Handler {
	loggerConfiguration := configuration.NewLoggerConfiguration()
	requestLogTemplate, err := gonja.FromString(
		loggerConfiguration.GetRequestTemplate(),
	)
	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			zap.L().Debug(getLogMessage(requestLogTemplate, r))
			next.ServeHTTP(w, r)
		})
	}

}
