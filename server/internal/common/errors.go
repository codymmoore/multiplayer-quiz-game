package common

// HTTPError Custom error type containing HTTP status code
type HTTPError struct {
	StatusCode int
	Message    string
}

// Error Error() implementation from error interface
func (e *HTTPError) Error() string {
	return e.Message
}
