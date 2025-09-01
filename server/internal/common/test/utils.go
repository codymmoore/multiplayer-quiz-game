package test

import (
	"common/errors"
	stderrors "errors"
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
