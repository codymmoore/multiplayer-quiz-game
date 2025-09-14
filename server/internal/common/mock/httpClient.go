package mock

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

// NewHttpClient creates http.Client that returns the specified status code and response body
func NewHttpClient(statusCode int, responseBody string) *http.Client {
	return &http.Client{
		Transport: &RoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: statusCode,
					Body: io.NopCloser(
						bytes.NewBufferString(responseBody),
					),
					Header: make(http.Header),
				}, nil
			},
		},
	}
}

// NewErrorHttpClient creates http.Client that returns an error
func NewErrorHttpClient() *http.Client {
	return &http.Client{
		Transport: &RoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("")
			},
		},
	}
}

// RoundTripper implementation used for mocking http.Client
type RoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}
