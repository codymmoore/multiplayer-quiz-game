package user

import (
	api "common/api/user"
	"common/errors"
	"common/validate"
	"context"
	"net/http"
	"reflect"
	"strings"
	"user/db/generated"
)

// ValidateCreateUserRequest Validate request for creating a new user
func ValidateCreateUserRequest(request *api.CreateUserRequest, service Service, context context.Context) error {
	username := request.Username
	if err := validate.Username(username); err != nil {
		return err
	}
	if !isUniqueUsername(username, service, context) {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "duplicate username",
		}
	}

	email := request.Email
	if err := validate.Email(email); err != nil {
		return err
	}
	if !isUniqueEmail(email, service, context) {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "duplicate email",
		}
	}

	if err := validate.Password(request.Password); err != nil {
		return err
	}

	return nil
}

// ValidateGetUserRequest Validate request for retrieving a user
func ValidateGetUserRequest(request *api.GetUserRequest) error {
	if request.UserId == nil && request.Username == nil && request.Email == nil {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "ID, username, or email is required",
		}
	}

	if request.UserId != nil && *request.UserId < 0 {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid user id",
		}
	}

	return nil
}

// ValidateGetUsersRequest Validate request for retrieving paginated users
func ValidateGetUsersRequest(request *api.GetUsersRequest) error {
	if request.Limit != nil && *request.Limit <= 0 {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "limit must be a positive number",
		}
	}

	if request.Offset != nil && *request.Offset < 0 {
		return &errors.HTTP{
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
			return &errors.HTTP{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid sort field",
			}
		}
	}

	if request.SortDirection != nil && (!strings.EqualFold(
		*request.SortDirection,
		"desc",
	) && !strings.EqualFold(*request.SortDirection, "asc")) {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid sort direction",
		}
	}

	return nil
}

// ValidateUpdateUserRequest Validate request for updating a user
func ValidateUpdateUserRequest(request *api.UpdateUserRequest, service Service, context context.Context) error {
	userId := request.UserId
	if userId < 0 {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid user id",
		}
	}

	getUserRequest := &api.GetUserRequest{UserId: &userId}
	if response, _ := service.GetUser(context, getUserRequest); response == nil {
		return &errors.HTTP{
			StatusCode: http.StatusNotFound,
			Message:    "user not found",
		}
	}

	if request.Username != nil {
		username := *request.Username
		if err := validate.Username(username); err != nil {
			return err
		}
		if !isUniqueUsername(username, service, context) {
			return &errors.HTTP{
				StatusCode: http.StatusBadRequest,
				Message:    "duplicate username",
			}
		}
	}

	if request.Email != nil {
		email := *request.Email
		if err := validate.Email(email); err != nil {
			return err
		}
		if !isUniqueEmail(email, service, context) {
			return &errors.HTTP{
				StatusCode: http.StatusBadRequest,
				Message:    "duplicate email",
			}
		}
	}

	if request.Password != nil {
		if err := validate.Password(*request.Password); err != nil {
			return err
		}
	}

	return nil
}

// ValidateDeleteUserRequest Validate request for deleting a user
func ValidateDeleteUserRequest(request *api.DeleteUserRequest, service Service, context context.Context) error {
	if request.UserId < 0 {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid user id",
		}
	}

	getUserRequest := api.GetUserRequest{UserId: &request.UserId}
	if response, _ := service.GetUser(context, &getUserRequest); response == nil {
		return &errors.HTTP{
			StatusCode: http.StatusNotFound,
			Message:    "user not found",
		}
	}

	return nil
}

// ValidateVerifyUserRequest validates request for verifying user
func ValidateVerifyUserRequest(request *api.VerifyUserRequest, service Service, context context.Context) error {
	if request.UserId < 0 {
		return &errors.HTTP{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid user id",
		}
	}

	getUserRequest := api.GetUserRequest{UserId: &request.UserId}
	if response, _ := service.GetUser(context, &getUserRequest); response == nil {
		return &errors.HTTP{
			StatusCode: http.StatusNotFound,
			Message:    "user not found",
		}
	}

	return nil
}

// isUniqueUsername determines if a username is already in use
func isUniqueUsername(username string, service Service, ctx context.Context) bool {
	getUserRequest := api.GetUserRequest{Username: &username}
	response, _ := service.GetUser(ctx, &getUserRequest)
	return response == nil
}

// isUniqueEmail determines if an email address is already in use
func isUniqueEmail(email string, service Service, ctx context.Context) bool {
	getUserRequest := api.GetUserRequest{Email: &email}
	response, _ := service.GetUser(ctx, &getUserRequest)
	return response == nil
}
