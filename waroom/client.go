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
	raw    *gen.Client
	common service

	Incidents *IncidentService
	Services  *ServiceService
	Labels    *LabelService
	Users     *UserService
}

// service holds the state shared by all resource services. Each service
// type is a cast of Client.common, so every service points back to the
// same client instance.
type service struct {
	client *Client
}

// The API spec includes the /api/v0 prefix in each path, so the base URL
// contains the host only.
const DefaultBaseURL = "https://api.staging.waroom.com"
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
			Timeout:   10 * time.Second, // TODO: make the timeout configurable
			Transport: http.DefaultTransport,
		},
		UserAgent:    DefaultUserAgent,
		RequestDump:  false,
		ResponseDump: false,
	}
	for _, opt := range opts {
		opt(cfg)
	}

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

	c := &Client{raw: raw}
	c.common.client = c
	c.Incidents = (*IncidentService)(&c.common)
	c.Services = (*ServiceService)(&c.common)
	c.Labels = (*LabelService)(&c.common)
	c.Users = (*UserService)(&c.common)
	return c, nil
}
