package edgeosclient

import (
	"net/http"
)

type csrfTransport struct {
	RoundTripper http.RoundTripper
	token        string
}

func (r *csrfTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if r.token != "" {
		request.Header.Set("X-CSRF-Token", r.token)
	}

	res, err := r.RoundTripper.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if res.Cookies() != nil {
		for _, cookie := range res.Cookies() {
			if cookie.Name == "X-CSRF-TOKEN" {
				r.token = cookie.Value
			}
		}
	}

	return res, err
}
