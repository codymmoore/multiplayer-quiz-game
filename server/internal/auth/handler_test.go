package auth

import (
	api "common/auth"
	"common/test"
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLoginHandler_Success(t *testing.T) {
	username := test.ValidUsername
	password := test.ValidPassword
	expectedResponse := api.LoginResponse{
		AccessToken:  "token",
		RefreshToken: "refresh",
		IssuedAt:     time.Now().Round(0),
		ExpiresIn:    30 * time.Minute,
		User: api.User{
			UserId:   1,
			Username: username,
			Email:    test.ValidEmail,
		},
	}
	service := &mockService{
		loginFunc: func(context context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {
			if request.Username != username {
				t.Errorf(`request.Username = "%v", expected "%v"`, request.Username, username)
			}
			if request.Password != password {
				t.Errorf(`request.Password = "%v", expected "%v"`, request.Password, password)
			}
			return &expectedResponse, nil
		},
	}

	payload := fmt.Sprintf(
		`{"username":"%s","password":"%s"}`,
		username,
		password,
	)
	request := httptest.NewRequest(http.MethodPost, api.LoginEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.LoginEndpoint, LoginHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusOK)
	}

	var response api.LoginResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Errorf(`json.NewDecoder(recorder.Body).Decode(&response) = "%v", expected "<nil>"`, err)
	}
	if response != expectedResponse {
		t.Errorf(`response = "%v", expected "%v"`, response, expectedResponse)
	}
}

