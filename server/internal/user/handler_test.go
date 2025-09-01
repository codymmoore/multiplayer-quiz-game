package user

import (
	"common"
	api "common/api/user"
	"common/test"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateUserHandler_Success(t *testing.T) {
	userId := 1
	service := &mockService{
		createUserFunc: func(context context.Context, request *api.CreateUserRequest) (
			*api.CreateUserResponse,
			error,
		) {
			return &api.CreateUserResponse{UserId: userId}, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return nil, nil
		},
	}

	payload := fmt.Sprintf(
		`{"username": "%s", "email": "%s", "password": "%s"}`,
		test.ValidUsername,
		test.ValidEmail,
		test.ValidPassword,
	)
	request := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/user", CreateUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusCreated)
	}

	var response api.CreateUserResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Errorf(`json.NewDecoder(recorder.Body).Decode(&response) = "%v", expected "<nil>"`, err)
	}

	if response.UserId != userId {
		t.Errorf(`response.UserId = "%d", expected "1"`, response.UserId)
	}
}

func TestCreateUserHandler_InvalidRequestBody(t *testing.T) {
	service := &mockService{
		createUserFunc: func(context context.Context, request *api.CreateUserRequest) (
			*api.CreateUserResponse,
			error,
		) {
			return &api.CreateUserResponse{UserId: 1}, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return nil, nil
		},
	}

	payload := `{"invalidField": "uh"}`
	request := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/user", CreateUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestCreateUserHandler_InvalidRequestObject(t *testing.T) {
	service := &mockService{
		createUserFunc: func(context context.Context, request *api.CreateUserRequest) (
			*api.CreateUserResponse,
			error,
		) {
			return &api.CreateUserResponse{UserId: 1}, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &api.GetUserResponse{}, nil
		},
	}

	payload := fmt.Sprintf(
		`{"username": "%s", "email": "%s", "password": "%s"}`,
		test.ValidUsername,
		test.ValidEmail,
		test.ValidPassword,
	)
	request := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/user", CreateUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestCreateUserHandler_ServiceFailure(t *testing.T) {
	service := &mockService{
		createUserFunc: func(context context.Context, request *api.CreateUserRequest) (
			*api.CreateUserResponse,
			error,
		) {
			return nil, errors.New("")
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return nil, nil
		},
	}

	payload := fmt.Sprintf(
		`{"username": "%s", "email": "%s", "password": "%s"}`,
		test.ValidUsername,
		test.ValidEmail,
		test.ValidPassword,
	)
	request := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/user", CreateUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}

func TestGetCurrentUserHandler_Success(t *testing.T) {
	currentTime := time.Now()
	mockResponse := api.GetUserResponse{
		UserId:       1,
		Username:     test.ValidUsername,
		Email:        test.ValidEmail,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	service := &mockService{
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &mockResponse, nil
		},
	}

	userClaims := &common.UserClaims{
		ID:       1,
		Username: test.ValidUsername,
		Email:    test.ValidEmail,
	}
	ctx := context.WithValue(context.Background(), common.UserClaimsCtxKey, userClaims)
	request := httptest.NewRequestWithContext(ctx, http.MethodGet, "/user/me", strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user/me", GetCurrentUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusOK)
	}

	var response api.GetUserResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Errorf(`json.NewDecoder(recorder.Body).Decode(&response) = "%v", expected "<nil>"`, err)
	}
	assertUserEqual(t, &response, &mockResponse)
}

func TestGetCurrentUserHandler_MissingUserClaims(t *testing.T) {
	currentTime := time.Now()
	mockResponse := api.GetUserResponse{
		UserId:       1,
		Username:     test.ValidUsername,
		Email:        test.ValidEmail,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	service := &mockService{
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &mockResponse, nil
		},
	}

	request := httptest.NewRequest(http.MethodGet, "/user/me", strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user/me", GetCurrentUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusUnauthorized)
	}
}

