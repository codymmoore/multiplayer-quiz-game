package user

import (
	"time"
)

type CreateUserResponse struct {
	UserId int `json:"userId"`
}

type User struct {
	UserId       int       `json:"userId"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	IsVerified   bool      `json:"isVerified"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type GetUserResponse = User

type GetUsersResponse struct {
	Users    []User  `json:"users"`
	PrevLink *string `json:"prevLink,omitempty"`
	NextLink *string `json:"nextLink,omitempty"`
}

type UpdateUserResponse = User

type DeleteUserResponse struct {
}
