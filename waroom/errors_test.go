package waroom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *APIError
		want string
	}{
		{
			name: "with single message",
			err:  &APIError{StatusCode: 404, Messages: []string{"label not found"}},
			want: "API error 404 Not Found: label not found",
		},
		{
			name: "with multiple messages",
			err:  &APIError{StatusCode: 401, Messages: []string{"Not Authenticated", "token expired"}},
			want: "API error 401 Unauthorized: Not Authenticated; token expired",
		},
		{
			name: "without messages",
			err:  &APIError{StatusCode: 402},
			want: "API error 402 Payment Required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Error())
		})
	}
}

func TestOptMessages(t *testing.T) {
	assert.Equal(t, []string{"boom"}, optMessages(NewOptString("boom")))
	assert.Nil(t, optMessages(OptString{}))
}
