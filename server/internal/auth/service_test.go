package auth

import (
	db "auth/db/generated"
	"common"
	api "common/api/auth"
	"common/api/user"
	"common/errors"
	"common/test"
	"context"
	"database/sql"
	stderrors "errors"
	"github.com/keighl/postmark"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	JWT          = "JWT"
	JWTSecretKey = "JWT_SECRET"
)

var querier = &mockQuerier{}
var userClient = &test.MockUserClient{}
var postmarkClient = postmark.NewClient("mockServerToken", "mockAuthToken")
var service = ServiceImpl{
	Queries:        querier,
	BaseUrl:        "https://mock-url.com",
	UserClient:     userClient,
	PostmarkClient: postmarkClient,
}

func TestService_Login_Success(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword

	ctx, request := loginSetup(t, userId, username, email, password)
	response, err := service.Login(ctx, request)

	if err != nil {
		t.Errorf(`service.Login(ctx, request) error = "%v", expected "<nil>"`, err)
	}

	if response == nil {
		t.Error(`service.Login(ctx, request) response = "<nil>", expected non-nil`)
		return
	}
	if response.AccessToken == "" {
		t.Error(`response.AccessToken = "", expected non-empty string`)
	}
	if response.RefreshToken == "" {
		t.Error(`response.RefreshToken = "", expected non-empty string`)
	}
	if response.IssuedAt.IsZero() && response.IssuedAt.Before(time.Now()) {
		t.Errorf(`response.IssuedAt = "%v", expected non-zero time`, response.IssuedAt)
	}

	expectedExpiresIn := time.Minute * 30
	if response.ExpiresIn != expectedExpiresIn {
		t.Errorf(`response.ExpiresIn = "%v", expected "%v"`, response.ExpiresIn, expectedExpiresIn)
	}

	expectedUser := api.User{
		UserId:   userId,
		Username: username,
		Email:    email,
	}
	if response.User != expectedUser {
		t.Errorf(`response.User = "%v", expected "%v"`, response.User, expectedUser)
	}
}

