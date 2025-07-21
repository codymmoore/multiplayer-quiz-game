// Package user contains the implementation for a user-oriented RESTful service
package user

import (
	"common"
	"context"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"user/db/generated"
	"user/dto"
)

const DefaultUsersPageLimit = 20

// Service Interface for performing user operations
type Service interface {
	CreateUser(context context.Context, request *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	GetUser(context context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error)
	GetUsers(context context.Context, request *dto.GetUsersRequest) (*dto.GetUsersResponse, error)
	UpdateUser(context context.Context, request *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
	DeleteUser(context context.Context, request *dto.DeleteUserRequest) (*dto.DeleteUserResponse, error)
}

// ServiceImpl Implementation for the Service
type ServiceImpl struct {
	Queries db.Querier
}

// CreateUser Create a new user
func (service *ServiceImpl) CreateUser(
	context context.Context,
	request *dto.CreateUserRequest,
) (*dto.CreateUserResponse, error) {
	params := db.CreateUserParams{
		Username: request.Username,
		Email:    request.Email,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to create password hash: %w", err)
	}
	params.PasswordHash = string(hashedPassword)

	user, err := service.Queries.CreateUser(context, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &dto.CreateUserResponse{
		UserId: int(user.ID),
	}, nil
}

// GetUser Retrieve a user based on ID, username, or email
// TODO Change query to return many
func (service *ServiceImpl) GetUser(context context.Context, request *dto.GetUserRequest) (
	*dto.GetUserResponse,
	error,
) {
	var params db.GetUserParams

	if userId := request.UserId; userId != nil {
		params.ID = sql.NullInt32{Int32: int32(*userId), Valid: true}
	} else {
		params.ID = sql.NullInt32{Valid: false}
	}

	if username := request.Username; username != nil {
		params.Username = sql.NullString{String: *username, Valid: true}
	} else {
		params.Username = sql.NullString{Valid: false}
	}

	if email := request.Email; email != nil {
		params.Email = sql.NullString{String: *email, Valid: true}
	} else {
		params.Email = sql.NullString{Valid: false}
	}

	user, err := service.Queries.GetUser(context, params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return &dto.GetUserResponse{
		UserId:       int(user.ID),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsVerified:   user.IsVerified,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// GetUsers Retrieve all users (paginated)
// TODO implement sorting
func (service *ServiceImpl) GetUsers(
	context context.Context,
	request *dto.GetUsersRequest,
) (*dto.GetUsersResponse, error) {
	params := db.GetUsersParams{}

	if request.Limit == nil {
		params.Limit = DefaultUsersPageLimit
	} else {
		params.Limit = int32(*request.Limit)
	}

	if request.Offset == nil {
		params.Offset = 0
	} else {
		params.Offset = int32(*request.Offset)
	}

	userCount, err := service.Queries.CountUsers(context)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	users, err := service.Queries.GetUsers(context, params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	response := dto.GetUsersResponse{Users: make([]dto.GetUserResponse, len(users))}
	for i, user := range users {
		response.Users[i] = dto.GetUserResponse{
			UserId:       int(user.ID),
			Username:     user.Username,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			IsVerified:   user.IsVerified,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		}
	}

	routeUrl, err := common.GetRouteUrl(context)
	if err != nil {
		return nil, fmt.Errorf("failed to get route url: %w", err)
	}

	if request.Offset != nil && *request.Offset > 0 {
		prevOffset := params.Offset - params.Limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		prevLink := fmt.Sprintf(
			"%s?limit=%d&offset=%d",
			routeUrl,
			params.Limit,
			prevOffset,
		)
		response.PrevLink = &prevLink
	}

	usersRemaining := int(userCount)
	usersRemaining -= int(params.Offset)
	usersRemaining -= int(params.Limit)

	if usersRemaining > 0 {
		nextLink := fmt.Sprintf(
			"%s?limit=%d&offset=%d",
			routeUrl,
			params.Limit,
			params.Offset+params.Limit,
		)
		response.NextLink = &nextLink
	}

	return &response, nil
}

// UpdateUser Update a user
func (service *ServiceImpl) UpdateUser(
	context context.Context,
	request *dto.UpdateUserRequest,
) (*dto.UpdateUserResponse, error) {
	params := db.UpdateUserParams{
		ID: int32(request.UserId),
	}

	if request.Username == nil {
		params.Username = sql.NullString{String: "", Valid: false}
	} else {
		params.Username = sql.NullString{String: *request.Username, Valid: true}
	}

	if request.Email == nil {
		params.Email = sql.NullString{String: "", Valid: false}
	} else {
		params.Email = sql.NullString{String: *request.Email, Valid: true}
	}

	if request.Password == nil {
		params.PasswordHash = sql.NullString{String: "", Valid: false}
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to create password hash: %w",
				err,
			)
		}
		params.PasswordHash = sql.NullString{String: string(hashedPassword), Valid: true}
	}

	user, err := service.Queries.UpdateUser(context, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &dto.UpdateUserResponse{
		UserId:       int(user.ID),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsVerified:   user.IsVerified,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// DeleteUser Delete user
func (service *ServiceImpl) DeleteUser(
	context context.Context,
	request *dto.DeleteUserRequest,
) (*dto.DeleteUserResponse, error) {
	err := service.Queries.DeleteUser(context, int32(request.UserId))
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	return &dto.DeleteUserResponse{}, nil
}
