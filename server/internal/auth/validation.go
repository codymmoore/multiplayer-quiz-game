package auth

import (
	api "common/api/auth"
	"common/errors"
	"common/validate"
	"net/http"
)

// ValidateLoginRequest validates a login request
func ValidateLoginRequest(request *api.LoginRequest) error {
	if err := validate.Username(request.Username); err != nil {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid username",
		}
	}
	if err := validate.Password(request.Password); err != nil {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid password",
		}
	}
	return nil
}

// ValidateSendVerificationEmailRequest validates a send verification email request
func ValidateSendVerificationEmailRequest(request *api.SendVerificationEmailRequest) error {
	if err := validate.Email(request.Email); err != nil {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid email",
		}
	}
	return nil
}
