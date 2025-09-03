package user

import (
	"common/test"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"
)

const (
	mockUrl   = "http://mockUrl"
	jwtString = "mockJwtString"
)

func TestClient_CreateUser_Success(t *testing.T) {
	mockResponse := &CreateUserResponse{
		UserId: 1,
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock create user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusCreated, string(responseBody)),
	}

	request := &CreateUserRequest{
		Username: test.ValidUsername,
		Email:    test.ValidEmail,
		Password: test.ValidPassword,
	}
	response, err := client.CreateUser(request, jwtString)
	if err != nil {
		t.Errorf(`client.CreateUser(request, jwtString) error = "%v", expected "<nil>"`, err)
	}
	if *response != *mockResponse {
		t.Errorf(`client.CreateUser(request, jwtString) response = "%v", expected "%v"`, *response, *mockResponse)
	}
}

func TestClient_CreateUser_NilRequest(t *testing.T) {
	mockResponse := &CreateUserResponse{
		UserId: 1,
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock create user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusCreated, string(responseBody)),
	}

	response, err := client.CreateUser(nil, jwtString)
	if err == nil {
		t.Error(`client.CreateUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.CreateUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_CreateUser_InvalidBaseUrl(t *testing.T) {
	mockResponse := &CreateUserResponse{
		UserId: 1,
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock create user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    ":",
		HttpClient: test.GetMockHttpClient(http.StatusCreated, string(responseBody)),
	}

	request := &CreateUserRequest{
		Username: test.ValidUsername,
		Email:    test.ValidEmail,
		Password: test.ValidPassword,
	}
	response, err := client.CreateUser(request, jwtString)
	if err == nil {
		t.Error(`client.CreateUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.CreateUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_CreateUser_HttpClientError(t *testing.T) {
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockErrorHttpClient(),
	}

	request := &CreateUserRequest{
		Username: test.ValidUsername,
		Email:    test.ValidEmail,
		Password: test.ValidPassword,
	}
	response, err := client.CreateUser(request, jwtString)
	if err == nil {
		t.Error(`client.CreateUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.CreateUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_CreateUser_BadStatusCode(t *testing.T) {
	mockResponse := &CreateUserResponse{
		UserId: 1,
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock create user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusBadRequest, string(responseBody)),
	}

	request := &CreateUserRequest{
		Username: test.ValidUsername,
		Email:    test.ValidEmail,
		Password: test.ValidPassword,
	}
	response, err := client.CreateUser(request, jwtString)
	if err == nil {
		t.Error(`client.CreateUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.CreateUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_CreateUser_BadResponseBody(t *testing.T) {
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusCreated, "badResponseBody"),
	}

	request := &CreateUserRequest{
		Username: test.ValidUsername,
		Email:    test.ValidEmail,
		Password: test.ValidPassword,
	}
	response, err := client.CreateUser(request, jwtString)
	if err == nil {
		t.Error(`client.CreateUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.CreateUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_GetUser_Success(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	mockResponse := &GetUserResponse{
		UserId:       userId,
		Username:     username,
		Email:        email,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    time.Now().Round(time.Nanosecond),
		UpdatedAt:    time.Now().Round(time.Nanosecond),
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusOK, string(responseBody)),
	}

	request := &GetUserRequest{
		UserId:   &userId,
		Username: &username,
		Email:    &email,
	}
	response, err := client.GetUser(request, jwtString)
	if err != nil {
		t.Errorf(`client.GetUser(request, jwtString) error = "%v", expected "<nil>"`, err)
	}
	if *response != *mockResponse {
		t.Errorf(`client.GetUser(request, jwtString) response = "%v", expected "%v"`, *response, *mockResponse)
	}
}

func TestClient_GetUser_NilRequest(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	mockResponse := &GetUserResponse{
		UserId:       userId,
		Username:     username,
		Email:        email,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    time.Now().Round(time.Nanosecond),
		UpdatedAt:    time.Now().Round(time.Nanosecond),
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusOK, string(responseBody)),
	}

	response, err := client.GetUser(nil, jwtString)
	if err == nil {
		t.Error(`client.GetUser(nil, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.GetUser(nil, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_GetUser_QueryStringFailure(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	mockResponse := &GetUserResponse{
		UserId:       userId,
		Username:     username,
		Email:        email,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    time.Now().Round(time.Nanosecond),
		UpdatedAt:    time.Now().Round(time.Nanosecond),
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusOK, string(responseBody)),
	}

	response, err := client.GetUser(&GetUserRequest{}, jwtString)
	if err == nil {
		t.Error(`client.GetUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.GetUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_GetUser_InvalidBaseUrl(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	mockResponse := &GetUserResponse{
		UserId:       userId,
		Username:     username,
		Email:        email,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    time.Now().Round(time.Nanosecond),
		UpdatedAt:    time.Now().Round(time.Nanosecond),
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    ":",
		HttpClient: test.GetMockHttpClient(http.StatusOK, string(responseBody)),
	}

	request := &GetUserRequest{
		UserId:   &userId,
		Username: &username,
		Email:    &email,
	}
	response, err := client.GetUser(request, jwtString)
	if err == nil {
		t.Error(`client.GetUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.GetUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_GetUser_HttpClientError(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockErrorHttpClient(),
	}

	request := &GetUserRequest{
		UserId:   &userId,
		Username: &username,
		Email:    &email,
	}
	response, err := client.GetUser(request, jwtString)
	if err == nil {
		t.Error(`client.GetUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.GetUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_GetUser_BadStatusCode(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	mockResponse := &GetUserResponse{
		UserId:       userId,
		Username:     username,
		Email:        email,
		PasswordHash: test.ValidPassword,
		IsVerified:   true,
		CreatedAt:    time.Now().Round(time.Nanosecond),
		UpdatedAt:    time.Now().Round(time.Nanosecond),
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusBadRequest, string(responseBody)),
	}

	request := &GetUserRequest{
		UserId:   &userId,
		Username: &username,
		Email:    &email,
	}
	response, err := client.GetUser(request, jwtString)
	if err == nil {
		t.Error(`client.GetUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.GetUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_GetUser_BadResponseBody(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusOK, "badResponseBody"),
	}

	request := &GetUserRequest{
		UserId:   &userId,
		Username: &username,
		Email:    &email,
	}
	response, err := client.GetUser(request, jwtString)
	if err == nil {
		t.Error(`client.GetUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.GetUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_GetUsers_Success(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	prevLink := mockUrl + "/users/prev"
	nextLink := mockUrl + "/users/next"
	mockResponse := &GetUsersResponse{
		Users: []User{
			{
				UserId:       userId,
				Username:     username,
				Email:        email,
				PasswordHash: test.ValidPassword,
				IsVerified:   true,
				CreatedAt:    time.Now().Round(time.Nanosecond),
				UpdatedAt:    time.Now().Round(time.Nanosecond),
			},
		},
		PrevLink: &prevLink,
		NextLink: &nextLink,
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusOK, string(responseBody)),
	}

	limit := 1
	offset := 2
	sortField := "CreatedAt"
	sortDirection := "ASC"
	request := &GetUsersRequest{
		Limit:         &limit,
		Offset:        &offset,
		SortField:     &sortField,
		SortDirection: &sortDirection,
	}
	response, err := client.GetUsers(request, jwtString)
	if err != nil {
		t.Errorf(`client.GetUsers(request, jwtString) error = "%v", expected "<nil>"`, err)
	}
	if response == nil {
		t.Error(`client.GetUsers(request, jwtString) response = "<nil>", expected non-nil`)
		return
	}
	if response.PrevLink != nil && *response.PrevLink != *mockResponse.PrevLink {
		t.Errorf(`response.PrevLink = "%v", expected "%v"`, response.PrevLink, mockResponse.PrevLink)
	}
	if response.NextLink != nil && *response.NextLink != *mockResponse.NextLink {
		t.Errorf(`response.PrevLink = "%v", expected "%v"`, response.NextLink, mockResponse.NextLink)
	}
	if !slices.Equal(response.Users, mockResponse.Users) {
		t.Errorf(`response.Users = "%v", expected "%v"`, response.Users, mockResponse.Users)
	}
}

func TestClient_GetUsers_NilRequest(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	prevLink := mockUrl + "/users/prev"
	nextLink := mockUrl + "/users/next"
	mockResponse := &GetUsersResponse{
		Users: []User{
			{
				UserId:       userId,
				Username:     username,
				Email:        email,
				PasswordHash: test.ValidPassword,
				IsVerified:   true,
				CreatedAt:    time.Now().Round(time.Nanosecond),
				UpdatedAt:    time.Now().Round(time.Nanosecond),
			},
		},
		PrevLink: &prevLink,
		NextLink: &nextLink,
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusOK, string(responseBody)),
	}

	response, err := client.GetUsers(nil, jwtString)
	if err == nil {
		t.Error(`client.GetUsers(nil, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.GetUsers(nil, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_GetUsers_InvalidBaseUrl(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	prevLink := mockUrl + "/users/prev"
	nextLink := mockUrl + "/users/next"
	mockResponse := &GetUsersResponse{
		Users: []User{
			{
				UserId:       userId,
				Username:     username,
				Email:        email,
				PasswordHash: test.ValidPassword,
				IsVerified:   true,
				CreatedAt:    time.Now().Round(time.Nanosecond),
				UpdatedAt:    time.Now().Round(time.Nanosecond),
			},
		},
		PrevLink: &prevLink,
		NextLink: &nextLink,
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    ":",
		HttpClient: test.GetMockHttpClient(http.StatusOK, string(responseBody)),
	}

	limit := 1
	offset := 2
	sortField := "CreatedAt"
	sortDirection := "ASC"
	request := &GetUsersRequest{
		Limit:         &limit,
		Offset:        &offset,
		SortField:     &sortField,
		SortDirection: &sortDirection,
	}
	response, err := client.GetUsers(request, jwtString)
	if err == nil {
		t.Error(`client.GetUsers(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.GetUsers(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_GetUsers_HttpClientError(t *testing.T) {
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockErrorHttpClient(),
	}

	limit := 1
	offset := 2
	sortField := "CreatedAt"
	sortDirection := "ASC"
	request := &GetUsersRequest{
		Limit:         &limit,
		Offset:        &offset,
		SortField:     &sortField,
		SortDirection: &sortDirection,
	}
	response, err := client.GetUsers(request, jwtString)
	if err == nil {
		t.Error(`client.GetUsers(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.GetUsers(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_GetUsers_BadStatusCode(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	prevLink := mockUrl + "/users/prev"
	nextLink := mockUrl + "/users/next"
	mockResponse := &GetUsersResponse{
		Users: []User{
			{
				UserId:       userId,
				Username:     username,
				Email:        email,
				PasswordHash: test.ValidPassword,
				IsVerified:   true,
				CreatedAt:    time.Now().Round(time.Nanosecond),
				UpdatedAt:    time.Now().Round(time.Nanosecond),
			},
		},
		PrevLink: &prevLink,
		NextLink: &nextLink,
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusBadRequest, string(responseBody)),
	}

	limit := 1
	offset := 2
	sortField := "CreatedAt"
	sortDirection := "ASC"
	request := &GetUsersRequest{
		Limit:         &limit,
		Offset:        &offset,
		SortField:     &sortField,
		SortDirection: &sortDirection,
	}
	response, err := client.GetUsers(request, jwtString)
	if err == nil {
		t.Error(`client.GetUsers(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.GetUsers(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_GetUsers_BadResponseBody(t *testing.T) {
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusOK, "badResponseBody"),
	}

	limit := 1
	offset := 2
	sortField := "CreatedAt"
	sortDirection := "ASC"
	request := &GetUsersRequest{
		Limit:         &limit,
		Offset:        &offset,
		SortField:     &sortField,
		SortDirection: &sortDirection,
	}
	response, err := client.GetUsers(request, jwtString)
	if err == nil {
		t.Error(`client.GetUsers(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.GetUsers(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_UpdateUser_Success(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword
	mockResponse := &UpdateUserResponse{
		UserId:       userId,
		Username:     username,
		Email:        email,
		PasswordHash: password,
		IsVerified:   true,
		CreatedAt:    time.Now().Round(time.Nanosecond),
		UpdatedAt:    time.Now().Round(time.Nanosecond),
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusOK, string(responseBody)),
	}

	request := &UpdateUserRequest{
		UserId:   userId,
		Username: &username,
		Email:    &email,
		Password: &password,
	}
	response, err := client.UpdateUser(request, jwtString)
	if err != nil {
		t.Errorf(`client.UpdateUser(request, jwtString) error = "%v", expected "<nil>"`, err)
	}
	if *response != *mockResponse {
		t.Errorf(`client.UpdateUser(request, jwtString) response = "%v", expected "%v"`, *response, *mockResponse)
	}
}

func TestClient_UpdateUser_NilRequest(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword
	mockResponse := &UpdateUserResponse{
		UserId:       userId,
		Username:     username,
		Email:        email,
		PasswordHash: password,
		IsVerified:   true,
		CreatedAt:    time.Now().Round(time.Nanosecond),
		UpdatedAt:    time.Now().Round(time.Nanosecond),
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusOK, string(responseBody)),
	}

	response, err := client.UpdateUser(nil, jwtString)
	if err == nil {
		t.Error(`client.UpdateUser(nil, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.UpdateUser(nil, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_UpdateUser_InvalidBaseUrl(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword
	mockResponse := &UpdateUserResponse{
		UserId:       userId,
		Username:     username,
		Email:        email,
		PasswordHash: password,
		IsVerified:   true,
		CreatedAt:    time.Now().Round(time.Nanosecond),
		UpdatedAt:    time.Now().Round(time.Nanosecond),
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    ":",
		HttpClient: test.GetMockHttpClient(http.StatusOK, string(responseBody)),
	}

	request := &UpdateUserRequest{
		UserId:   userId,
		Username: &username,
		Email:    &email,
		Password: &password,
	}
	response, err := client.UpdateUser(request, jwtString)
	if err == nil {
		t.Error(`client.UpdateUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.UpdateUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_UpdateUser_HttpClientError(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockErrorHttpClient(),
	}

	request := &UpdateUserRequest{
		UserId:   userId,
		Username: &username,
		Email:    &email,
		Password: &password,
	}
	response, err := client.UpdateUser(request, jwtString)
	if err == nil {
		t.Error(`client.UpdateUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.UpdateUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_UpdateUser_BadStatusCode(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword
	mockResponse := &UpdateUserResponse{
		UserId:       userId,
		Username:     username,
		Email:        email,
		PasswordHash: password,
		IsVerified:   true,
		CreatedAt:    time.Now().Round(time.Nanosecond),
		UpdatedAt:    time.Now().Round(time.Nanosecond),
	}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusBadRequest, string(responseBody)),
	}

	request := &UpdateUserRequest{
		UserId:   userId,
		Username: &username,
		Email:    &email,
		Password: &password,
	}
	response, err := client.UpdateUser(request, jwtString)
	if err == nil {
		t.Error(`client.UpdateUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.UpdateUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_UpdateUser_BadResponseBody(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusOK, "badResponseBody"),
	}

	request := &UpdateUserRequest{
		UserId:   userId,
		Username: &username,
		Email:    &email,
		Password: &password,
	}
	response, err := client.UpdateUser(request, jwtString)
	if err == nil {
		t.Error(`client.UpdateUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.UpdateUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_DeleteUser_Success(t *testing.T) {
	userId := 1
	mockResponse := &DeleteUserResponse{}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusNoContent, string(responseBody)),
	}

	request := &DeleteUserRequest{
		UserId: userId,
	}
	response, err := client.DeleteUser(request, jwtString)
	if err != nil {
		t.Errorf(`client.DeleteUser(request, jwtString) error = "%v", expected "<nil>"`, err)
	}
	if *response != *mockResponse {
		t.Errorf(`client.DeleteUser(request, jwtString) response = "%v", expected "%v"`, *response, *mockResponse)
	}
}

func TestClient_DeleteUser_NilRequest(t *testing.T) {
	mockResponse := &DeleteUserResponse{}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusNoContent, string(responseBody)),
	}

	response, err := client.DeleteUser(nil, jwtString)
	if err == nil {
		t.Error(`client.DeleteUser(nil, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.DeleteUser(nil, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_DeleteUser_InvalidBaseUrl(t *testing.T) {
	userId := 1
	mockResponse := &DeleteUserResponse{}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    ":",
		HttpClient: test.GetMockHttpClient(http.StatusNoContent, string(responseBody)),
	}

	request := &DeleteUserRequest{
		UserId: userId,
	}
	response, err := client.DeleteUser(request, jwtString)
	if err == nil {
		t.Error(`client.DeleteUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.DeleteUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_DeleteUser_HttpClientError(t *testing.T) {
	userId := 1
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockErrorHttpClient(),
	}

	request := &DeleteUserRequest{
		UserId: userId,
	}
	response, err := client.DeleteUser(request, jwtString)
	if err == nil {
		t.Error(`client.DeleteUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.DeleteUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_DeleteUser_BadStatusCode(t *testing.T) {
	userId := 1
	mockResponse := &DeleteUserResponse{}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusBadRequest, string(responseBody)),
	}

	request := &DeleteUserRequest{
		UserId: userId,
	}
	response, err := client.DeleteUser(request, jwtString)
	if err == nil {
		t.Error(`client.DeleteUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.DeleteUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_VerifyUser_Success(t *testing.T) {
	userId := 1
	mockResponse := &VerifyUserResponse{}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusNoContent, string(responseBody)),
	}

	request := &VerifyUserRequest{
		UserId: userId,
	}
	response, err := client.VerifyUser(request, jwtString)
	if err != nil {
		t.Errorf(`client.VerifyUser(request, jwtString) error = "%v", expected "<nil>"`, err)
	}
	if *response != *mockResponse {
		t.Errorf(`client.VerifyUser(request, jwtString) response = "%v", expected "%v"`, *response, *mockResponse)
	}
}

func TestClient_VerifyUser_NilRequest(t *testing.T) {
	mockResponse := &VerifyUserResponse{}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusNoContent, string(responseBody)),
	}

	response, err := client.VerifyUser(nil, jwtString)
	if err == nil {
		t.Error(`client.VerifyUser(nil, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.VerifyUser(nil, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_VerifyUser_InvalidBaseUrl(t *testing.T) {
	userId := 1
	mockResponse := &VerifyUserResponse{}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    ":",
		HttpClient: test.GetMockHttpClient(http.StatusNoContent, string(responseBody)),
	}

	request := &VerifyUserRequest{
		UserId: userId,
	}
	response, err := client.VerifyUser(request, jwtString)
	if err == nil {
		t.Error(`client.VerifyUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.VerifyUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_VerifyUser_HttpClientError(t *testing.T) {
	userId := 1
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockErrorHttpClient(),
	}

	request := &VerifyUserRequest{
		UserId: userId,
	}
	response, err := client.VerifyUser(request, jwtString)
	if err == nil {
		t.Error(`client.VerifyUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.VerifyUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_VerifyUser_BadStatusCode(t *testing.T) {
	userId := 1
	mockResponse := &VerifyUserResponse{}
	responseBody, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("unable to marshal mock get user response: %v", err)
	}
	client := ClientImpl{
		BaseUrl:    mockUrl,
		HttpClient: test.GetMockHttpClient(http.StatusBadRequest, string(responseBody)),
	}

	request := &VerifyUserRequest{
		UserId: userId,
	}
	response, err := client.VerifyUser(request, jwtString)
	if err == nil {
		t.Error(`client.VerifyUser(request, jwtString) error = "<nil>", expected non-nil`)
	}
	if response != nil {
		t.Errorf(`client.VerifyUser(request, jwtString) response = "%v", expected "<nil>"`, response)
	}
}

func TestClient_createGetUserQueryString(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	expectedQueryString := fmt.Sprintf(
		"%s=%s&%s=%d&%s=%s",
		EmailKey,
		strings.ReplaceAll(email, "@", "%40"),
		UserIdKey,
		userId,
		UsernameKey,
		username,
	)

	request := &GetUserRequest{
		UserId:   &userId,
		Username: &username,
		Email:    &email,
	}
	queryString, err := createGetUserQueryString(request)
	if err != nil {
		t.Errorf("failed to create query string: %v", err)
	}
	if queryString != expectedQueryString {
		t.Errorf(`queryString: "%s", expected "%s"`, queryString, expectedQueryString)
	}
}

func TestClient_createGetUsersQueryString(t *testing.T) {
	limit := 1
	offset := 2
	sortField := "CreatedAt"
	sortDirection := "asc"
	expectedQueryString := fmt.Sprintf(
		"%s=%d&%s=%d&%s=%s&%s=%s",
		LimitKey,
		limit,
		OffsetKey,
		offset,
		SortDirectionKey,
		sortDirection,
		SortFieldKey,
		sortField,
	)

	request := &GetUsersRequest{
		Limit:         &limit,
		Offset:        &offset,
		SortField:     &sortField,
		SortDirection: &sortDirection,
	}
	queryString, err := createGetUsersQueryString(request)
	if err != nil {
		t.Errorf("failed to create query string: %v", err)
	}
	if queryString != expectedQueryString {
		t.Errorf(`queryString: "%s", expected "%s"`, queryString, expectedQueryString)
	}
}
