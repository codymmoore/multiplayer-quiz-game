package dto

import (
	"time"
)

type CreateUserResponse struct {
	UserId int `json:"userId"`
}

type GetUserResponse struct {
	UserId     int       `json:"userId"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	IsActive   bool      `json:"isActive"`
	IsVerified bool      `json:"isVerified"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type GetUsersResponse struct {
	Users    []GetUserResponse `json:"users"`
	PrevLink string            `json:"prevLink"`
	NextLink string            `json:"nextLink"`
}

type UpdateUserResponse struct {
	UserId     int       `json:"userId"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	IsActive   bool      `json:"isActive"`
	IsVerified bool      `json:"isVerified"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type DeleteUserResponse struct {
}
