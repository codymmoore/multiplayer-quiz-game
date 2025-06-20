package user

import (
	"common"
	"context"
	"errors"
	"testing"
	"user/dto"
)

const (
	ValidUsername = "test-username"
	ValidEmail    = "test@email.com"
	ValidPassword = "testPassword1234#?!@$%^&*-"
)

type mockService struct {
	createUserFunc func(context context.Context, request *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	getUserFunc    func(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error)
	getUsersFunc   func(context context.Context, request *dto.GetUsersRequest) (*dto.GetUsersResponse, error)
	updateUserFunc func(context context.Context, request *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
	deleteUserFunc func(context context.Context, request *dto.DeleteUserRequest) (*dto.DeleteUserResponse, error)
}

func (m *mockService) CreateUser(context context.Context, request *dto.CreateUserRequest) (
	*dto.CreateUserResponse,
	error,
) {
	return m.createUserFunc(context, request)
}

func (m *mockService) GetUser(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
	return m.getUserFunc(context, request)
}

func (m *mockService) GetUsers(ctx context.Context, request *dto.GetUsersRequest) (*dto.GetUsersResponse, error) {
	return m.getUsersFunc(ctx, request)
}

func (m *mockService) UpdateUser(context context.Context, request *dto.UpdateUserRequest) (
	*dto.UpdateUserResponse,
	error,
) {
	return m.updateUserFunc(context, request)
}

func (m *mockService) DeleteUser(context context.Context, request *dto.DeleteUserRequest) (
	*dto.DeleteUserResponse,
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
