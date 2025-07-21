package user

import (
	"context"
	"net/http"
	"testing"
	"user/dto"
)

func TestValidateCreateUserRequest_Success(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return nil, nil
		},
	}
	request := dto.CreateUserRequest{
		Username: ValidUsername,
		Email:    ValidEmail,
		Password: ValidPassword,
	}

	if err := ValidateCreateUserRequest(&request, service, nil); err != nil {
		t.Errorf(`ValidateCreateUserRequest(&request, service, nil) = "%v", expected "<nil>"`, err)
	}
}

func TestValidateGetUserRequest_Success(t *testing.T) {
	userId := 1
	username := ValidUsername
	email := ValidEmail

	request := dto.GetUserRequest{
		UserId:   &userId,
		Username: &username,
		Email:    &email,
	}

	if err := ValidateGetUserRequest(&request); err != nil {
		t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "<nil>"`, err)
	}
}

func TestValidateGetUserRequest_MissingParameter(t *testing.T) {
	request := dto.GetUserRequest{}
	err := ValidateGetUserRequest(&request)
	if err == nil {
		t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "ID, username, or email is required"`, err)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateGetUserRequest_InvalidUserId(t *testing.T) {
	userId := -1
	request := dto.GetUserRequest{
		UserId: &userId,
	}

	err := ValidateGetUserRequest(&request)
	if err == nil {
		t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "invalid user id"`, err)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateGetUsersRequest_Success(t *testing.T) {
	limit := 1
	offset := 1
	sortField := "CreatedAt"
	sortDirection := "desc"

	request := dto.GetUsersRequest{
		Limit:         &limit,
		Offset:        &offset,
		SortField:     &sortField,
		SortDirection: &sortDirection,
	}

	if err := ValidateGetUsersRequest(&request); err != nil {
		t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "<nil>"`, err)
	}
}

func TestValidateGetUsersRequest_InvalidLimit(t *testing.T) {
	limit := -1

	request := dto.GetUsersRequest{
		Limit: &limit,
	}

	err := ValidateGetUsersRequest(&request)
	if err == nil {
		t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "limit must be a positive number"`, err)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateGetUsersRequest_InvalidOffset(t *testing.T) {
	offset := -1

	request := dto.GetUsersRequest{
		Offset: &offset,
	}

	err := ValidateGetUsersRequest(&request)
	if err == nil {
		t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "offset must be a positive number"`, err)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateGetUsersRequest_InvalidSortField(t *testing.T) {
	sortField := "invalidField"

	request := dto.GetUsersRequest{
		SortField: &sortField,
	}

	err := ValidateGetUsersRequest(&request)
	if err == nil {
		t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "invalid sort field"`, err)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateGetUsersRequest_InvalidSortDirection(t *testing.T) {
	sortDirection := "invalidDirection"

	request := dto.GetUsersRequest{
		SortDirection: &sortDirection,
	}

	err := ValidateGetUsersRequest(&request)
	if err == nil {
		t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "invalid sort direction"`, err)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateUpdateUserRequest_Success(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			if request.UserId != nil {
				return &dto.GetUserResponse{}, nil
			}
			return nil, nil
		},
	}

	username := ValidUsername
	email := ValidEmail
	password := ValidPassword

	request := dto.UpdateUserRequest{
		UserId:   1,
		Username: &username,
		Email:    &email,
		Password: &password,
	}

	if err := ValidateUpdateUserRequest(&request, service, nil); err != nil {
		t.Errorf(`ValidateUpdateUserRequest(&request, service, nil) = "%v", expected "<nil>"`, err)
	}
}

func TestValidateUpdateUserRequest_InvalidUserId(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			if request.UserId != nil {
				return &dto.GetUserResponse{}, nil
			}
			return nil, nil
		},
	}

	username := ValidUsername
	email := ValidEmail
	password := ValidPassword

	request := dto.UpdateUserRequest{
		UserId:   -1,
		Username: &username,
		Email:    &email,
		Password: &password,
	}

	err := ValidateUpdateUserRequest(&request, service, nil)
	if err == nil {
		t.Errorf(`ValidateUpdateUserRequest(&request, service, nil) = "%v", expected "invalid user id"`, err)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateUpdateUserRequest_UserNotFound(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return nil, nil
		},
	}

	username := ValidUsername
	email := ValidEmail
	password := ValidPassword

	request := dto.UpdateUserRequest{
		UserId:   1,
		Username: &username,
		Email:    &email,
		Password: &password,
	}

	err := ValidateUpdateUserRequest(&request, service, nil)
	if err == nil {
		t.Errorf(`ValidateUpdateUserRequest(&request, service, nil) = "%v", expected "user not found"`, err)
	}
	assertHTTPError(t, err, http.StatusNotFound)
}

func TestValidateDeleteUserRequest_Success(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return &dto.GetUserResponse{}, nil
		},
	}

	request := dto.DeleteUserRequest{
		UserId: 1,
	}

	if err := ValidateDeleteUserRequest(&request, service, nil); err != nil {
		t.Errorf(`ValidateDeleteUserRequest(&request, service, nil) = "%v", expected "<nil>"`, err)
	}
}

