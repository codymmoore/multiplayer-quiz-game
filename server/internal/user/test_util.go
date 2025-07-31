package user

import (
	"common"
	api "common/api/user"
	"context"
	"errors"
	"testing"
)

const (
	ValidUsername = "test-username"
	ValidEmail    = "test@email.com"
	ValidPassword = "testPassword1234#?!@$%^&*-"
)

type mockService struct {
	createUserFunc func(context context.Context, request *api.CreateUserRequest) (
		*api.CreateUserResponse,
		error,
	)
	getUserFunc    func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error)
	getUsersFunc   func(context context.Context, request *api.GetUsersRequest) (*api.GetUsersResponse, error)
	updateUserFunc func(context context.Context, request *api.UpdateUserRequest) (
		*api.UpdateUserResponse,
		error,
	)
	deleteUserFunc func(context context.Context, request *api.DeleteUserRequest) (
		*api.DeleteUserResponse,
		error,
	)
}

func (m *mockService) CreateUser(context context.Context, request *api.CreateUserRequest) (
	*api.CreateUserResponse,
	error,
) {
	return m.createUserFunc(context, request)
}

func (m *mockService) GetUser(context context.Context, request *api.GetUserRequest) (
	*api.GetUserResponse,
	error,
) {
	return m.getUserFunc(context, request)
}

func (m *mockService) GetUsers(ctx context.Context, request *api.GetUsersRequest) (
	*api.GetUsersResponse,
	error,
) {
	return m.getUsersFunc(ctx, request)
}

func (m *mockService) UpdateUser(context context.Context, request *api.UpdateUserRequest) (
	*api.UpdateUserResponse,
	error,
) {
	return m.updateUserFunc(context, request)
}

func (m *mockService) DeleteUser(context context.Context, request *api.DeleteUserRequest) (
	*api.DeleteUserResponse,
	error,
) {
	return m.deleteUserFunc(context, request)
}

func assertHTTPError(t *testing.T, err error, statusCode int) {
	var httpErr *common.HTTPError
	if ok := errors.As(err, &httpErr); !ok {
		t.Errorf(`errors.As(err, &httpErr) = "%v", expected "true"`, ok)
	}

	if httpErr.StatusCode != statusCode {
		t.Errorf(`httpErr.StatusCode = "%d", expected "%d"`, httpErr.StatusCode, statusCode)
	}
}
