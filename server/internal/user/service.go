// Package user contains the implementation for a user-oriented RESTful service
package user

import (
    "context"
    "fmt"
    "golang.org/x/crypto/bcrypt"
    "os"
    "user/db"
    "user/dto"
)

const DEFAULT_USERS_PAGE_LIMIT = 20

var BaseUrl string = os.Getenv("HOST") + "/" + os.Getenv("PORT")

// Service Interface for performing user operations
type Service interface {
    CreateUser(context context.Context, request dto.CreateUserRequest) (dto.CreateUserResponse, error)
    GetUser(context context.Context, request dto.GetUserRequest) (dto.GetUserResponse, error)
    GetUsers(context context.Context, request dto.GetUsersRequest) (dto.GetUsersResponse, error)
    UpdateUser(context context.Context, request dto.UpdateUserRequest) (dto.UpdateUserResponse, error)
    DeleteUser(context context.Context, request dto.DeleteUserRequest) (dto.DeleteUserResponse, error)
}

// ServiceImpl Implementation for the Service
type ServiceImpl struct {
    Queries db.Queries
}

// CreateUser Create a new user
func (service *ServiceImpl) CreateUser(
    context context.Context,
    request dto.CreateUserRequest,
) (dto.CreateUserResponse, error) {
    params := db.CreateUserParams{
        Username: request.Username,
        Email:    request.Email,
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
    if err != nil {
        return dto.CreateUserResponse{}, fmt.Errorf("failed to create user: failed to create password hash: %w", err)
    }
    params.PasswordHash = string(hashedPassword)

    user, err := service.Queries.CreateUser(context, params)
    if err != nil {
        return dto.CreateUserResponse{}, fmt.Errorf("failed to create user: %w", err)
    }

    return dto.CreateUserResponse{
        UserId: int(user.ID),
    }, nil
}

// GetUser Retrieve a user based on ID, username, or email
func (service *ServiceImpl) GetUser(context context.Context, request dto.GetUserRequest) (
    dto.GetUserResponse,
    error,
) {
    params := db.GetUserParams{
        ID:       int32(*request.UserId),
        Username: *request.Username,
        Email:    *request.Email,
    }

    user, err := service.Queries.GetUser(context, params)
    if err != nil {
        return dto.GetUserResponse{}, fmt.Errorf("failed to retrieve user: %w", err)
    }

    return dto.GetUserResponse{
        UserId:     int(user.ID),
        Username:   user.Username,
        Email:      user.Email,
        IsActive:   user.IsActive,
        IsVerified: user.IsVerified,
        CreatedAt:  user.CreatedAt,
        UpdatedAt:  user.UpdatedAt,
    }, nil
}

// GetUsers Retrieve all users (paginated)
func (service *ServiceImpl) GetUsers(
    context context.Context,
    request dto.GetUsersRequest,
) (dto.GetUsersResponse, error) {
    params := db.GetUsersParams{}

    if request.Limit == nil {
        params.Limit = DEFAULT_USERS_PAGE_LIMIT
    } else {
        params.Limit = int32(*request.Limit)
    }

    if request.Offset == nil {
        params.Offset = 0
    } else {
        params.Offset = int32(*request.Offset)
    }

    userCount, err := service.Queries.CountUsers(context)
    users, err := service.Queries.GetUsers(context, params)
    if err != nil {
        return dto.GetUsersResponse{}, fmt.Errorf("failed to retrieve users: %w", err)
    }

    response := dto.GetUsersResponse{Users: make([]dto.GetUserResponse, len(users))}
    for i, user := range users {
        userResponse := dto.GetUserResponse{
            UserId:     int(user.ID),
            Username:   user.Username,
            Email:      user.Email,
            IsVerified: user.IsVerified,
            CreatedAt:  user.CreatedAt,
            UpdatedAt:  user.UpdatedAt,
        }
        response.Users[i] = userResponse
    }

    if request.Offset != nil && *request.Offset > 0 {
        response.PrevLink = BaseUrl + "?limit=" + string(params.Limit) + "&offset=" + string(params.Offset-params.Limit)
    }

    usersRemaining := int(userCount)
    usersRemaining -= int(params.Offset)
    usersRemaining -= int(params.Limit)

    if usersRemaining > 0 {
        response.NextLink = BaseUrl + "?limit=" + string(params.Limit) + "&offset=" + string(params.Offset+params.Limit)
    }

    return response, nil
}

// UpdateUser Update a user
func (service *ServiceImpl) UpdateUser(
    context context.Context,
    request dto.UpdateUserRequest,
) (dto.UpdateUserResponse, error) {
    params := db.UpdateUserParams{
        ID: int32(request.UserId),
    }

    if request.Username != nil {
        params.Username = *request.Username
    }

    if request.Email != nil {
        params.Email = *request.Email
    }

    if request.Password != nil {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
        if err != nil {
            return dto.UpdateUserResponse{}, fmt.Errorf(
                "failed to update user: failed to create password hash: %w",
                err,
            )
        }
        params.PasswordHash = string(hashedPassword)
    }

    user, err := service.Queries.UpdateUser(context, params)
    if err != nil {
        return dto.UpdateUserResponse{}, fmt.Errorf("failed to update user: %w", err)
    }

    response := dto.UpdateUserResponse{
        UserId:     int(user.ID),
        Username:   user.Username,
        Email:      user.Email,
        IsActive:   user.IsActive,
        IsVerified: user.IsVerified,
        CreatedAt:  user.CreatedAt,
        UpdatedAt:  user.UpdatedAt,
    }

    return response, nil
}

// DeleteUser Delete user
func (service *ServiceImpl) DeleteUser(
    context context.Context,
    request dto.DeleteUserRequest,
) (dto.DeleteUserResponse, error) {
    err := service.Queries.DeactivateUser(context, int32(request.UserId))
    if err != nil {
        return dto.DeleteUserResponse{}, fmt.Errorf("failed to delete user: %w", err)
    }
    return dto.DeleteUserResponse{}, nil
}
