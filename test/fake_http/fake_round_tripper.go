package fake_http

import "net/http"

type RoundTripper struct {
	RoundTripFn func(*http.Request) (*http.Response, error)
}

func (f *RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if f.RoundTripFn == nil {
		return nil, nil
	}
	return f.RoundTripFn(r)
}
