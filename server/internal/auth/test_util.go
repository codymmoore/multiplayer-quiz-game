package auth

import (
	api "common/api/auth"
	"context"
)

// mockService mock implementation of Service
type mockService struct {
	loginFunc                 func(context context.Context, request *api.LoginRequest) (*api.LoginResponse, error)
	logoutFunc                func(context context.Context, request *api.LogoutRequest) (*api.LogoutResponse, error)
	renewFunc                 func(context context.Context, request *api.RenewRequest) (*api.RenewResponse, error)
	sendVerificationEmailFunc func(
		context context.Context,
		request *api.SendVerificationEmailRequest,
	) (*api.SendVerificationEmailResponse, error)
	verifyEmailFunc func(context context.Context, request *api.VerifyEmailRequest) (
		*api.VerifyEmailResponse,
		error,
	)
}

func (m *mockService) Login(context context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {
	return m.loginFunc(context, request)
}

func (m *mockService) Logout(context context.Context, request *api.LogoutRequest) (*api.LogoutResponse, error) {
	return m.logoutFunc(context, request)
}

func (m *mockService) Renew(context context.Context, request *api.RenewRequest) (*api.RenewResponse, error) {
	return m.renewFunc(context, request)
}

func (m *mockService) SendVerificationEmail(
	context context.Context,
	request *api.SendVerificationEmailRequest,
) (*api.SendVerificationEmailResponse, error) {
	return m.sendVerificationEmailFunc(context, request)
}

func (m *mockService) VerifyEmail(context context.Context, request *api.VerifyEmailRequest) (
	*api.VerifyEmailResponse,
	error,
) {
	return m.verifyEmailFunc(context, request)
}
