package user

import (
    "common"
    "context"
    "fmt"
    "net/http"
    "reflect"
    "regexp"
    "strings"
    "user/db/generated"
    "user/dto"
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

// ValidateCreateUserRequest Validate request for creating a new user
func ValidateCreateUserRequest(request *dto.CreateUserRequest, service Service, context context.Context) error {
    if err := validateUsername(request.Username, service, context); err != nil {
        return err
    }

    if err := validateEmail(request.Email, service, context); err != nil {
        return err
    }

    if err := validatePassword(request.Password); err != nil {
        return err
    }

    return nil
}

// ValidateGetUserRequest Validate request for retrieving a user
func ValidateGetUserRequest(request *dto.GetUserRequest) error {
    if request.UserId == nil && request.Username == nil && request.Email == nil {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "ID, username, or email is required",
        }
    }

    if request.UserId != nil && *request.UserId < 0 {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "invalid user id",
        }
    }

    return nil
}

// ValidateGetUsersRequest Validate request for retrieving paginated users
func ValidateGetUsersRequest(request *dto.GetUsersRequest) error {
    if request.Limit != nil && *request.Limit <= 0 {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "limit must be a positive number",
        }
    }

    if request.Offset != nil && *request.Offset < 0 {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "offset must be a positive number",
        }
    }

    userType := reflect.TypeOf(db.User{})
    if request.SortField != nil {
        sortFieldExists := false
        for i := 0; i < userType.NumField(); i++ {
            field := userType.Field(i)
            if strings.EqualFold(field.Name, *request.SortField) {
                sortFieldExists = true
                break
            }
        }

        if !sortFieldExists {
            return &common.HTTPError{
                StatusCode: http.StatusBadRequest,
                Message:    "invalid sort field",
            }
        }
    }

    if request.SortDirection != nil && (!strings.EqualFold(
        *request.SortDirection,
        "desc",
    ) && !strings.EqualFold(*request.SortDirection, "asc")) {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "invalid sort direction",
        }
    }

    return nil
}

// ValidateUpdateUserRequest Validate request for updating a user
func ValidateUpdateUserRequest(request *dto.UpdateUserRequest, service Service, context context.Context) error {
    if request.UserId < 0 {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "invalid user id",
        }
    }

    getUserRequest := dto.GetUserRequest{UserId: &request.UserId}
    if response, _ := service.GetUser(context, &getUserRequest); response == nil {
        return &common.HTTPError{
            StatusCode: http.StatusNotFound,
            Message:    "user not found",
        }
    }

    if request.Username != nil {
        if err := validateUsername(*request.Username, service, context); err != nil {
            return err
        }
    }

    if request.Email != nil {
        if err := validateEmail(*request.Email, service, context); err != nil {
            return err
        }
    }

    if request.Password != nil {
        if err := validatePassword(*request.Password); err != nil {
            return err
        }
    }

    return nil
}

// ValidateDeleteUserRequest Validate request for deleting a user
func ValidateDeleteUserRequest(request *dto.DeleteUserRequest, service Service, context context.Context) error {
    if request.UserId < 0 {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "invalid user id",
        }
    }

    getUserRequest := dto.GetUserRequest{UserId: &request.UserId}
    if response, _ := service.GetUser(context, &getUserRequest); response == nil {
        return &common.HTTPError{
            StatusCode: http.StatusNotFound,
            Message:    "user not found",
        }
    }

    return nil
}

// validateUsername Validate a username
func validateUsername(username string, service Service, context context.Context) error {
    if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message: fmt.Sprintf(
                "username must be between %d and %d characters",
                MinUsernameLength,
                MaxUsernameLength,
            ),
        }
    }

    if !usernameRegex.MatchString(username) {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "illegal character. username must contain only letters, numbers, underscores, and hyphens",
        }
    }

    getUserRequest := dto.GetUserRequest{Username: &username}
    if response, _ := service.GetUser(context, &getUserRequest); response != nil {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "username already exists",
        }
    }

    return nil
}

// validateEmail Validate an email address
func validateEmail(email string, service Service, context context.Context) error {
    if !emailRegex.MatchString(email) {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "invalid email format",
        }
    }

    getUserRequest := dto.GetUserRequest{Email: &email}
    if response, _ := service.GetUser(context, &getUserRequest); response != nil {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "email already exists",
        }
    }

    return nil
}

// validatePassword Validate a password
func validatePassword(password string) error {
    if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message: fmt.Sprintf(
                "password must be between %d and %d characters",
                MinPasswordLength,
                MaxPasswordLength,
            ),
        }
    }

    if !hasUpper.MatchString(password) || !hasLower.MatchString(password) || !hasNumber.MatchString(password) || !hasSymbol.MatchString(password) {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message: "invalid password. Must contain at least one of each of the following: upper" +
                " case English character, lower case English character, number, special character",
        }
    }

    if hasIllegalCharacters.MatchString(password) {
        return &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "password contains illegal characters",
        }
    }

    return nil
}
