package common

import (
	"time"
)

// UserClaims Stores JWT information for User
type UserClaims struct {
	ID       int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type getUserResponse struct {
	UserId       int       `json:"userId"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	IsVerified   bool      `json:"isVerified"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
