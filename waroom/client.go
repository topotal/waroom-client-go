package waroom

import (
	"context"
	"net/http"

	gen "waroom-client-go/internal/ogen"
)

type Option func(*config)

type config struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

func WithBaseURL(u string) Option          { return func(c *config) { c.BaseURL = u } }
func WithHTTPClient(h *http.Client) Option { return func(c *config) { c.HTTPClient = h } }
func WithUserAgent(ua string) Option       { return func(c *config) { c.UserAgent = ua } }

type Client struct {
	raw *gen.Client
}

const DefaultBaseURL = "https://api.staging.waroom.com/api/v0"

func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	//TODO: configあれば追加していく
	// cfg := &config{
	// 	HTTPClient: &http.Client{
	// 		Timeout:   30 * time.Second,
	// 		Transport: http.DefaultTransport,
	// 	},
	// }
	// for _, opt := range opts {
	// 	opt(cfg)
	// }

	// ogen 側の Client を作成
	auth := WaroomBearerAuth{}
	raw, err := gen.NewClient(DefaultBaseURL, &auth)

	if err != nil {
		return nil, err
	}
	return &Client{raw: raw}, nil
}
