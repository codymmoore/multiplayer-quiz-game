package validate

import (
    "common/errors"
    "fmt"
    "net/http"
    "regexp"
)

const (
    MinUsernameLength = 3
    MaxUsernameLength = 15
    MinPasswordLength = 15
    MaxPasswordLength = 64
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,15}$`)
var emailRegex = regexp.MustCompile(`^[\w\-.]+@([\w-]+\.)+[\w-]{2,}$`)
var hasUpper = regexp.MustCompile(`[A-Z]`)
var hasLower = regexp.MustCompile(`[a-z]`)
var hasNumber = regexp.MustCompile(`[0-9]`)
var hasSymbol = regexp.MustCompile(`[#?!@$%^&*-]`)
var hasValidCharacters = regexp.MustCompile(`^[ -~]+$`)
var hasIllegalCharacters = regexp.MustCompile(`[^ -~]`)

// Username validates a username
func Username(username string) error {
    if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
        return &errors.HTTP{
            StatusCode: http.StatusBadRequest,
            Message: fmt.Sprintf(
                "username must be between %d and %d characters",
                MinUsernameLength,
                MaxUsernameLength,
            ),
        }
    }

    if !usernameRegex.MatchString(username) {
        return &errors.HTTP{
            StatusCode: http.StatusBadRequest,
            Message:    "illegal character. username must contain only letters, numbers, underscores, and hyphens",
        }
    }

    return nil
}

// Email validates an email
func Email(email string) error {
    if !emailRegex.MatchString(email) {
        return &errors.HTTP{
            StatusCode: http.StatusBadRequest,
            Message:    "invalid email format",
        }
    }
    return nil
}

// Password validates a password
func Password(password string) error {
    if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
        return &errors.HTTP{
            StatusCode: http.StatusBadRequest,
            Message: fmt.Sprintf(
                "password must be between %d and %d characters",
                MinPasswordLength,
                MaxPasswordLength,
            ),
        }
    }

    if !hasUpper.MatchString(password) || !hasLower.MatchString(password) || !hasNumber.MatchString(password) || !hasSymbol.MatchString(password) {
        return &errors.HTTP{
            StatusCode: http.StatusBadRequest,
            Message: "invalid password. Must contain at least one of each of the following: upper" +
                " case English character, lower case English character, number, special character",
        }
    }

    if hasIllegalCharacters.MatchString(password) {
        return &errors.HTTP{
            StatusCode: http.StatusBadRequest,
            Message:    "password contains illegal characters",
        }
    }

    return nil
}
