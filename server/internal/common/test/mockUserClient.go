package test

import (
	api "common/api/user"
)

// MockUserClient mock implementation of user.Client
type MockUserClient struct {
	CreateUserFunc func(request *api.CreateUserRequest, jwt string) (*api.CreateUserResponse, error)
	GetUserFunc    func(request *api.GetUserRequest, jwt string) (*api.GetUserResponse, error)
	GetUsersFunc   func(request *api.GetUsersRequest, jwt string) (*api.GetUsersResponse, error)
	UpdateUserFunc func(request *api.UpdateUserRequest, jwt string) (*api.UpdateUserResponse, error)
	DeleteUserFunc func(request *api.DeleteUserRequest, jwt string) (*api.DeleteUserResponse, error)
	VerifyUserFunc func(request *api.VerifyUserRequest, jwt string) (*api.VerifyUserResponse, error)
}

func (m *MockUserClient) CreateUser(request *api.CreateUserRequest, jwt string) (*api.CreateUserResponse, error) {
	return m.CreateUserFunc(request, jwt)
}

func (m *MockUserClient) GetUser(request *api.GetUserRequest, jwt string) (*api.GetUserResponse, error) {
	return m.GetUserFunc(request, jwt)
}

func (m *MockUserClient) GetUsers(request *api.GetUsersRequest, jwt string) (*api.GetUsersResponse, error) {
	return m.GetUsersFunc(request, jwt)
}

func (m *MockUserClient) UpdateUser(request *api.UpdateUserRequest, jwt string) (*api.UpdateUserResponse, error) {
	return m.UpdateUserFunc(request, jwt)
}

func (m *MockUserClient) DeleteUser(request *api.DeleteUserRequest, jwt string) (*api.DeleteUserResponse, error) {
	return m.DeleteUserFunc(request, jwt)
}

func (m *MockUserClient) VerifyUser(request *api.VerifyUserRequest, jwt string) (*api.VerifyUserResponse, error) {
	return m.VerifyUserFunc(request, jwt)
}
