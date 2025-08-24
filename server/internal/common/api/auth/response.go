package auth

import (
	"time"
)

type User struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AuthInfo struct {
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	IssuedAt     time.Time     `json:"issuedAt"`
	ExpiresIn    time.Duration `json:"expiresIn"`
	User         User          `json:"user"`
}

type LoginResponse = AuthInfo

type LogoutResponse struct{}

type RenewResponse = AuthInfo

type SendVerificationEmailResponse struct {
}

type VerifyEmailResponse struct {
}
