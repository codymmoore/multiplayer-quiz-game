package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	CreateUser(request *CreateUserRequest) (*CreateUserResponse, error)
	GetUser(request *GetUserRequest) (*GetUserResponse, error)
	GetUsers(request *GetUsersRequest) (*GetUsersResponse, error)
	UpdateUser(request *UpdateUserRequest) (*UpdateUserResponse, error)
	DeleteUser(request *DeleteUserRequest) (*DeleteUserResponse, error)
}

type ClientImpl struct {
	BaseUrl    string
	HttpClient http.Client
}

// TODO Add authentication
func (client *ClientImpl) CreateUser(request *CreateUserRequest) (*CreateUserResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CreateUserRequest: %w", err)
	}

	url := client.BaseUrl + "/user"
	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.HttpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	var response CreateUserResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode CreateUserResponse: %w", err)
	}

	return &response, nil
}
