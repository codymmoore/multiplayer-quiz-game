package user

import (
	"bytes"
	"common/errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Client Interface for user service client
type Client interface {
	CreateUser(request *CreateUserRequest, jwt string) (*CreateUserResponse, error)
	GetUser(request *GetUserRequest, jwt string) (*GetUserResponse, error)
	GetUsers(request *GetUsersRequest, jwt string) (*GetUsersResponse, error)
	UpdateUser(request *UpdateUserRequest, jwt string) (*UpdateUserResponse, error)
	DeleteUser(request *DeleteUserRequest, jwt string) (*DeleteUserResponse, error)
	VerifyUser(request *VerifyUserRequest, jwt string) (*VerifyUserResponse, error)
}

// ClientImpl Implementation for user service client
type ClientImpl struct {
	BaseUrl    string
	HttpClient *http.Client
}

// CreateUser Create a new user
func (client *ClientImpl) CreateUser(request *CreateUserRequest, jwt string) (*CreateUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CreateUserRequest: %w", err)
	}

	urlString := client.BaseUrl + "/user"
	httpRequest, err := http.NewRequest(http.MethodPost, urlString, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+jwt)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.HttpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusCreated {
		return nil, &errors.HTTP{
			StatusCode: httpResponse.StatusCode,
			Message:    getErrorString(&httpResponse.Body),
		}
	}

	var response CreateUserResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode CreateUserResponse: %w", err)
	}

	return &response, nil
}

// GetUser Retrieve a user by ID, username, and/or email
func (client *ClientImpl) GetUser(request *GetUserRequest, jwt string) (*GetUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	paramString, err := createGetUserQueryString(request)
	if err != nil {
		return nil, fmt.Errorf("failed to create parameter string: %w", err)
	}
	urlString := client.BaseUrl + "/user?" + paramString

	httpRequest, err := http.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+jwt)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.HttpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, &errors.HTTP{
			StatusCode: httpResponse.StatusCode,
			Message:    getErrorString(&httpResponse.Body),
		}
	}

	var response GetUserResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode GetUserResponse: %w", err)
	}

	return &response, nil
}

// GetUsers Retrieve all users (paginated)
func (client *ClientImpl) GetUsers(request *GetUsersRequest, jwt string) (*GetUsersResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	paramString, err := createGetUsersQueryString(request)
	if err != nil {
		return nil, fmt.Errorf("failed to create parameter string: %w", err)
	}

	urlString := client.BaseUrl + "/user/all"
	if paramString != "" {
		urlString += "?" + paramString
	}

	httpRequest, err := http.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+jwt)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.HttpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, &errors.HTTP{
			StatusCode: httpResponse.StatusCode,
			Message:    getErrorString(&httpResponse.Body),
		}
	}

	var response GetUsersResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode GetUsersResponse: %w", err)
	}

	return &response, nil
}

// UpdateUser Update a user
func (client *ClientImpl) UpdateUser(request *UpdateUserRequest, jwt string) (*UpdateUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal UpdateUserRequest: %w", err)
	}

	urlString := client.BaseUrl + "/user/" + strconv.Itoa(request.UserId)
	httpRequest, err := http.NewRequest(http.MethodPatch, urlString, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+jwt)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.HttpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, &errors.HTTP{
			StatusCode: httpResponse.StatusCode,
			Message:    getErrorString(&httpResponse.Body),
		}
	}

	var response UpdateUserResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode UpdateUserResponse: %w", err)
	}

	return &response, nil
}

// DeleteUser Delete a user
func (client *ClientImpl) DeleteUser(request *DeleteUserRequest, jwt string) (*DeleteUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	urlString := client.BaseUrl + "/user/" + strconv.Itoa(request.UserId)
	httpRequest, err := http.NewRequest(http.MethodDelete, urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+jwt)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.HttpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusNoContent {
		return nil, &errors.HTTP{
			StatusCode: httpResponse.StatusCode,
			Message:    getErrorString(&httpResponse.Body),
		}
	}

	return &DeleteUserResponse{}, nil
}

// VerifyUser verifies a user
func (client *ClientImpl) VerifyUser(request *VerifyUserRequest, jwt string) (*VerifyUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	urlString := client.BaseUrl + fmt.Sprintf(VerifyUserEndpoint, request.UserId)
	httpRequest, err := http.NewRequest(http.MethodPatch, urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+jwt)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.HttpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusNoContent {
		return nil, &errors.HTTP{
			StatusCode: httpResponse.StatusCode,
			Message:    getErrorString(&httpResponse.Body),
		}
	}

	return &VerifyUserResponse{}, nil
}

// createGetUserQueryString Create string containing query parameters from GetUserRequest
func createGetUserQueryString(request *GetUserRequest) (string, error) {
	if request == nil {
		return "", fmt.Errorf("request cannot be nil")
	}

	if request.UserId == nil && request.Username == nil && request.Email == nil {
		return "", fmt.Errorf("user ID, username, or email must be populated")
	}

	query := url.Values{}
	if request.UserId != nil {
		query.Set(UserIdKey, strconv.Itoa(*request.UserId))
	}
	if request.Username != nil {
		query.Set(UsernameKey, *request.Username)
	}
	if request.Email != nil {
		query.Set(EmailKey, *request.Email)
	}

	return query.Encode(), nil
}

// Create string containing query parameters from GetUsersRequest
func createGetUsersQueryString(request *GetUsersRequest) (string, error) {
	if request == nil {
		return "", fmt.Errorf("request cannot be nil")
	}

	query := url.Values{}
	if request.Limit != nil {
		query.Set(LimitKey, strconv.Itoa(*request.Limit))
	}
	if request.Offset != nil {
		query.Set(OffsetKey, strconv.Itoa(*request.Offset))
	}
	if request.SortField != nil {
		query.Set(SortFieldKey, *request.SortField)
	}
	if request.SortDirection != nil {
		query.Set(SortDirectionKey, *request.SortDirection)
	}

	return query.Encode(), nil
}

// getErrorString get error string from HTTP response body
func getErrorString(responseBody *io.ReadCloser) string {
	if bodyBytes, err := io.ReadAll(*responseBody); err == nil {
		return string(bodyBytes)
	} else {
		return fmt.Sprintf("failed to read response body: %v", err)
	}
}
