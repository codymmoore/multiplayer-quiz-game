package auth

import (
	db "auth/db/generated"
	"common"
	api "common/api/auth"
	"common/api/user"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"errors"
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

	accessToken, err := generateJWT(usr.UserId, usr.Username, usr.Email)
	if err != nil {
		return nil, fmt.Errorf("unable to generate JWT: %v", err)
	}

	refreshToken, err := s.generateRefreshToken(ctx, usr.UserId)
	if err != nil {
		return nil, fmt.Errorf("unable to generate refresh token: %v", err)
	}

	return &api.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout logs out the current authenticated user
func (s *ServiceImpl) Logout(ctx context.Context, request *api.LogoutRequest) (*api.LogoutResponse, error) {
	tokenHash := hashToken(request.RefreshToken)
	dbToken, err := s.Queries.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid refresh token")
		}
		return nil, fmt.Errorf("unable to fetch refresh token: %w", err)
	}

	if err := s.Queries.DeactivateRefreshToken(ctx, dbToken.ID); err != nil {
		return nil, fmt.Errorf("unable to deactivate refresh token: %v", err)
	}

	return &api.LogoutResponse{}, nil
}

// Renew provides a new JWT for a user given a refresh token
func (s *ServiceImpl) Renew(ctx context.Context, request *api.RenewRequest) (*api.RenewResponse, error) {
	tokenHash := hashToken(request.RefreshToken)
	dbToken, err := s.Queries.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid refresh token")
		}
		return nil, fmt.Errorf("unable to fetch refresh token: %w", err)
	}

	userClaims, ok := ctx.Value(common.UserClaimsCtxKey).(common.UserClaims)
	if !ok {
		return nil, errors.New("unable to retrieve user claims")
	}

	if int(dbToken.UserID) != userClaims.ID {
		return nil, errors.New("invalid refresh token for authenticated user")
	}

	accessToken, err := generateJWT(userClaims.ID, userClaims.Username, userClaims.Email)
	if err != nil {
		return nil, fmt.Errorf("unable to generate JWT: %v", err)
	}

	newRefreshToken, err := s.generateRefreshToken(ctx, userClaims.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to generate refresh token: %v", err)
	}

	if err := s.Queries.DeactivateRefreshToken(ctx, dbToken.ID); err != nil {
		return nil, fmt.Errorf("unable to deactivate refresh token: %v", err)
	}

	return &api.RenewResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *ServiceImpl) generateRefreshToken(ctx context.Context, userId int) (string, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	token := fmt.Sprintf("%x", data)

	params := db.CreateRefreshTokenParams{
		TokenHash: hashToken(token),
		UserID:    int32(userId),
	}

	_, err := s.Queries.CreateRefreshToken(ctx, params)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateJWT(userId int, username string, email string) (string, error) {
	token, err := jwt.NewBuilder().
		Issuer("quizchief-auth").
		Subject(username).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(30*time.Minute)).
		Claim(common.UserIdClaimsKey, userId).
		Claim(common.UsernameClaimsKey, username).
		Claim(common.EmailClaimsKey, email).
		Build()
	if err != nil {
		return "", err
	}

	signedToken, err := jwt.Sign(token, jwt.WithKey(common.JWTAlg, common.JWTSecret))
	if err != nil {
		return "", err
	}

	return string(signedToken), nil
}

func hashToken(token string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
}
