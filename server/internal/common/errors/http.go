package errors

// HTTP custom error type containing status code
type HTTP struct {
	StatusCode int
	Message    string
}

// Error implementation from error interface
func (e *HTTP) Error() string {
	return e.Message
}
