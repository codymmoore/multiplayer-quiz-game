package auth

import (
    db "auth/db/generated"
    "common"
    api "common/api/auth"
    "common/api/user"
    "common/errors"
    "context"
    "crypto/rand"
    "database/sql"
    stderrors "errors"
    "fmt"
    "github.com/keighl/postmark"
    "golang.org/x/crypto/bcrypt"
    "net/http"
)

// Service Interface for performing authentication operations
// TODO endpoint for client JWT?
type Service interface {
    Login(ctx context.Context, request *api.LoginRequest) (*api.LoginResponse, error)
    Logout(ctx context.Context, request *api.LogoutRequest) (*api.LogoutResponse, error)
    Renew(ctx context.Context, request *api.RenewRequest) (*api.RenewResponse, error)
    SendVerificationEmail(
        ctx context.Context,
        request *api.SendVerificationEmailRequest,
    ) (*api.SendVerificationEmailResponse, error)
    VerifyEmail(ctx context.Context, request *api.VerifyEmailRequest) (*api.VerifyEmailResponse, error)
}

// ServiceImpl Implementation for Service
type ServiceImpl struct {
    Queries        db.Querier
    BaseUrl        string
    UserClient     user.Client
    PostmarkClient postmark.Client
}

// Login provides JWT for a user
func (s *ServiceImpl) Login(ctx context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {
    getUserRequest := &user.GetUserRequest{Username: &request.Username}
    jwtStr, err := common.JWTFromContext(ctx)
    if err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to retrieve JWT from context: %v", err),
        }
    }
    usr, err := s.UserClient.GetUser(getUserRequest, jwtStr)
    if err != nil {
        var httpErr *errors.HTTP
        if stderrors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
            return nil, &errors.HTTP{
                StatusCode: http.StatusUnauthorized,
                Message:    "invalid credentials",
            }
        }
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to retrieve user: %v", err),
        }
    }

    if err := bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(request.Password)); err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusUnauthorized,
            Message:    "invalid credentials",
        }
    }

    accessToken, err := GenerateJWT(usr.UserId, usr.Username, usr.Email)
    if err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to generate JWT: %v", err),
        }
    }

    refreshToken, err := s.generateRefreshToken(ctx, usr.UserId)
    if err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to generate refresh token: %v", err),
        }
    }

    return &api.LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}

// Logout logs out the current authenticated user
func (s *ServiceImpl) Logout(ctx context.Context, request *api.LogoutRequest) (*api.LogoutResponse, error) {
    tokenHash := SecureHash(request.RefreshToken)
    dbToken, err := s.Queries.GetRefreshToken(ctx, tokenHash)
    if err != nil {
        if stderrors.Is(err, sql.ErrNoRows) {
            return &api.LogoutResponse{}, nil
        }
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to fetch refresh token: %v", err),
        }
    }

    userClaims, ok := ctx.Value(common.UserClaimsCtxKey).(common.UserClaims)
    if !ok {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    "unable to retrieve user claims",
        }
    }

    if int(dbToken.UserID) != userClaims.ID {
        return nil, &errors.HTTP{
            StatusCode: http.StatusUnauthorized,
            Message:    "invalid refresh token",
        }
    }

    if err := s.Queries.DeactivateRefreshToken(ctx, dbToken.ID); err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to deactivate refresh token: %v", err),
        }
    }

    return &api.LogoutResponse{}, nil
}

// Renew provides a new JWT for a user given a refresh token
func (s *ServiceImpl) Renew(ctx context.Context, request *api.RenewRequest) (*api.RenewResponse, error) {
    tokenHash := SecureHash(request.RefreshToken)
    dbToken, err := s.Queries.GetRefreshToken(ctx, tokenHash)
    if err != nil {
        if stderrors.Is(err, sql.ErrNoRows) {
            return nil, &errors.HTTP{
                StatusCode: http.StatusUnauthorized,
                Message:    "invalid refresh token",
            }
        }
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to fetch refresh token: %v", err),
        }
    }

    userClaims, ok := ctx.Value(common.UserClaimsCtxKey).(common.UserClaims)
    if !ok {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    "unable to retrieve user claims",
        }
    }

    if int(dbToken.UserID) != userClaims.ID {
        return nil, &errors.HTTP{
            StatusCode: http.StatusUnauthorized,
            Message:    "invalid refresh token for authenticated user",
        }
    }

    accessToken, err := GenerateJWT(userClaims.ID, userClaims.Username, userClaims.Email)
    if err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to generate JWT: %v", err),
        }
    }

    newRefreshToken, err := s.generateRefreshToken(ctx, userClaims.ID)
    if err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to generate refresh token: %v", err),
        }
    }

    if err := s.Queries.DeactivateRefreshToken(ctx, dbToken.ID); err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to deactivate refresh token: %v", err),
        }
    }

    return &api.RenewResponse{
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken,
    }, nil
}

