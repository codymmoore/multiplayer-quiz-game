package mock

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

// GetMockHttpClient gets http.Client that returns the specified status code and response body
func GetMockHttpClient(statusCode int, responseBody string) *http.Client {
	return &http.Client{
		Transport: &MockRoundTripper{
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

// GetMockErrorHttpClient gets http.Client that returns an error
func GetMockErrorHttpClient() *http.Client {
	return &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("")
			},
		},
	}
}

// MockRoundTripper RoundTripper implementation used for mocking http.Client
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}
