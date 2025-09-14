package mock

import (
	"common/user"
)

// UserClient mock implementation of user.Client
type UserClient struct {
	CreateUserFunc func(request *user.CreateUserRequest) (*user.CreateUserResponse, error)
	GetUserFunc    func(request *user.GetUserRequest) (*user.GetUserResponse, error)
	GetUsersFunc   func(request *user.GetUsersRequest) (*user.GetUsersResponse, error)
	UpdateUserFunc func(request *user.UpdateUserRequest, jwt string) (*user.UpdateUserResponse, error)
	DeleteUserFunc func(request *user.DeleteUserRequest, jwt string) (*user.DeleteUserResponse, error)
	VerifyUserFunc func(request *user.VerifyUserRequest, jwt string) (*user.VerifyUserResponse, error)
}

func (m *UserClient) CreateUser(request *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	return m.CreateUserFunc(request)
}

func (m *UserClient) GetUser(request *user.GetUserRequest) (*user.GetUserResponse, error) {
	return m.GetUserFunc(request)
}

func (m *UserClient) GetUsers(request *user.GetUsersRequest) (*user.GetUsersResponse, error) {
	return m.GetUsersFunc(request)
}

func (m *UserClient) UpdateUser(request *user.UpdateUserRequest, jwt string) (*user.UpdateUserResponse, error) {
	return m.UpdateUserFunc(request, jwt)
}

func (m *UserClient) DeleteUser(request *user.DeleteUserRequest, jwt string) (*user.DeleteUserResponse, error) {
	return m.DeleteUserFunc(request, jwt)
}

func (m *UserClient) VerifyUser(request *user.VerifyUserRequest, jwt string) (*user.VerifyUserResponse, error) {
	return m.VerifyUserFunc(request, jwt)
}