func TestGetCurrentUserHandler_ServiceFailure(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return nil, errors.New("")
		},
	}

	userClaims := &common.UserClaims{
		ID:       1,
		Username: test.ValidUsername,
		Email:    test.ValidEmail,
	}
	ctx := context.WithValue(context.Background(), common.UserClaimsCtxKey, userClaims)
	request := httptest.NewRequestWithContext(ctx, http.MethodGet, "/user/me", strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user/me", GetCurrentUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}

func TestGetCurrentUserHandler_UserNotFound(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return nil, nil
		},
	}

	userClaims := &common.UserClaims{
		ID:       1,
		Username: test.ValidUsername,
		Email:    test.ValidEmail,
	}
	ctx := context.WithValue(context.Background(), common.UserClaimsCtxKey, userClaims)
	request := httptest.NewRequestWithContext(ctx, http.MethodGet, "/user/me", strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user/me", GetCurrentUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestGetUserHandler_Success(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	currentTime := time.Now()
	mockResponse := api.GetUserResponse{
		UserId:       userId,
		Username:     username,
		Email:        email,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	service := &mockService{
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &mockResponse, nil
		},
	}

	urlString := fmt.Sprintf("/user?id=%d&username=%s&email=%s", userId, username, email)
	request := httptest.NewRequest(http.MethodGet, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user", GetUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusOK)
	}

	var response api.GetUserResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Errorf(`json.NewDecoder(recorder.Body).Decode(&response) = "%v", expected "<nil>"`, err)
	}
	assertUserEqual(t, &response, &mockResponse)
}