func TestLoginHandler_InvalidRequestBody(t *testing.T) {
	service := &mockService{
		loginFunc: func(context context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {
			return &api.LoginResponse{}, nil
		},
	}

	payload := "invalidRequest"
	request := httptest.NewRequest(http.MethodPost, api.LoginEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.LoginEndpoint, LoginHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestLoginHandler_InvalidRequestObject(t *testing.T) {
	username := test.ValidUsername
	password := test.ValidPassword
	expectedResponse := api.LoginResponse{
		AccessToken:  "token",
		RefreshToken: "refresh",
		IssuedAt:     time.Now().Round(0),
		ExpiresIn:    30 * time.Minute,
		User: api.User{
			UserId:   1,
			Username: username,
			Email:    test.ValidEmail,
		},
	}
	service := &mockService{
		loginFunc: func(context context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {
			return &expectedResponse, nil
		},
	}

	payload := fmt.Sprintf(
		`{"username":"%s","password":"%s"}`,
		"!@#$%^&*()",
		password,
	)
	request := httptest.NewRequest(http.MethodPost, api.LoginEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.LoginEndpoint, LoginHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestLoginHandler_ServiceFailure(t *testing.T) {
	username := test.ValidUsername
	password := test.ValidPassword
	service := &mockService{
		loginFunc: func(context context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {
			return nil, stderrors.New("")
		},
	}

	payload := fmt.Sprintf(
		`{"username":"%s","password":"%s"}`,
		username,
		password,
	)
	request := httptest.NewRequest(http.MethodPost, api.LoginEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.LoginEndpoint, LoginHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}

func TestLogoutHandler_Success(t *testing.T) {
	refreshToken := "refreshToken"
	service := &mockService{
		logoutFunc: func(context context.Context, request *api.LogoutRequest) (*api.LogoutResponse, error) {
			if request.RefreshToken != refreshToken {
				t.Errorf(`request.RefreshToken = "%v", expected "%v"`, request.RefreshToken, refreshToken)
			}
			return &api.LogoutResponse{}, nil
		},
	}

	payload := fmt.Sprintf(`{"refreshToken":"%s"}`, refreshToken)
	request := httptest.NewRequest(http.MethodPost, api.LogoutEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.LogoutEndpoint, LogoutHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusNoContent)
	}
}

func TestLogoutHandler_InvalidRequestBody(t *testing.T) {
	service := &mockService{
		logoutFunc: func(context context.Context, request *api.LogoutRequest) (*api.LogoutResponse, error) {
			return &api.LogoutResponse{}, nil
		},
	}

	payload := fmt.Sprintf(`{"refreshToken":"%s"}`, "token")
	request := httptest.NewRequest(http.MethodPost, api.LogoutEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.LogoutEndpoint, LogoutHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusNoContent)
	}
}

func TestLogoutHandler_ServiceFailure(t *testing.T) {
	service := &mockService{
		logoutFunc: func(context context.Context, request *api.LogoutRequest) (*api.LogoutResponse, error) {
			return nil, stderrors.New("")
		},
	}

	payload := fmt.Sprintf(`{"refreshToken":"%s"}`, "token")
	request := httptest.NewRequest(http.MethodPost, api.LogoutEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.LogoutEndpoint, LogoutHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}

func TestRenewHandler_Success(t *testing.T) {
	refreshToken := "refreshToken"
	expectedResponse := api.LoginResponse{
		AccessToken:  "token",
		RefreshToken: "refresh",
		IssuedAt:     time.Now().Round(0),
		ExpiresIn:    30 * time.Minute,
		User: api.User{
			UserId:   1,
			Username: test.ValidUsername,
			Email:    test.ValidEmail,
		},
	}
	service := &mockService{
		renewFunc: func(context context.Context, request *api.RenewRequest) (*api.RenewResponse, error) {
			if request.RefreshToken != refreshToken {
				t.Errorf(`request.RefreshToken = "%v", expected "%v"`, request.RefreshToken, refreshToken)
			}
			return &expectedResponse, nil
		},
	}

	payload := fmt.Sprintf(`{"refreshToken":"%s"}`, refreshToken)
	request := httptest.NewRequest(http.MethodPost, api.RenewEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.RenewEndpoint, RenewHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusOK)
	}

	var response api.LoginResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Errorf(`json.NewDecoder(recorder.Body).Decode(&response) = "%v", expected "<nil>"`, err)
	}
	if response != expectedResponse {
		t.Errorf(`response = "%v", expected "%v"`, response, expectedResponse)
	}
}

func TestRenewHandler_InvalidRequestBody(t *testing.T) {
	expectedResponse := api.LoginResponse{
		AccessToken:  "token",
		RefreshToken: "refresh",
		IssuedAt:     time.Now().Round(0),
		ExpiresIn:    30 * time.Minute,
		User: api.User{
			UserId:   1,
			Username: test.ValidUsername,
			Email:    test.ValidEmail,
		},
	}
	service := &mockService{
		renewFunc: func(context context.Context, request *api.RenewRequest) (*api.RenewResponse, error) {
			return &expectedResponse, nil
		},
	}

	payload := `invalidRefresh`
	request := httptest.NewRequest(http.MethodPost, api.RenewEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.RenewEndpoint, RenewHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestRenewHandler_ServiceFailure(t *testing.T) {
	service := &mockService{
		renewFunc: func(context context.Context, request *api.RenewRequest) (*api.RenewResponse, error) {
			return nil, stderrors.New("")
		},
	}

	payload := fmt.Sprintf(`{"refreshToken":"%s"}`, "token")
	request := httptest.NewRequest(http.MethodPost, api.RenewEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.RenewEndpoint, RenewHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}

func TestSendVerificationEmailHandler_Success(t *testing.T) {
	email := test.ValidEmail
	service := &mockService{
		sendVerificationEmailFunc: func(
			context context.Context,
			request *api.SendVerificationEmailRequest,
		) (*api.SendVerificationEmailResponse, error) {
			if request.Email != email {
				t.Errorf(`request.Email = "%v", expected "%v"`, request.Email, email)
			}
			return &api.SendVerificationEmailResponse{}, nil
		},
	}

	payload := fmt.Sprintf(`{"email":"%s"}`, email)
	request := httptest.NewRequest(http.MethodPost, api.SendVerificationEmailEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.SendVerificationEmailEndpoint, SendVerificationEmailHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusNoContent)
	}
}

func TestSendVerificationEmailHandler_InvalidRequestBody(t *testing.T) {
	service := &mockService{
		sendVerificationEmailFunc: func(
			context context.Context,
			request *api.SendVerificationEmailRequest,
		) (*api.SendVerificationEmailResponse, error) {
			return &api.SendVerificationEmailResponse{}, nil
		},
	}

	payload := fmt.Sprintf(`invalidRequest`)
	request := httptest.NewRequest(http.MethodPost, api.SendVerificationEmailEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.SendVerificationEmailEndpoint, SendVerificationEmailHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestSendVerificationEmailHandler_ServiceFailure(t *testing.T) {
	service := &mockService{
		sendVerificationEmailFunc: func(
			context context.Context,
			request *api.SendVerificationEmailRequest,
		) (*api.SendVerificationEmailResponse, error) {
			return nil, stderrors.New("")
		},
	}

	payload := fmt.Sprintf(`{"email":"%s"}`, test.ValidEmail)
	request := httptest.NewRequest(http.MethodPost, api.SendVerificationEmailEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.SendVerificationEmailEndpoint, SendVerificationEmailHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}

func TestVerifyEmailHandler_Success(t *testing.T) {
	verificationCode := "verificationCode"
	userId := 1
	service := &mockService{
		verifyEmailFunc: func(context context.Context, request *api.VerifyEmailRequest) (
			*api.VerifyEmailResponse,
			error,
		) {
			if request.VerificationCode != verificationCode {
				t.Errorf(`request.VerificationCode = "%v", expected "%v"`, request.VerificationCode, verificationCode)
			}
			return &api.VerifyEmailResponse{UserId: userId}, nil
		},
	}

	payload := fmt.Sprintf(`{"verificationCode":"%s"}`, verificationCode)
	request := httptest.NewRequest(http.MethodPost, api.VerifyEmailEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.VerifyEmailEndpoint, VerifyEmailHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusOK)
	}

	var response api.VerifyEmailResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Errorf(`json.NewDecoder(recorder.Body).Decode(&response) = "%v", expected "<nil>"`, err)
	}
	if response.UserId != userId {
		t.Errorf(`response.UserId = "%v", expected "%v"`, response.UserId, userId)
	}
}

func TestVerifyEmailHandler_InvalidRequestBody(t *testing.T) {
	service := &mockService{
		verifyEmailFunc: func(context context.Context, request *api.VerifyEmailRequest) (
			*api.VerifyEmailResponse,
			error,
		) {
			return &api.VerifyEmailResponse{UserId: 1}, nil
		},
	}

	payload := `invalidRequest`
	request := httptest.NewRequest(http.MethodPost, api.VerifyEmailEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.VerifyEmailEndpoint, VerifyEmailHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestVerifyEmailHandler_ServiceFailure(t *testing.T) {
	verificationCode := "verificationCode"
	service := &mockService{
		verifyEmailFunc: func(context context.Context, request *api.VerifyEmailRequest) (
			*api.VerifyEmailResponse,
			error,
		) {
			return nil, stderrors.New("")
		},
	}

	payload := fmt.Sprintf(`{"verificationCode":"%s"}`, verificationCode)
	request := httptest.NewRequest(http.MethodPost, api.VerifyEmailEndpoint, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post(api.VerifyEmailEndpoint, VerifyEmailHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}
