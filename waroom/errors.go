package waroom

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// Sentinel errors usable with errors.Is to detect API error responses by
// status code, e.g. errors.Is(err, waroom.ErrNotFound).
var (
	ErrBadRequest          = errors.New("waroom: bad request")
	ErrUnauthorized        = errors.New("waroom: unauthorized")
	ErrPaymentRequired     = errors.New("waroom: payment required")
	ErrNotFound            = errors.New("waroom: not found")
	ErrUnprocessableEntity = errors.New("waroom: unprocessable entity")
)

// APIError represents an error response from the Waroom API. Use
// errors.As to inspect it, or errors.Is with the Err* sentinels to match
// by status code.
type APIError struct {
	StatusCode  int
	OperationID string
	Messages    []string
}

func (e *APIError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "waroom: %s: %d %s", e.OperationID, e.StatusCode, http.StatusText(e.StatusCode))
	if len(e.Messages) > 0 {
		b.WriteString(": ")
		b.WriteString(strings.Join(e.Messages, "; "))
	}
	return b.String()
}

// Is maps StatusCode to the sentinel errors above.
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrBadRequest:
		return e.StatusCode == 400
	case ErrUnauthorized:
		return e.StatusCode == 401
	case ErrPaymentRequired:
		return e.StatusCode == 402
	case ErrNotFound:
		return e.StatusCode == 404
	case ErrUnprocessableEntity:
		return e.StatusCode == 422
	}
	return false
}

// newAPIError is the only error constructor generated code calls. It
// extracts messages from any of the error body shapes the Waroom API
// uses, via the getters ogen generates:
//
//	errors: [string]  → GetErrors() []string
//	error: string?    → GetError() OptString
//	error: string     → GetError() string
func newAPIError(operationID string, statusCode int, body any) *APIError {
	var messages []string
	switch b := body.(type) {
	case interface{ GetErrors() []string }:
		messages = b.GetErrors()
	case interface{ GetError() gen.OptString }:
		if v, ok := b.GetError().Get(); ok && v != "" {
			messages = []string{v}
		}
	case interface{ GetError() string }:
		if v := b.GetError(); v != "" {
			messages = []string{v}
		}
	}
	return &APIError{
		StatusCode:  statusCode,
		OperationID: operationID,
		Messages:    messages,
	}
}