func TestGetUserHandler_RequestGenerationError_InvalidUserId(t *testing.T) {
	username := test.ValidUsername
	email := test.ValidEmail
	currentTime := time.Now()
	mockResponse := api.GetUserResponse{
		UserId:       1,
		Username:     username,
		Email:        email,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	service := &mockService{
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &mockResponse, nil
		},
	}

	urlString := fmt.Sprintf("/user?id=%s&username=%s&email=%s", "invalidId", username, email)
	request := httptest.NewRequest(http.MethodGet, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user", GetUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestGetUserHandler_InvalidRequest(t *testing.T) {
	userId := -1
	username := test.ValidUsername
	email := test.ValidEmail
	currentTime := time.Now()
	mockResponse := api.GetUserResponse{
		UserId:       userId,
		Username:     username,
		Email:        email,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	service := &mockService{
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &mockResponse, nil
		},
	}

	urlString := fmt.Sprintf("/user?id=%d&username=%s&email=%s", userId, username, email)
	request := httptest.NewRequest(http.MethodGet, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user", GetUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestGetUserHandler_ServiceFailure(t *testing.T) {
	service := &mockService{
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return nil, errors.New("")
		},
	}

	urlString := fmt.Sprintf("/user?id=%d&username=%s&email=%s", 1, test.ValidUsername, test.ValidEmail)
	request := httptest.NewRequest(http.MethodGet, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user", GetUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}

func TestGetUsersHandler_Success(t *testing.T) {
	currentTime := time.Now()
	mockUser := api.User{
		UserId:       1,
		Username:     test.ValidUsername,
		Email:        test.ValidEmail,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	prevLink := "prevLink"
	nextLink := "nextLink"
	mockResponse := api.GetUsersResponse{
		Users: []api.User{
			mockUser,
		},
		PrevLink: &prevLink,
		NextLink: &nextLink,
	}
	service := &mockService{
		getUsersFunc: func(context context.Context, request *api.GetUsersRequest) (
			*api.GetUsersResponse,
			error,
		) {
			return &mockResponse, nil
		},
	}

	urlString := fmt.Sprintf("/user/all?limit=%d&offset=%d&sortField=%s&sortDirection=%s", 10, 20, "CreatedAt", "desc")
	request := httptest.NewRequest(http.MethodGet, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user/all", GetUsersHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusOK)
	}

	var response api.GetUsersResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Errorf(`json.NewDecoder(recorder.Body).Decode(&response) = "%v", expected "<nil>"`, err)
	}
	assertUserEqual(t, &response.Users[0], &mockUser)

	if *response.PrevLink != *mockResponse.PrevLink {
		t.Errorf(`response.PrevLink = "%s", expected "%s"`, *response.PrevLink, *mockResponse.PrevLink)
	}
	if *response.NextLink != *mockResponse.NextLink {
		t.Errorf(`response.NextLink = "%s", expected "%s"`, *response.NextLink, *mockResponse.NextLink)
	}
}

func TestGetUsersHandler_RequestGenerationError_InvalidLimit(t *testing.T) {
	currentTime := time.Now()
	mockUser := api.User{
		UserId:       1,
		Username:     test.ValidUsername,
		Email:        test.ValidEmail,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	mockResponse := api.GetUsersResponse{
		Users: []api.User{
			mockUser,
		},
	}
	service := &mockService{
		getUsersFunc: func(context context.Context, request *api.GetUsersRequest) (
			*api.GetUsersResponse,
			error,
		) {
			return &mockResponse, nil
		},
	}

	urlString := fmt.Sprintf(
		"/user/all?limit=%s&offset=%d&sortField=%s&sortDirection=%s",
		"invalidLimit",
		20,
		"CreatedAt",
		"desc",
	)
	request := httptest.NewRequest(http.MethodGet, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user/all", GetUsersHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestGetUsersHandler_RequestGenerationError_InvalidOffset(t *testing.T) {
	currentTime := time.Now()
	mockUser := api.User{
		UserId:       1,
		Username:     test.ValidUsername,
		Email:        test.ValidEmail,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	mockResponse := api.GetUsersResponse{
		Users: []api.User{
			mockUser,
		},
	}
	service := &mockService{
		getUsersFunc: func(context context.Context, request *api.GetUsersRequest) (
			*api.GetUsersResponse,
			error,
		) {
			return &mockResponse, nil
		},
	}

	urlString := fmt.Sprintf(
		"/user/all?limit=%d&offset=%s&sortField=%s&sortDirection=%s",
		10,
		"invalidOffset",
		"CreatedAt",
		"desc",
	)
	request := httptest.NewRequest(http.MethodGet, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user/all", GetUsersHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestGetUsersHandler_InvalidRequest(t *testing.T) {
	currentTime := time.Now()
	mockUser := api.User{
		UserId:       1,
		Username:     test.ValidUsername,
		Email:        test.ValidEmail,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	mockResponse := api.GetUsersResponse{
		Users: []api.User{
			mockUser,
		},
	}
	service := &mockService{
		getUsersFunc: func(context context.Context, request *api.GetUsersRequest) (
			*api.GetUsersResponse,
			error,
		) {
			return &mockResponse, nil
		},
	}

	urlString := fmt.Sprintf("/user/all?limit=%d&offset=%d&sortField=%s&sortDirection=%s", -1, 20, "CreatedAt", "desc")
	request := httptest.NewRequest(http.MethodGet, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user/all", GetUsersHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestGetUsersHandler_ServiceFailure(t *testing.T) {
	service := &mockService{
		getUsersFunc: func(context context.Context, request *api.GetUsersRequest) (
			*api.GetUsersResponse,
			error,
		) {
			return nil, errors.New("")
		},
	}

	urlString := fmt.Sprintf("/user/all?limit=%d&offset=%d&sortField=%s&sortDirection=%s", 1, 20, "CreatedAt", "desc")
	request := httptest.NewRequest(http.MethodGet, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/user/all", GetUsersHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}

func TestUpdateUserHandler_Success(t *testing.T) {
	userId := 1
	currentTime := time.Now()
	mockUser := api.UpdateUserResponse{
		UserId:       userId,
		Username:     test.ValidUsername,
		Email:        test.ValidEmail,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}

	service := &mockService{
		updateUserFunc: func(context context.Context, request *api.UpdateUserRequest) (
			*api.UpdateUserResponse,
			error,
		) {
			return &mockUser, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			if request.UserId != nil {
				return &api.GetUserResponse{}, nil
			} else {
				return nil, nil
			}
		},
	}

	payload := fmt.Sprintf(
		`{"username": "%s", "email": "%s", "password": "%s"}`,
		test.ValidUsername,
		test.ValidEmail,
		test.ValidPassword,
	)
	urlString := fmt.Sprintf("/user/%d", userId)
	request := httptest.NewRequest(http.MethodPatch, urlString, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Patch("/user/{id}", UpdateUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusOK)
	}

	var response api.UpdateUserResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Errorf(`json.NewDecoder(recorder.Body).Decode(&response) = "%v", expected "<nil>"`, err)
	}
	assertUserEqual(t, &response, &mockUser)
}

func TestUpdateUserHandler_RequestGenerationError_InvalidUserId(t *testing.T) {
	userId := 1
	currentTime := time.Now()
	mockUser := api.UpdateUserResponse{
		UserId:       userId,
		Username:     test.ValidUsername,
		Email:        test.ValidEmail,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}

	service := &mockService{
		updateUserFunc: func(context context.Context, request *api.UpdateUserRequest) (
			*api.UpdateUserResponse,
			error,
		) {
			return &mockUser, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			if request.UserId != nil {
				return &api.GetUserResponse{}, nil
			} else {
				return nil, nil
			}
		},
	}

	payload := fmt.Sprintf(
		`{"username": "%s", "email": "%s", "password": "%s"}`,
		test.ValidUsername,
		test.ValidEmail,
		test.ValidPassword,
	)
	urlString := fmt.Sprintf("/user/%s", "invalidId")
	request := httptest.NewRequest(http.MethodPatch, urlString, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Patch("/user/{id}", UpdateUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestUpdateUserHandler_RequestGenerationError_InvalidRequestBody(t *testing.T) {
	userId := 1
	currentTime := time.Now()
	mockUser := api.UpdateUserResponse{
		UserId:       userId,
		Username:     test.ValidUsername,
		Email:        test.ValidEmail,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}

	service := &mockService{
		updateUserFunc: func(context context.Context, request *api.UpdateUserRequest) (
			*api.UpdateUserResponse,
			error,
		) {
			return &mockUser, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			if request.UserId != nil {
				return &api.GetUserResponse{}, nil
			} else {
				return nil, nil
			}
		},
	}

	payload := `"invalidRequest"`
	urlString := fmt.Sprintf("/user/%d", userId)
	request := httptest.NewRequest(http.MethodPatch, urlString, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Patch("/user/{id}", UpdateUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestUpdateUserHandler_InvalidRequest(t *testing.T) {
	userId := -1
	currentTime := time.Now()
	mockUser := api.UpdateUserResponse{
		UserId:       userId,
		Username:     test.ValidUsername,
		Email:        test.ValidEmail,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}

	service := &mockService{
		updateUserFunc: func(context context.Context, request *api.UpdateUserRequest) (
			*api.UpdateUserResponse,
			error,
		) {
			return &mockUser, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			if request.UserId != nil {
				return &api.GetUserResponse{}, nil
			} else {
				return nil, nil
			}
		},
	}

	payload := fmt.Sprintf(
		`{"username": "%s", "email": "%s", "password": "%s"}`,
		test.ValidUsername,
		test.ValidEmail,
		test.ValidPassword,
	)
	urlString := fmt.Sprintf("/user/%d", userId)
	request := httptest.NewRequest(http.MethodPatch, urlString, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Patch("/user/{id}", UpdateUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestUpdateUserHandler_ServiceFailure(t *testing.T) {
	userId := 1
	service := &mockService{
		updateUserFunc: func(context context.Context, request *api.UpdateUserRequest) (
			*api.UpdateUserResponse,
			error,
		) {
			return nil, errors.New("")
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			if request.UserId != nil {
				return &api.GetUserResponse{}, nil
			} else {
				return nil, nil
			}
		},
	}

	payload := fmt.Sprintf(
		`{"username": "%s", "email": "%s", "password": "%s"}`,
		test.ValidUsername,
		test.ValidEmail,
		test.ValidPassword,
	)
	urlString := fmt.Sprintf("/user/%d", userId)
	request := httptest.NewRequest(http.MethodPatch, urlString, strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Patch("/user/{id}", UpdateUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}

func TestDeleteUserHandler_Success(t *testing.T) {
	userId := 1
	service := &mockService{
		deleteUserFunc: func(context context.Context, request *api.DeleteUserRequest) (
			*api.DeleteUserResponse,
			error,
		) {
			return &api.DeleteUserResponse{}, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &api.GetUserResponse{}, nil
		},
	}

	urlString := fmt.Sprintf("/user/%d", userId)
	request := httptest.NewRequest(http.MethodDelete, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/user/{id}", DeleteUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusNoContent)
	}
}

func TestDeleteUserHandler_RequestGenerationError_InvalidUserId(t *testing.T) {
	service := &mockService{
		deleteUserFunc: func(context context.Context, request *api.DeleteUserRequest) (
			*api.DeleteUserResponse,
			error,
		) {
			return &api.DeleteUserResponse{}, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &api.GetUserResponse{}, nil
		},
	}

	urlString := fmt.Sprintf("/user/%s", "invalidId")
	request := httptest.NewRequest(http.MethodDelete, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/user/{id}", DeleteUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestDeleteUserHandler_InvalidRequest(t *testing.T) {
	userId := -1
	service := &mockService{
		deleteUserFunc: func(context context.Context, request *api.DeleteUserRequest) (
			*api.DeleteUserResponse,
			error,
		) {
			return &api.DeleteUserResponse{}, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &api.GetUserResponse{}, nil
		},
	}

	urlString := fmt.Sprintf("/user/%d", userId)
	request := httptest.NewRequest(http.MethodDelete, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/user/{id}", DeleteUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestDeleteUserHandler_ServiceFailure(t *testing.T) {
	userId := 1
	service := &mockService{
		deleteUserFunc: func(context context.Context, request *api.DeleteUserRequest) (
			*api.DeleteUserResponse,
			error,
		) {
			return nil, errors.New("")
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &api.GetUserResponse{}, nil
		},
	}

	urlString := fmt.Sprintf("/user/%d", userId)
	request := httptest.NewRequest(http.MethodDelete, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/user/{id}", DeleteUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}

func TestVerifyUserHandler_Success(t *testing.T) {
	userId := 1
	service := &mockService{
		verifyUserFunc: func(context context.Context, request *api.VerifyUserRequest) (
			*api.VerifyUserResponse,
			error,
		) {
			return &api.VerifyUserResponse{}, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &api.GetUserResponse{}, nil
		},
	}

	urlString := fmt.Sprintf(api.VerifyUserEndpoint, userId)
	request := httptest.NewRequest(http.MethodPatch, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Patch(fmt.Sprintf(api.VerifyUserEndpoint, "{id}"), VerifyUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusNoContent)
	}
}

func TestVerifyUserHandler_RequestGenerationError_InvalidUserId(t *testing.T) {
	service := &mockService{
		verifyUserFunc: func(context context.Context, request *api.VerifyUserRequest) (
			*api.VerifyUserResponse,
			error,
		) {
			return &api.VerifyUserResponse{}, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &api.GetUserResponse{}, nil
		},
	}

	urlString := fmt.Sprintf(api.VerifyUserEndpoint, "invalidId")
	request := httptest.NewRequest(http.MethodPatch, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Patch(fmt.Sprintf(api.VerifyUserEndpoint, "{id}"), VerifyUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestVerifyUserHandler_InvalidRequest(t *testing.T) {
	userId := -1
	service := &mockService{
		verifyUserFunc: func(context context.Context, request *api.VerifyUserRequest) (
			*api.VerifyUserResponse,
			error,
		) {
			return &api.VerifyUserResponse{}, nil
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &api.GetUserResponse{}, nil
		},
	}

	urlString := fmt.Sprintf(api.VerifyUserEndpoint, userId)
	request := httptest.NewRequest(http.MethodPatch, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Patch(fmt.Sprintf(api.VerifyUserEndpoint, "{id}"), VerifyUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusBadRequest)
	}
}

func TestVerifyUserHandler_ServiceFailure(t *testing.T) {
	userId := 1
	service := &mockService{
		verifyUserFunc: func(context context.Context, request *api.VerifyUserRequest) (
			*api.VerifyUserResponse,
			error,
		) {
			return nil, errors.New("")
		},
		getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
			return &api.GetUserResponse{}, nil
		},
	}

	urlString := fmt.Sprintf(api.VerifyUserEndpoint, userId)
	request := httptest.NewRequest(http.MethodPatch, urlString, strings.NewReader(""))
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Patch(fmt.Sprintf(api.VerifyUserEndpoint, "{id}"), VerifyUserHandler(service))
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf(`recorder.Code = "%v", expected "%v"`, recorder.Code, http.StatusInternalServerError)
	}
}

func TestGenerateGetUsersRequest_Success(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail

	urlString := fmt.Sprintf("/user?id=%d&username=%s&email=%s", userId, username, email)
	request := httptest.NewRequest(http.MethodGet, urlString, strings.NewReader(""))

	getUserRequest, err := generateGetUserRequest(request)
	if err != nil {
		t.Errorf(`generateGetUserRequest(request) return error = "%v", expected "<nil>"`, err)
	}

	if getUserRequest == nil {
		t.Errorf(`generateGetUserRequest(request) = "<nil>", expected non-nil`)
		return
	}

	if *getUserRequest.UserId != userId {
		t.Errorf(`getUserRequest.UserId = "%d", expected "%d"`, *getUserRequest.UserId, userId)
	}

	if *getUserRequest.Username != username {
		t.Errorf(`getUserRequest.Username = "%s", expected "%s"`, *getUserRequest.Username, username)
	}

	if *getUserRequest.Email != email {
		t.Errorf(`*getUserRequest.Email = "%s", expected "%s"`, *getUserRequest.Email, email)
	}
}

func assertUserEqual(t *testing.T, actual *api.User, expected *api.User) {
	if actual.UserId != expected.UserId {
		t.Errorf(`actual.UserId = "%d", expected "%d"`, actual.UserId, expected.UserId)
	}

	if actual.Username != expected.Username {
		t.Errorf(`actual.Username = "%s", expected "%s"`, actual.Username, expected.Username)
	}

	if actual.Email != expected.Email {
		t.Errorf(`actual.Email = "%s", expected "%s"`, actual.Email, expected.Email)
	}

	if actual.PasswordHash != expected.PasswordHash {
		t.Errorf(`actual.PasswordHash = "%s", expected "%s"`, actual.PasswordHash, expected.PasswordHash)
	}

	if actual.IsVerified != true {
		t.Errorf(`actual.IsVerified = "%t"`, actual.IsVerified)
	}

	if !actual.CreatedAt.Equal(expected.CreatedAt) {
		t.Errorf(`actual.CreatedAt = "%s", expected "%s"`, actual.CreatedAt, expected.CreatedAt)
	}

	if !actual.UpdatedAt.Equal(expected.UpdatedAt) {
		t.Errorf(`actual.UpdatedAt = "%s", expected "%s"`, actual.UpdatedAt, expected.UpdatedAt)
	}
}
