package waroom

import (
	"context"
	"net/http"
	"time"

	"waroom-client-go/internal/httpx"
	gen "waroom-client-go/internal/ogen"
)

type Option func(*config)

type config struct {
	BaseURL    string
	HTTPClient *http.Client
	TokenSrc   TokenSource
	UserAgent  string
	Retry      *httpx.RetryPolicy
}

func WithBaseURL(u string) Option           { return func(c *config) { c.BaseURL = u } }
func WithHTTPClient(h *http.Client) Option  { return func(c *config) { c.HTTPClient = h } }
func WithTokenSource(ts TokenSource) Option { return func(c *config) { c.TokenSrc = ts } }
func WithUserAgent(ua string) Option        { return func(c *config) { c.UserAgent = ua } }
func WithRetryPolicy(p httpx.RetryPolicy) Option {
	return func(c *config) { c.Retry = &p }
}

type Client struct {
	raw *gen.Client
}

const DefaultBaseURL = "https://api.app.waroom.com/api/v0"

func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	cfg := &config{
		BaseURL: "https://api.app.waroom.com/api/v0",
		HTTPClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: http.DefaultTransport,
		},
	}
	for _, o := range opts {
		o(cfg)
	}

	// middleware chain
	transport := cfg.HTTPClient.Transport
	var mws []func(http.RoundTripper) http.RoundTripper
	if cfg.TokenSrc != nil {
		mws = append(mws, httpx.WithAuthBearer(cfg.TokenSrc.Token))
	}
	if cfg.Retry != nil {
		mws = append(mws, httpx.WithRetry(*cfg.Retry))
	}
	cfg.HTTPClient.Transport = httpx.Chain(transport, mws...)

	// ogen 側の Client を作成（生成物のAPIに合わせて調整）
	raw, err := gen.NewClient(DefaultBaseURL, )
	if err != nil {
		return nil, err
	}
	return &Client{raw: raw}, nil
}

// 便利メソッド例（下は概念。生成された名前に合わせて薄くラップ）
func (c *Client) CreateIncident(ctx context.Context, in CreateIncidentRequest) (*Incident, error) {
	out, err := c.raw.CreateIncident(ctx, gen.CreateIncidentRequest(in))
	return (*Incident)(out), err
}