// SendVerificationEmail sends a verification email to the specified address
func (s *ServiceImpl) SendVerificationEmail(
    ctx context.Context,
    request *api.SendVerificationEmailRequest,
) (*api.SendVerificationEmailResponse, error) {
    getUserRequest := &user.GetUserRequest{Email: &request.Email}
    jwtStr, err := common.JWTFromContext(ctx)
    if err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to retrieve JWT from context: %v", err),
        }
    }
    usr, err := s.UserClient.GetUser(getUserRequest, jwtStr)
    if err != nil {
        var httpErr *errors.HTTP
        if stderrors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
            return nil, &errors.HTTP{
                StatusCode: http.StatusNotFound,
                Message:    "no user with specified email found",
            }
        }
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to retrieve user: %v", err),
        }
    }

    verificationCode, err := s.generateVerificationCode(ctx, usr.UserId)
    if err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to generate verification code: %v", err),
        }
    }

    verificationLink := fmt.Sprintf(
        "%s%s?%s=%s",
        s.BaseUrl,
        api.VerifyEmailEndpoint,
        api.VerificationCodeQueryParamId,
        verificationCode,
    )
    email := postmark.Email{
        From:    "no-reply@quizchief.gg",
        To:      usr.Email,
        Subject: "Email Verification",
        HtmlBody: fmt.Sprintf(
            "<p>Click here to verify your email: <a href='%s'>%s</a></p>",
            verificationLink,
            verificationLink,
        ),
        TextBody:   fmt.Sprintf("Click here to verify your email: %s", verificationLink),
        Tag:        "email-verification",
        TrackOpens: true,
    }
    if _, err := s.PostmarkClient.SendEmail(email); err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to send verification email: %v", err),
        }
    }

    return &api.SendVerificationEmailResponse{}, nil
}

// VerifyEmail verifies the user associated with the specified verification code
func (s *ServiceImpl) VerifyEmail(ctx context.Context, request *api.VerifyEmailRequest) (
    *api.VerifyEmailResponse,
    error,
) {
    codeHash := SecureHash(request.VerificationCode)
    dbCode, err := s.Queries.GetVerificationCode(ctx, codeHash)
    if err != nil {
        if stderrors.Is(err, sql.ErrNoRows) {
            return nil, &errors.HTTP{
                StatusCode: http.StatusNotFound,
                Message:    "invalid verification code",
            }
        }
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to retrieve verification code: %v", err),
        }
    }

    jwtStr, err := common.JWTFromContext(ctx)
    if err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to retrieve JWT from context: %v", err),
        }
    }

    userId := int(dbCode.UserID)
    verifyUserRequest := &user.VerifyUserRequest{
        UserId: userId,
    }
    if _, err := s.UserClient.VerifyUser(verifyUserRequest, jwtStr); err != nil {
        return nil, &errors.HTTP{
            StatusCode: http.StatusInternalServerError,
            Message:    fmt.Sprintf("unable to verify user: %v", err),
        }
    }

    return &api.VerifyEmailResponse{
        UserId: userId,
    }, nil
}

func (s *ServiceImpl) generateRefreshToken(ctx context.Context, userId int) (string, error) {
    data := make([]byte, 32)
    if _, err := rand.Read(data); err != nil {
        return "", err
    }
    token := fmt.Sprintf("%x", data)

    params := db.CreateRefreshTokenParams{
        TokenHash: SecureHash(token),
        UserID:    int32(userId),
    }

    _, err := s.Queries.CreateRefreshToken(ctx, params)
    if err != nil {
        return "", err
    }

    return token, nil
}

func (s *ServiceImpl) generateVerificationCode(ctx context.Context, userId int) (string, error) {
    data := make([]byte, 16)
    if _, err := rand.Read(data); err != nil {
        return "", err
    }
    verificationCode := fmt.Sprintf("%x", data)

    params := db.UpsertVerificationCodeParams{
        UserID:               int32(userId),
        VerificationCodeHash: SecureHash(verificationCode),
    }
    _, err := s.Queries.UpsertVerificationCode(ctx, params)
    if err != nil {
        return "", err
    }

    return verificationCode, nil
}
