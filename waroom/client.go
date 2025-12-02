package waroom

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/topotal/waroom-client-go/internal/httpx"
	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

type Option func(*config)

type config struct {
	BaseURL      string
	HTTPClient   *http.Client
	UserAgent    string
	RequestDump  bool
	ResponseDump bool
	APIKey       string
}

func WithBaseURL(u string) Option          { return func(c *config) { c.BaseURL = u } }
func WithHTTPClient(h *http.Client) Option { return func(c *config) { c.HTTPClient = h } }
func WithUserAgent(ua string) Option       { return func(c *config) { c.UserAgent = ua } }
func WithRequestDump() Option              { return func(c *config) { c.RequestDump = true } }
func WithResponseDump() Option             { return func(c *config) { c.ResponseDump = true } }
func WithAPIKey(apiKey string) Option      { return func(c *config) { c.APIKey = apiKey } }

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

// NewClient creates a new Waroom API client.
func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	cfg := &config{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout:   10 * time.Second, //TODO: 外部から変更可能に
			Transport: http.DefaultTransport,
		},
		UserAgent:    DefaultUserAgent,
		RequestDump:  false,
		ResponseDump: false,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// Build Http Transport Chain
	tp := cfg.HTTPClient.Transport
	var trippers []func(http.RoundTripper) http.RoundTripper

	trippers = append(trippers, httpx.WithUA(cfg.UserAgent))
	if cfg.RequestDump {
		trippers = append(trippers, httpx.WithRequestDump())
	}

	if cfg.ResponseDump {
		trippers = append(trippers, httpx.WithResponseDump())
	}

	cfg.HTTPClient.Transport = buildTripperChain(tp, trippers...)

	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("WAROOM_API_KEY")
	}
	if apiKey == "" {
		return nil, errors.New("API key is required: set it via WithAPIKey option or WAROOM_API_KEY environment variable")
	}

	auth := waroomBearerAuth{apiKey: apiKey}
	raw, err := gen.NewClient(cfg.BaseURL, &auth, gen.WithClient(cfg.HTTPClient))

	if err != nil {
		return nil, err
	}
	return &Client{raw: raw}, nil
}
