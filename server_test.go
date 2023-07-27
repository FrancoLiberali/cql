package badaas

// This files holds the tests for the server.go file.

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/configuration"
)

func TestCreateServer(t *testing.T) {
	handl := http.NewServeMux()

	srv := createServer(handl, configuration.NewHTTPServerConfiguration())
	assert.NotNil(t, srv)
}
