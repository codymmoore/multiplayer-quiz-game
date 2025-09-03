package test

import (
	"bytes"
	"common/errors"
	stderrors "errors"
	"io"
	"net/http"
	"testing"
)

func AssertHTTPError(t *testing.T, err error, statusCode int) {
	var httpErr *errors.HTTP
	if ok := stderrors.As(err, &httpErr); !ok {
		t.Errorf(`errors.As(err, &httpErr) = "%v", expected "true"`, ok)
	}

	if httpErr.StatusCode != statusCode {
		t.Errorf(`httpErr.StatusCode = "%d", expected "%d"`, httpErr.StatusCode, statusCode)
	}
}

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
				return nil, stderrors.New("")
			},
		},
	}
}
