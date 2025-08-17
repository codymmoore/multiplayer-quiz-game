package api

import (
	"net/http"
)

const (
	ValidUsername = "test-username"
	ValidEmail    = "test@email.com"
	ValidPassword = "testPassword1234#?!@$%^&*-"
)

type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}
