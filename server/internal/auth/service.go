package auth

import (
	db "auth/db/generated"
	"common"
	api "common/api/auth"
	"common/api/user"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Service Interface for performing authentication operations
type Service interface {
	Login(context context.Context, request *api.LoginRequest) (*api.LoginResponse, error)
	Logout(context context.Context, request *api.LogoutRequest) (*api.LogoutResponse, error)
	Renew(context context.Context, request *api.RenewRequest) (*api.RenewResponse, error)
	SendVerificationEmail(
		context context.Context,
		request *api.SendVerificationEmailRequest,
	) (*api.SendVerificationEmailResponse, error)
	VerifyEmail(context context.Context, request *api.VerifyEmailRequest) (*api.VerifyEmailResponse, error)
	// endpoint for client JWT?
}

// ServiceImpl Implementation for Service
type ServiceImpl struct {
	Queries db.Querier
	BaseUrl string
}

// Login provides JWT for a user
func (s *ServiceImpl) Login(ctx context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {
	userClient := user.ClientImpl{
		BaseUrl:    s.BaseUrl,
		HttpClient: http.DefaultClient,
	}

	getUserRequest := &user.GetUserRequest{Username: &request.Username}
	jwtStr, err := common.JWTFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve JWT from context: %v", err)
	}
	usr, err := userClient.GetUser(getUserRequest, jwtStr)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve user: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(request.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	token, err := jwt.NewBuilder().
		Issuer("quizchief-auth").
		Subject(usr.Username).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(30*time.Minute)).
		Claim(common.UserIdClaimsKey, usr.UserId).
		Claim(common.UsernameClaimsKey, usr.Username).
		Claim(common.EmailClaimsKey, usr.Email).
		Build()
	if err != nil {
		return nil, fmt.Errorf("unable to build JWT token: %v", err)
	}

	signedToken, err := jwt.Sign(token, jwt.WithKey(common.JWTAlg, common.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("unable to sign JWT token: %v", err)
	}

	refreshToken, err := s.generateRefreshToken(ctx, usr.UserId)
	if err != nil {
		return nil, fmt.Errorf("unable to generate refresh token: %v", err)
	}

	return &api.LoginResponse{
		AccessToken:  string(signedToken),
		RefreshToken: refreshToken,
	}, nil
}

// generateRefreshToken generates refresh token for the user associated with the specified ID
func (s *ServiceImpl) generateRefreshToken(ctx context.Context, userId int) (string, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(data)

	params := db.CreateRefreshTokenParams{
		TokenHash: token,
		UserID:    int32(userId),
	}
	refreshToken, err := s.Queries.CreateRefreshToken(ctx, params)
	if err != nil {
		return "", err
	}
	return refreshToken.TokenHash, nil
}
