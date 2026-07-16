package waroom

import (
	"fmt"
	"net/http"
	"strings"
)

// APIError represents an error response returned by the Waroom API server.
// Transport-layer errors (connection or decode failures) are never APIErrors.
type APIError struct {
	StatusCode int      // 400, 401, 402, 404, or 422
	Messages   []string // normalized from the response body's error / errors field
}

func (e *APIError) Error() string {
	msg := fmt.Sprintf("API error %d", e.StatusCode)
	if text := http.StatusText(e.StatusCode); text != "" {
		msg += " " + text
	}
	if len(e.Messages) > 0 {
		msg += ": " + strings.Join(e.Messages, "; ")
	}
	return msg
}

func optMessages(s OptString) []string {
	if v, ok := s.Get(); ok {
		return []string{v}
	}
	return nil
}
