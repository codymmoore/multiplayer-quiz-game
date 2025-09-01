package test

import (
	"net/http"
)

// MockRoundTripper RoundTripper implementation used for mocking http.Client
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}
