package user

import (
	"errors"
	"fmt"
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

// ValidateCreateUserRequest Validate request for creating a new user
// TODO: check is username or email is already being used?
func ValidateCreateUserRequest(request *dto.CreateUserRequest) error {
	err := ValidateUsername(request.Username)
	if err != nil {
		return err
	}

	err = ValidateEmail(request.Email)
	if err != nil {
		return err
	}

	err = ValidatePassword(request.Password)
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetUserRequest Validate request for retrieving a user
func ValidateGetUserRequest(request *dto.GetUserRequest) error {
	if request.UserId == nil && request.Username == nil && request.Email == nil {
		return errors.New("ID, username, or email is required")
	}
	return nil
}

// ValidateGetUsersRequest Validate request for retrieving paginated users
func ValidateGetUsersRequest(request *dto.GetUsersRequest) error {
	userType := reflect.TypeOf(db.User{})

	if request.Limit != nil && *request.Limit <= 0 {
		return errors.New("limit must be greater than zero")
	}

	if request.Offset != nil && *request.Offset < 0 {
		return errors.New("offset must be a positive number")
	}

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
			return errors.New("invalid sort field")
		}
	}

	if request.SortField != nil && (!strings.EqualFold(
		*request.SortField,
		"desc",
	) || !strings.EqualFold(*request.SortField, "asc")) {
		return errors.New("invalid sort direction")
	}

	return nil
}

// ValidateUpdateUserRequest Validate request for updating a user
func ValidateUpdateUserRequest(request *dto.UpdateUserRequest) error {
	if request.Username != nil {
		err := ValidateUsername(*request.Username)
		if err != nil {
			return err
		}
	}

	if request.Email != nil {
		err := ValidateEmail(*request.Email)
		if err != nil {
			return err
		}
	}

	if request.Password != nil {
		err := ValidatePassword(*request.Password)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateUsername Validate a username
func ValidateUsername(username string) error {
	if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
		return fmt.Errorf("username must be between %d and %d characters", MinUsernameLength, MaxUsernameLength)
	}

	if !usernameRegex.MatchString(username) {
		return errors.New("invalid username format. username must contain only letters, numbers, underscores, and hyphens")
	}

	return nil
}

// ValidateEmail Validate an email address
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}

// ValidatePassword Validate a password
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
		return fmt.Errorf("password must be between %d and %d characters", MinPasswordLength, MaxPasswordLength)
	}
	return nil
}