func TestValidateDeleteUserRequest_InvalidUserId(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return &dto.GetUserResponse{}, nil
		},
	}

	request := dto.DeleteUserRequest{
		UserId: -1,
	}

	err := ValidateDeleteUserRequest(&request, service, nil)
	if err == nil {
		t.Errorf(`ValidateDeleteUserRequest(&request, service, nil) = "%v", expected "invalid user id"`, err)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateDeleteUserRequest_UserNotFound(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return nil, nil
		},
	}

	request := dto.DeleteUserRequest{
		UserId: 1,
	}

	err := ValidateDeleteUserRequest(&request, service, nil)
	if err == nil {
		t.Errorf(`ValidateDeleteUserRequest(&request, service, nil) = "%v", expected "user not found"`, err)
	}
	assertHTTPError(t, err, http.StatusNotFound)
}

func TestValidateUsername_Success(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return nil, nil
		},
	}

	username := ValidUsername
	if err := validateUsername(username, service, nil); err != nil {
		t.Errorf(`validateUsername("%s", service, nil) = "%v", expected "<nil>"`, username, err)
	}
}

func TestValidateUsername_Short(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return nil, nil
		},
	}

	username := "uh"
	err := validateUsername(username, service, nil)
	if err == nil {
		t.Errorf(
			`validateUsername("%s", service, nil) = "%v", expected "username must be between %d and %d characters"`,
			username,
			err,
			MinUsernameLength,
			MaxUsernameLength,
		)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateUsername_Long(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return nil, nil
		},
	}

	username := "testInvalidUsername"
	err := validateUsername(username, service, nil)
	if err == nil {
		t.Errorf(
			`validateUsername("%s", service, nil) = "%v", expected "username must be between %d and %d characters"`,
			username,
			err,
			MinUsernameLength,
			MaxUsernameLength,
		)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateUsername_IllegalCharacter(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return nil, nil
		},
	}

	username := "username*"
	err := validateUsername(username, service, nil)
	if err == nil {
		t.Errorf(
			`validateUsername("%s", service, nil) = "%v", expected "illegal character. username must contain only letters, numbers, underscores, and hyphens"`,
			username,
			err,
		)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateUsername_Duplicate(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return &dto.GetUserResponse{}, nil
		},
	}

	username := ValidUsername
	err := validateUsername(username, service, nil)
	if err == nil {
		t.Errorf(`validateUsername("%s", service, nil) = "%v", expected "username already exists"`, username, err)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateEmail_Success(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return nil, nil
		},
	}

	email := ValidEmail
	if err := validateEmail(email, service, nil); err != nil {
		t.Errorf(`validateEmail("%s", service, nil) = "%v", expected "<nil>"`, email, err)
	}
}

func TestValidateEmail_Format(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return nil, nil
		},
	}

	email := "invalidEmail"
	err := validateEmail(email, service, nil)
	if err == nil {
		t.Errorf(`validateEmail("%s", service, nil) = "%v", expected "invalid email format"`, email, err)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateEmail_Duplicate(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
			return &dto.GetUserResponse{}, nil
		},
	}

	email := ValidEmail
	err := validateEmail(email, service, nil)
	if err == nil {
		t.Errorf(`validateEmail("%s", service, nil) = "%v", expected "email already exists"`, email, err)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidatePassword_Success(t *testing.T) {
	password := ValidPassword
	if err := validatePassword(password); err != nil {
		t.Errorf(`validatePassword("%s") = "%v", expected "<nil>"`, password, err)
	}
}

func TestValidatePassword_Short(t *testing.T) {
	password := "Invalid123456!"
	err := validatePassword(password)
	if err == nil {
		t.Errorf(
			`validatePassword("%s") = "%v", expected "password must be between %d and %d characters"`,
			password,
			err,
			MinPasswordLength,
			MaxPasswordLength,
		)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidatePassword_Long(t *testing.T) {
	password := "InvalidPassword-1234567890123457890123456789012345678901234567-1!"
	err := validatePassword(password)
	if err == nil {
		t.Errorf(
			`validatePassword("%s") = "%v", expected "password must be between %d and %d characters"`,
			password,
			err,
			MinPasswordLength,
			MaxPasswordLength,
		)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidatePassword_MissingUpper(t *testing.T) {
	password := "testpassword1234!"
	err := validatePassword(password)
	if err == nil {
		t.Errorf(
			`validatePassword("%s") = "%v", expected "invalid password. Must contain at least one of the following: upper case English character, lower case English character, number, special character"`,
			password,
			err,
		)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidatePassword_MissingLower(t *testing.T) {
	password := "TESTPASSWORD1234!"
	err := validatePassword(password)
	if err == nil {
		t.Errorf(
			`validatePassword("%s") = "%v", expected "invalid password. Must contain at least one of the following: upper case English character, lower case English character, number, special character"`,
			password,
			err,
		)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidatePassword_MissingNumber(t *testing.T) {
	password := "testPassword#?!@$%^&*-"
	err := validatePassword(password)
	if err == nil {
		t.Errorf(
			`validatePassword("%s") = "%v", expected "invalid password. Must contain at least one of the following: upper case English character, lower case English character, number, special character"`,
			password,
			err,
		)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidatePassword_MissingSymbol(t *testing.T) {
	password := "testPassword12345"
	err := validatePassword(password)
	if err == nil {
		t.Errorf(
			`validatePassword("%s") = "%v", expected "invalid password. Must contain at least one of the following: upper case English character, lower case English character, number, special character"`,
			password,
			err,
		)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidatePassword_IllegalCharacter(t *testing.T) {
	password := ValidPassword + "😘"
	err := validatePassword(password)
	if err == nil {
		t.Errorf(
			`validatePassword("%s") = "%v", expected "password contains illegal characters"`,
			password,
			err,
		)
	}
	assertHTTPError(t, err, http.StatusBadRequest)
}
