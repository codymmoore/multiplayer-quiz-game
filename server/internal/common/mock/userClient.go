package mock

import (
	api "common/api/user"
)

// UserClient mock implementation of user.Client
type UserClient struct {
	CreateUserFunc func(request *api.CreateUserRequest) (*api.CreateUserResponse, error)
	GetUserFunc    func(request *api.GetUserRequest) (*api.GetUserResponse, error)
	GetUsersFunc   func(request *api.GetUsersRequest) (*api.GetUsersResponse, error)
	UpdateUserFunc func(request *api.UpdateUserRequest, jwt string) (*api.UpdateUserResponse, error)
	DeleteUserFunc func(request *api.DeleteUserRequest, jwt string) (*api.DeleteUserResponse, error)
	VerifyUserFunc func(request *api.VerifyUserRequest, jwt string) (*api.VerifyUserResponse, error)
}

func (m *UserClient) CreateUser(request *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	return m.CreateUserFunc(request)
}

func (m *UserClient) GetUser(request *api.GetUserRequest) (*api.GetUserResponse, error) {
	return m.GetUserFunc(request)
}

func (m *UserClient) GetUsers(request *api.GetUsersRequest) (*api.GetUsersResponse, error) {
	return m.GetUsersFunc(request)
}

func (m *UserClient) UpdateUser(request *api.UpdateUserRequest, jwt string) (*api.UpdateUserResponse, error) {
	return m.UpdateUserFunc(request, jwt)
}

func (m *UserClient) DeleteUser(request *api.DeleteUserRequest, jwt string) (*api.DeleteUserResponse, error) {
	return m.DeleteUserFunc(request, jwt)
}

func (m *UserClient) VerifyUser(request *api.VerifyUserRequest, jwt string) (*api.VerifyUserResponse, error) {
	return m.VerifyUserFunc(request, jwt)
}
