package auth

import (
    api "common/api/auth"
    "common/test"
    "net/http"
    "testing"
)

func TestValidateLoginRequest_Success(t *testing.T) {
    request := &api.LoginRequest{
        Username: test.ValidUsername,
        Password: test.ValidPassword,
    }
    if err := ValidateLoginRequest(request); err != nil {
        t.Errorf(`ValidateLoginRequest(request) error = "%v", expected "<nil>"`, err)
    }
}

func TestValidateLoginRequest_InvalidUsername(t *testing.T) {
    request := &api.LoginRequest{
        Username: "!@#$%^&*()",
        Password: test.ValidPassword,
    }
    err := ValidateLoginRequest(request)
    if err == nil {
        t.Error(`ValidateLoginRequest(request) error = "<nil>", expected non-nil`)
    }
    test.AssertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateLoginRequest_InvalidPassword(t *testing.T) {
    request := &api.LoginRequest{
        Username: test.ValidUsername,
        Password: "invalid password",
    }
    err := ValidateLoginRequest(request)
    if err == nil {
        t.Error(`ValidateLoginRequest(request) error = "<nil>", expected non-nil`)
    }
    test.AssertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateSendVerificationEmailRequest_Success(t *testing.T) {
    request := &api.SendVerificationEmailRequest{
        Email: test.ValidEmail,
    }
    if err := ValidateSendVerificationEmailRequest(request); err != nil {
        t.Errorf(`ValidateSendVerificationEmailRequest(request) error = "%v", expected "<nil>"`, err)
    }
}

func TestValidateSendVerificationEmailRequest_InvalidEmail(t *testing.T) {
    request := api.SendVerificationEmailRequest{
        Email: "uh",
    }
    err := ValidateSendVerificationEmailRequest(&request)
    if err == nil {
        t.Error(`ValidateSendVerificationEmailRequest(request) error = "<nil>", expected non-nil`)
    }
    test.AssertHTTPError(t, err, http.StatusBadRequest)
}
