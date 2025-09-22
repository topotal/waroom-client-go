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

func WithRequestDump() func(http.RoundTripper) http.RoundTripper {
	return func(nextTripper http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			println("Tripper Dump Request:", r.Method, r.URL.String())
			b, err := httputil.DumpRequest(r, true) // 2つ目の引数はbodyをダンプするかどうか
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(b))
			}
			println("---")

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
