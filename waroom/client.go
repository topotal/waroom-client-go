package waroom

import (
	"context"
	"net/http"
	"time"

	httpx "github.com/topotal/waroom-client-go/internal"
	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

type Option func(*config)

type config struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
	Dump       bool
}

func WithBaseURL(u string) Option          { return func(c *config) { c.BaseURL = u } }
func WithHTTPClient(h *http.Client) Option { return func(c *config) { c.HTTPClient = h } }
func WithUserAgent(ua string) Option       { return func(c *config) { c.UserAgent = ua } }
func WithDump(b bool) Option               { return func(c *config) { c.Dump = b } }

type Client struct {
	raw *gen.Client
}

const DefaultBaseURL = "https://api.staging.waroom.com/api/v0"
const DefaultUserAgent = "waroom-sdk-go/0.0.1"

func buildTripperChain(rt http.RoundTripper, trippers ...func(http.RoundTripper) http.RoundTripper) http.RoundTripper {
	for i := len(trippers) - 1; i >= 0; i-- {
		rt = trippers[i](rt)
	}
	return rt
}

func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	cfg := &config{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: http.DefaultTransport,
		},
		UserAgent: DefaultUserAgent,
		Dump:      false,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// Build Http Transport Chain
	tp := cfg.HTTPClient.Transport
	var trippers []func(http.RoundTripper) http.RoundTripper

	trippers = append(trippers, httpx.WithUA(cfg.UserAgent))
	if cfg.Dump {
		trippers = append(trippers, httpx.WithDump())
	}

	cfg.HTTPClient.Transport = buildTripperChain(tp, trippers...)

	// ogen 側の Client を作成
	auth := WaroomBearerAuth{}
	raw, err := gen.NewClient(cfg.BaseURL, &auth, gen.WithClient(cfg.HTTPClient))

	if err != nil {
		return nil, err
	}
	return &Client{raw: raw}, nil
}
