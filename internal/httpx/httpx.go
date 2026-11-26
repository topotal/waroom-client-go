package httpx

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func Chain(rt http.RoundTripper, mws ...func(http.RoundTripper) http.RoundTripper) http.RoundTripper {
	for i := len(mws) - 1; i >= 0; i-- {
		rt = mws[i](rt)
	}
	return rt
}

func WithUA(ua string) func(http.RoundTripper) http.RoundTripper {
	return func(nextTripper http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			if r.Header.Get("User-Agent") == "" {
				r.Header.Add("User-Agent", ua)
				r.Header.Add("Accept", "application/json")
			}

			if nextTripper == nil {
				return http.DefaultTransport.RoundTrip(r)
			}
			return nextTripper.RoundTrip(r)

		})
	}
}

// sensitiveHeaders are the request headers whose values WithRequestDump
// replaces with redactedValue. Dumps end up in bug reports, so they must
// never contain credentials.
var sensitiveHeaders = []string{
	"Authorization",
	"Proxy-Authorization",
	"Cookie",
	"X-Api-Key",
}

const redactedValue = "[REDACTED]"

// redactHeader returns a copy of h with the values of sensitiveHeaders
// redacted. The auth scheme (e.g. "Bearer") is kept so the dump still shows
// how the request authenticated.
func redactHeader(h http.Header) http.Header {
	out := h.Clone()
	for _, name := range sensitiveHeaders {
		vals := out.Values(name)
		if len(vals) == 0 {
			continue
		}
		redacted := make([]string, len(vals))
		for i, v := range vals {
			if scheme, _, ok := strings.Cut(v, " "); ok && name != "Cookie" {
				redacted[i] = scheme + " " + redactedValue
			} else {
				redacted[i] = redactedValue
			}
		}
		out[http.CanonicalHeaderKey(name)] = redacted
	}
	return out
}

// WithRequestDump dumps each outgoing request (headers and body) to w.
// Credential headers are redacted; see sensitiveHeaders.
func WithRequestDump(w io.Writer) func(http.RoundTripper) http.RoundTripper {
	return func(nextTripper http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			_, _ = fmt.Fprintln(w, "Tripper Dump Request:", r.Method, r.URL.String())

			// Dump a shallow copy carrying the redacted headers. DumpRequest
			// consumes the body and replaces it with an equivalent reader,
			// so hand that reader back to the original request.
			dr := *r
			dr.Header = redactHeader(r.Header)
			b, err := httputil.DumpRequest(&dr, true)
			r.Body = dr.Body
			if err != nil {
				_, _ = fmt.Fprintln(w, err)
			} else {
				_, _ = fmt.Fprintln(w, string(b))
			}
			_, _ = fmt.Fprintln(w, "---")

			if nextTripper == nil {
				return http.DefaultTransport.RoundTrip(r)
			}
			return nextTripper.RoundTrip(r)
		})
	}
}

// NOTE RoundTripper型定義
// type RoundTripper interface {
// 	RoundTrip(*Request) (*Response, error)
// }

// WithResponseDump wraps an http.RoundTripper, dumps response body (gzip対応) and restores resp.Body.
func WithResponseDump() func(http.RoundTripper) http.RoundTripper {
	return func(nextTripper http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			start := time.Now()

			resp, err := nextTripper.RoundTrip(r)
			if err != nil {
				return nil, err
			}

			// Read body (gzip-aware) and restore resp.Body
			var raw []byte
			if resp.Body != nil {
				var rd io.Reader = resp.Body
				if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
					if gr, gzErr := gzip.NewReader(resp.Body); gzErr == nil {
						defer gr.Close()
						rd = gr
					}
				}
				raw, _ = io.ReadAll(rd)
				_ = resp.Body.Close()
				resp.Body = io.NopCloser(bytes.NewReader(raw))
			}

			// Summary line
			fmt.Fprintf(os.Stderr, "[%s] %s %s -> %s (%d bytes, %s)\n",
				time.Now().Format(time.RFC3339), r.Method, r.URL, resp.Status,
				len(raw), time.Since(start).Truncate(time.Millisecond))

			// Dump body: pretty-print if JSON, otherwise print text/* as-is
			if len(raw) > 0 {
				ct := resp.Header.Get("Content-Type")
				if strings.HasPrefix(ct, "application/json") {
					var buf bytes.Buffer
					if json.Indent(&buf, raw, "", "  ") == nil {
						fmt.Fprintln(os.Stderr, buf.String())
					} else {
						fmt.Fprintln(os.Stderr, string(raw))
					}
				} else if strings.HasPrefix(ct, "text/") {
					fmt.Fprintln(os.Stderr, string(raw))
				}
			}

			if nextTripper == nil {
				nextTripper = http.DefaultTransport
			}

			return resp, nil
		})
	}
}
