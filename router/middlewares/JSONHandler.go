package middlewares

import (
	"net/http"

	"github.com/ditrit/badaas/httperrors"
)

// This handler return a Marshable object and/or an [github.com/ditrit/badaas/services/httperrors.HTTPError]
type JSONHandler func(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
