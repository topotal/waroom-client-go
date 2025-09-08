package httpx

import (
	"fmt"
	"net/http"
	"net/http/httputil"
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

func WithDump() func(http.RoundTripper) http.RoundTripper {
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
