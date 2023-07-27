package main

// This files holds the tests for the server.go file.

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/configuration"
)

func Test_addrFromConf(t *testing.T) {
	expected := "192.168.236.222:25100"
	addr := addrFromConf("192.168.236.222", 25100)
	assert.Equal(t, expected, addr)
}

func Test_createServer(t *testing.T) {
	handl := http.NewServeMux()
	timeout := time.Duration(time.Second)
	srv := createServer(
		handl,
		"localhost:8000",
		timeout, timeout,
	)
	assert.NotNil(t, srv)
}

func TestCreateServerFromConfigurationHolder(t *testing.T) {
	handl := http.NewServeMux()

	srv := createServerFromConfigurationHolder(handl, configuration.NewHTTPServerConfiguration())
	assert.NotNil(t, srv)
}
