package auth

import (
	db "auth/db/generated"
	"common"
	api "common/auth"
	"common/user"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/keighl/postmark"
	"log"
	"net/http"
	"os"
	"time"
)

// RunServer starts the user service and listen for requests
func RunServer() {
	if err := common.InitJWT(); err != nil {
		log.Fatalf("failed to initialize JWT data: %v", err)
		return
	}

	database, err := common.GetDatabaseConnection()
	if err != nil {
		log.Fatalf("error establishing database connection: %v", err)
		return
	}

	baseUrl, err := common.GetBaseUrl()
	if err != nil {
		log.Fatalf("error getting base url: %v", err)
		return
	}

	userServiceUrl := os.Getenv("USER_SERVICE_URL")
	if userServiceUrl == "" {
		log.Fatal("USER_SERVICE_URL environment variable not set")
		return
	}
	userClient := &user.ClientImpl{
		BaseUrl:    userServiceUrl,
		HttpClient: http.DefaultClient,
	}

	postmarkServerToken := os.Getenv("POSTMARK_SERVER_TOKEN")
	if postmarkServerToken == "" {
		log.Fatal("POSTMARK_SERVER_TOKEN environment variable not set")
		return
	}
	postmarkAccountToken := os.Getenv("POSTMARK_ACCOUNT_TOKEN")
	if postmarkAccountToken == "" {
		log.Fatal("POSTMARK_ACCOUNT_TOKEN environment variable not set")
	}

	service := &ServiceImpl{
		Queries:        db.New(database),
		BaseUrl:        baseUrl,
		UserClient:     userClient,
		PostmarkClient: postmark.NewClient(postmarkServerToken, postmarkAccountToken),
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Minute))
	// r.Use(jwtauth.Verifier(common.TokenAuth))
	// r.Use(jwtauth.Authenticator(common.TokenAuth))
	r.Use(common.JWTMiddleware())

	r.Post(api.LoginEndpoint, LoginHandler(service))
	r.Post(api.LogoutEndpoint, LogoutHandler(service))
	r.Post(api.RenewEndpoint, RenewHandler(service))
	r.Post(api.SendVerificationEmailEndpoint, SendVerificationEmailHandler(service))
	r.Post(api.VerifyEmailEndpoint, VerifyEmailHandler(service))

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
		return
	}

	fmt.Println("Running on port " + port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