func TestService_Login_MissingUser(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword

	ctx, request := loginSetup(t, userId, username, email, password)
	userClient.GetUserFunc = func(request *user.GetUserRequest) (*user.GetUserResponse, error) {
		return nil, &errors.HTTP{StatusCode: http.StatusNotFound, Message: "Not Found"}
	}

	response, err := service.Login(ctx, request)

	if err == nil {
		t.Error(`service.Login(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusUnauthorized)

	if response != nil {
		t.Errorf(`service.Login(ctx, request) response = "%v", expected nil`, response)
	}
}

func TestService_Login_GetUserError(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword

	ctx, request := loginSetup(t, userId, username, email, password)
	userClient.GetUserFunc = func(request *user.GetUserRequest) (*user.GetUserResponse, error) {
		return nil, &errors.HTTP{StatusCode: http.StatusInternalServerError, Message: "Error"}
	}

	response, err := service.Login(ctx, request)

	if err == nil {
		t.Error(`service.Login(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.Login(ctx, request) response = "%v", expected nil`, response)
	}
}

func TestService_Login_UnverifiedUser(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword

	ctx, request := loginSetup(t, userId, username, email, password)
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	userClient.GetUserFunc = func(request *user.GetUserRequest) (*user.GetUserResponse, error) {
		return &user.GetUserResponse{
			UserId:       userId,
			Username:     username,
			Email:        email,
			PasswordHash: string(passwordHash),
			IsVerified:   false,
		}, nil
	}

	response, err := service.Login(ctx, request)

	if err == nil {
		t.Error(`service.Login(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusUnauthorized)

	if response != nil {
		t.Errorf(`service.Login(ctx, request) response = "%v", expected nil`, response)
	}
}

func TestService_Login_IncorrectPassword(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword

	ctx, request := loginSetup(t, userId, username, email, password)

	userClient.GetUserFunc = func(request *user.GetUserRequest) (*user.GetUserResponse, error) {
		return &user.GetUserResponse{
			UserId:       userId,
			Username:     username,
			Email:        email,
			PasswordHash: "invalid",
		}, nil
	}

	response, err := service.Login(ctx, request)

	if err == nil {
		t.Error(`service.Login(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusUnauthorized)

	if response != nil {
		t.Errorf(`service.Login(ctx, request) response = "%v", expected nil`, response)
	}
}

func TestService_Login_JWTSignError(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword

	ctx, request := loginSetup(t, userId, username, email, password)
	common.JWTSecret = nil
	response, err := service.Login(ctx, request)

	if err == nil {
		t.Error(`service.Login(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.Login(ctx, request) response = "%v", expected nil`, response)
	}
}

func TestService_Login_RefreshTokenGenerationError(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	password := test.ValidPassword

	ctx, request := loginSetup(t, userId, username, email, password)

	querier.createRefreshTokenFunc = func(ctx context.Context, arg db.CreateRefreshTokenParams) (
		db.RefreshToken,
		error,
	) {
		return db.RefreshToken{}, stderrors.New("")
	}

	response, err := service.Login(ctx, request)

	if err == nil {
		t.Error(`service.Login(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.Login(ctx, request) response = "%v", expected nil`, response)
	}
}

func TestService_Logout_Success(t *testing.T) {
	userId := 1
	ctx := logoutSetup(t, userId)
	request := &api.LogoutRequest{RefreshToken: "refreshToken"}
	response, err := service.Logout(ctx, request)

	if err != nil {
		t.Errorf(`service.Logout(ctx, request) error = "%v", expected "<nil>"`, err)
	}

	if response == nil {
		t.Error(`service.Logout(ctx, request) response = "<nil>", expected non-nil`)
	}
}

func TestService_Logout_TokenNotFound(t *testing.T) {
	userId := 1
	ctx := logoutSetup(t, userId)

	querier.getRefreshTokenFunc = func(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
		return db.RefreshToken{}, sql.ErrNoRows
	}

	request := &api.LogoutRequest{RefreshToken: "refreshToken"}
	response, err := service.Logout(ctx, request)

	if err != nil {
		t.Errorf(`service.Logout(ctx, request) error = "%v", expected "<nil>"`, err)
	}

	if response == nil {
		t.Error(`service.Logout(ctx, request) response = "<nil>", expected non-nil`)
	}
}

func TestService_Logout_RetrieveTokenError(t *testing.T) {
	userId := 1
	ctx := logoutSetup(t, userId)

	querier.getRefreshTokenFunc = func(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
		return db.RefreshToken{}, stderrors.New("")
	}

	request := &api.LogoutRequest{RefreshToken: "refreshToken"}
	response, err := service.Logout(ctx, request)

	if err == nil {
		t.Error(`service.Logout(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.Logout(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_Logout_MissingUserClaims(t *testing.T) {
	userId := 1
	logoutSetup(t, userId)
	request := &api.LogoutRequest{RefreshToken: "refreshToken"}
	response, err := service.Logout(context.Background(), request)

	if err == nil {
		t.Error(`service.Logout(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusUnauthorized)

	if response != nil {
		t.Errorf(`service.Logout(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_Logout_WrongAuthenticatedUser(t *testing.T) {
	userId := 1
	logoutSetup(t, userId)

	ctx := context.WithValue(context.Background(), common.UserClaimsCtxKey, userId+1)

	request := &api.LogoutRequest{RefreshToken: "refreshToken"}
	response, err := service.Logout(ctx, request)

	if err == nil {
		t.Error(`service.Logout(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusUnauthorized)

	if response != nil {
		t.Errorf(`service.Logout(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_Logout_DeactivateTokenError(t *testing.T) {
	userId := 1
	ctx := logoutSetup(t, userId)

	querier.deactivateRefreshTokenFunc = func(ctx context.Context, id int32) error {
		return stderrors.New("")
	}

	request := &api.LogoutRequest{RefreshToken: "refreshToken"}
	response, err := service.Logout(ctx, request)

	if err == nil {
		t.Error(`service.Logout(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.Logout(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_Renew_Success(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail

	ctx, request := renewSetup(t, userId, username, email)
	response, err := service.Renew(ctx, request)

	if err != nil {
		t.Errorf(`service.Renew(ctx, request) error = "%v", expected "<nil>"`, err)
	}

	if response == nil {
		t.Error(`service.Renew(ctx, request) response = "<nil>", expected non-nil`)
		return
	}
	if response.AccessToken == "" {
		t.Error(`response.AccessToken = "", expected non-empty string`)
	}
	if response.RefreshToken == "" {
		t.Error(`response.RefreshToken = "", expected non-empty string`)
	}
	if response.IssuedAt.IsZero() && response.IssuedAt.Before(time.Now()) {
		t.Errorf(`response.IssuedAt = "%v", expected non-zero time`, response.IssuedAt)
	}

	expectedExpiresIn := time.Minute * 30
	if response.ExpiresIn != expectedExpiresIn {
		t.Errorf(`response.ExpiresIn = "%v", expected "%v"`, response.ExpiresIn, expectedExpiresIn)
	}

	expectedUser := api.User{
		UserId:   userId,
		Username: username,
		Email:    email,
	}
	if response.User != expectedUser {
		t.Errorf(`response.User = "%v", expected "%v"`, response.User, expectedUser)
	}
}

func TestService_Renew_TokenNotFound(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail

	ctx, request := renewSetup(t, userId, username, email)
	querier.getRefreshTokenFunc = func(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
		return db.RefreshToken{}, sql.ErrNoRows
	}

	response, err := service.Renew(ctx, request)

	if err == nil {
		t.Error(`service.Renew(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusUnauthorized)

	if response != nil {
		t.Errorf(`service.Renew(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_Renew_GetRefreshTokenError(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail

	ctx, request := renewSetup(t, userId, username, email)
	querier.getRefreshTokenFunc = func(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
		return db.RefreshToken{}, stderrors.New("")
	}

	response, err := service.Renew(ctx, request)

	if err == nil {
		t.Error(`service.Renew(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.Renew(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_Renew_MissingUserClaims(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail

	_, request := renewSetup(t, userId, username, email)
	response, err := service.Renew(context.Background(), request)

	if err == nil {
		t.Error(`service.Renew(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusUnauthorized)

	if response != nil {
		t.Errorf(`service.Renew(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_Renew_JWTSignError(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail

	ctx, request := renewSetup(t, userId, username, email)
	common.JWTSecret = nil

	response, err := service.Renew(ctx, request)

	if err == nil {
		t.Error(`service.Renew(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.Renew(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_Renew_CreateRefreshTokenError(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail

	ctx, request := renewSetup(t, userId, username, email)
	querier.createRefreshTokenFunc = func(ctx context.Context, args db.CreateRefreshTokenParams) (
		db.RefreshToken,
		error,
	) {
		return db.RefreshToken{}, stderrors.New("")
	}

	response, err := service.Renew(ctx, request)

	if err == nil {
		t.Error(`service.Renew(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.Renew(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_Renew_DeactivateRefreshTokenError(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail

	ctx, request := renewSetup(t, userId, username, email)
	querier.deactivateRefreshTokenFunc = func(ctx context.Context, id int32) error {
		return stderrors.New("")
	}

	response, err := service.Renew(ctx, request)

	if err == nil {
		t.Error(`service.Renew(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.Renew(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_SendVerificationEmail_Success(t *testing.T) {
	userId := 1
	email := test.ValidEmail

	ctx, request := sendVerificationEmailSetup(t, userId, email)
	response, err := service.SendVerificationEmail(ctx, request)

	if err != nil {
		t.Errorf(`service.SendVerificationEmail(ctx, request) error = "%v", expected "<nil>"`, err)
	}

	if response == nil {
		t.Error(`service.SendVerificationEmail(ctx, request) response = "<nil>", expected non-nil`)
		return
	}
}

func TestService_SendVerificationEmail_MissingUser(t *testing.T) {
	userId := 1
	email := test.ValidEmail

	ctx, request := sendVerificationEmailSetup(t, userId, email)
	userClient.GetUserFunc = func(request *user.GetUserRequest) (*user.GetUserResponse, error) {
		return nil, &errors.HTTP{StatusCode: http.StatusNotFound, Message: ""}
	}

	response, err := service.SendVerificationEmail(ctx, request)

	if err == nil {
		t.Error(`service.SendVerificationEmail(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusNotFound)

	if response != nil {
		t.Errorf(`service.SendVerificationEmail(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_SendVerificationEmail_GetUserError(t *testing.T) {
	userId := 1
	email := test.ValidEmail

	ctx, request := sendVerificationEmailSetup(t, userId, email)
	userClient.GetUserFunc = func(request *user.GetUserRequest) (*user.GetUserResponse, error) {
		return nil, stderrors.New("")
	}

	response, err := service.SendVerificationEmail(ctx, request)

	if err == nil {
		t.Error(`service.SendVerificationEmail(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.SendVerificationEmail(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_SendVerificationEmail_GenerateVerificationCodeError(t *testing.T) {
	userId := 1
	email := test.ValidEmail

	ctx, request := sendVerificationEmailSetup(t, userId, email)
	querier.upsertVerificationCodeFunc = func(ctx context.Context, arg db.UpsertVerificationCodeParams) (
		db.VerificationCode,
		error,
	) {
		return db.VerificationCode{}, stderrors.New("")
	}

	response, err := service.SendVerificationEmail(ctx, request)

	if err == nil {
		t.Error(`service.SendVerificationEmail(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.SendVerificationEmail(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_SendVerificationEmail_SendEmailError(t *testing.T) {
	userId := 1
	email := test.ValidEmail

	ctx, request := sendVerificationEmailSetup(t, userId, email)
	postmarkClient.HTTPClient = test.GetMockErrorHttpClient()

	response, err := service.SendVerificationEmail(ctx, request)

	if err == nil {
		t.Error(`service.SendVerificationEmail(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.SendVerificationEmail(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_VerifyEmail_Success(t *testing.T) {
	verificationCode := "verificationCode"
	userId := 1

	ctx, request := verifyEmailSetup(t, verificationCode, userId)
	response, err := service.VerifyEmail(ctx, request)

	if err != nil {
		t.Errorf(`service.VerifyEmail(ctx, request) error = "%v", expected "<nil>"`, err)
	}

	if response == nil {
		t.Error(`service.VerifyEmail(ctx, request) response = "<nil>", expected non-nil`)
		return
	}
	if response.UserId != userId {
		t.Errorf(`response.UserId = "%v", expected "%v"`, response.UserId, userId)
	}
}

func TestService_VerifyEmail_MissingVerificationCode(t *testing.T) {
	verificationCode := "verificationCode"
	userId := 1

	ctx, request := verifyEmailSetup(t, verificationCode, userId)
	querier.getVerificationCodeFunc = func(ctx context.Context, codeHash string) (db.VerificationCode, error) {
		return db.VerificationCode{}, sql.ErrNoRows
	}

	response, err := service.VerifyEmail(ctx, request)

	if err == nil {
		t.Error(`service.VerifyEmail(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusNotFound)

	if response != nil {
		t.Errorf(`service.VerifyEmail(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_VerifyEmail_GetVerificationCodeError(t *testing.T) {
	verificationCode := "verificationCode"
	userId := 1

	ctx, request := verifyEmailSetup(t, verificationCode, userId)
	querier.getVerificationCodeFunc = func(ctx context.Context, codeHash string) (db.VerificationCode, error) {
		return db.VerificationCode{}, stderrors.New("")
	}

	response, err := service.VerifyEmail(ctx, request)

	if err == nil {
		t.Error(`service.VerifyEmail(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusInternalServerError)

	if response != nil {
		t.Errorf(`service.VerifyEmail(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_VerifyEmail_JWTMissingFromContext(t *testing.T) {
	verificationCode := "verificationCode"
	userId := 1

	_, request := verifyEmailSetup(t, verificationCode, userId)
	response, err := service.VerifyEmail(context.Background(), request)

	if err == nil {
		t.Error(`service.VerifyEmail(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusUnauthorized)

	if response != nil {
		t.Errorf(`service.VerifyEmail(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func TestService_VerifyEmail_VerifyUserError(t *testing.T) {
	verificationCode := "verificationCode"
	userId := 1

	_, request := verifyEmailSetup(t, verificationCode, userId)
	userClient.VerifyUserFunc = func(request *user.VerifyUserRequest, jwt string) (*user.VerifyUserResponse, error) {
		return nil, stderrors.New("")
	}

	response, err := service.VerifyEmail(context.Background(), request)

	if err == nil {
		t.Error(`service.VerifyEmail(ctx, request) error = "<nil>", expected non-nil`)
	}
	test.AssertHTTPError(t, err, http.StatusUnauthorized)

	if response != nil {
		t.Errorf(`service.VerifyEmail(ctx, request) response = "%v", expected "<nil>"`, response)
	}
}

func loginSetup(t *testing.T, userId int, username string, email string, password string) (
	context.Context,
	*api.LoginRequest,
) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	userClient.GetUserFunc = func(request *user.GetUserRequest) (*user.GetUserResponse, error) {
		if *request.Username != username {
			t.Errorf(`request.Username = "%v", expected "%v"`, *request.Username, username)
		}

		return &user.GetUserResponse{
			UserId:       userId,
			Username:     username,
			Email:        email,
			PasswordHash: string(passwordHash),
			IsVerified:   true,
		}, nil
	}

	_ = os.Setenv(JWTSecretKey, JWT)
	_ = common.InitJWT()

	querier.createRefreshTokenFunc = func(ctx context.Context, arg db.CreateRefreshTokenParams) (
		db.RefreshToken,
		error,
	) {
		if arg.UserID != int32(userId) {
			t.Errorf(`arg.UserID = "%v", expected "%v"`, arg.UserID, userId)
		}
		return db.RefreshToken{}, nil
	}

	ctx := context.WithValue(context.Background(), common.JWTCtxKey, JWT)
	request := &api.LoginRequest{
		Username: username,
		Password: password,
	}

	return ctx, request
}

func logoutSetup(t *testing.T, userId int) context.Context {
	refreshTokenId := int32(1)
	querier.getRefreshTokenFunc = func(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
		return db.RefreshToken{ID: refreshTokenId, UserID: int32(userId)}, nil
	}

	userClaims := common.UserClaims{ID: userId}
	ctx := context.WithValue(context.Background(), common.UserClaimsCtxKey, userClaims)

	querier.deactivateRefreshTokenFunc = func(ctx context.Context, id int32) error {
		if id != refreshTokenId {
			t.Errorf(`id = "%v", expected "%v"`, id, refreshTokenId)
		}
		return nil
	}

	return ctx
}

func renewSetup(t *testing.T, userId int, username string, email string) (context.Context, *api.RenewRequest) {
	var refreshTokenID int32 = 1
	refreshToken := "refreshToken"

	querier.getRefreshTokenFunc = func(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
		expectedTokenHash := SecureHash(refreshToken)
		if tokenHash != expectedTokenHash {
			t.Errorf(`tokenHash = "%v", expected "%v"`, tokenHash, expectedTokenHash)
		}
		return db.RefreshToken{ID: refreshTokenID, UserID: int32(userId)}, nil
	}

	_ = os.Setenv(JWTSecretKey, JWT)
	_ = common.InitJWT()

	querier.createRefreshTokenFunc = func(ctx context.Context, arg db.CreateRefreshTokenParams) (
		db.RefreshToken,
		error,
	) {
		return db.RefreshToken{}, nil
	}

	querier.deactivateRefreshTokenFunc = func(ctx context.Context, id int32) error {
		if id != refreshTokenID {
			t.Errorf(`id = "%v", expected "%v"`, id, refreshTokenID)
		}
		return nil
	}

	userClaims := common.UserClaims{
		ID:       userId,
		Username: username,
		Email:    email,
	}
	ctx := context.WithValue(context.Background(), common.UserClaimsCtxKey, userClaims)
	request := &api.RenewRequest{RefreshToken: refreshToken}

	return ctx, request
}

func sendVerificationEmailSetup(t *testing.T, userId int, email string) (
	context.Context,
	*api.SendVerificationEmailRequest,
) {
	userClient.GetUserFunc = func(request *user.GetUserRequest) (*user.GetUserResponse, error) {
		return &user.GetUserResponse{
			UserId: 1,
			Email:  email,
		}, nil
	}

	querier.upsertVerificationCodeFunc = func(
		ctx context.Context,
		arg db.UpsertVerificationCodeParams,
	) (db.VerificationCode, error) {
		if arg.UserID != int32(userId) {
			t.Errorf(`arg.UserID = "%v", expected "%v"`, arg.UserID, userId)
		}
		return db.VerificationCode{}, nil
	}

	postmarkClient.HTTPClient = test.GetMockHttpClient(http.StatusOK, "{}")

	ctx := context.WithValue(context.Background(), common.JWTCtxKey, JWT)
	request := &api.SendVerificationEmailRequest{Email: email}

	return ctx, request
}

func verifyEmailSetup(t *testing.T, verificationCode string, userId int) (
	context.Context,
	*api.VerifyEmailRequest,
) {
	querier.getVerificationCodeFunc = func(ctx context.Context, tokenHash string) (db.VerificationCode, error) {
		expectedTokenHash := SecureHash(verificationCode)
		if tokenHash != expectedTokenHash {
			t.Errorf(`tokenHash = "%v", expected "%v"`, verificationCode, expectedTokenHash)
		}
		return db.VerificationCode{UserID: int32(userId)}, nil
	}

	userClient.VerifyUserFunc = func(request *user.VerifyUserRequest, jwt string) (*user.VerifyUserResponse, error) {
		if request.UserId != userId {
			t.Errorf(`request.UserID = "%v", expected "%v"`, request.UserId, userId)
		}
		return &user.VerifyUserResponse{}, nil
	}

	ctx := context.WithValue(context.Background(), common.JWTCtxKey, JWT)
	request := &api.VerifyEmailRequest{VerificationCode: verificationCode}
	return ctx, request
}

type mockQuerier struct {
	createRefreshTokenFunc     func(ctx context.Context, arg db.CreateRefreshTokenParams) (db.RefreshToken, error)
	deactivateRefreshTokenFunc func(ctx context.Context, id int32) error
	getRefreshTokenFunc        func(ctx context.Context, tokenHash string) (db.RefreshToken, error)
	getVerificationCodeFunc    func(ctx context.Context, verificationCodeHash string) (db.VerificationCode, error)
	upsertVerificationCodeFunc func(ctx context.Context, arg db.UpsertVerificationCodeParams) (
		db.VerificationCode,
		error,
	)
}

func (m *mockQuerier) CreateRefreshToken(ctx context.Context, arg db.CreateRefreshTokenParams) (
	db.RefreshToken,
	error,
) {
	return m.createRefreshTokenFunc(ctx, arg)
}

func (m *mockQuerier) DeactivateRefreshToken(ctx context.Context, id int32) error {
	return m.deactivateRefreshTokenFunc(ctx, id)
}

func (m *mockQuerier) GetRefreshToken(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
	return m.getRefreshTokenFunc(ctx, tokenHash)
}

func (m *mockQuerier) GetVerificationCode(ctx context.Context, verificationCodeHash string) (
	db.VerificationCode,
	error,
) {
	return m.getVerificationCodeFunc(ctx, verificationCodeHash)
}

func (m *mockQuerier) UpsertVerificationCode(
	ctx context.Context,
	arg db.UpsertVerificationCodeParams,
) (db.VerificationCode, error) {
	return m.upsertVerificationCodeFunc(ctx, arg)
}
