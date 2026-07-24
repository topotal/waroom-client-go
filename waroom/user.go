package waroom

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// UserService handles communication with the user related
// methods of the Waroom API. Use it via Client.Users.
type UserService service

// NewUserService returns client.Users.
//
// Deprecated: Use Client.Users instead.
func NewUserService(client *Client) *UserService {
	return client.Users
}

func (s *UserService) List(ctx context.Context, params ListUsersParams) ([]User, error) {
	res, err := s.client.raw.ListUsers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("waroom: list users: %w", err)
	}

	switch res := res.(type) {
	case *gen.ListUsersOK:
		return res.Users, nil
	case *gen.ListUsersUnauthorized:
		return nil, fmt.Errorf("waroom: list users: %w", &APIError{StatusCode: http.StatusUnauthorized, Messages: res.Errors})
	case *gen.ListUsersPaymentRequired:
		return nil, fmt.Errorf("waroom: list users: %w", &APIError{StatusCode: http.StatusPaymentRequired, Messages: []string{res.Error}})
	default:
		return nil, fmt.Errorf("waroom: list users: unhandled response %T", res)
	}
}
