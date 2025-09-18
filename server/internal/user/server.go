package user

import (
	"common"
	api "common/user"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/lib/pq" // registers "postgres" driver
	"log"
	"net/http"
	"os"
	"time"
	"user/db/generated"
)

// RunServer Start the user service and listen for requests
func RunServer() {
	if err := common.InitJWT(); err != nil {
		log.Fatal("Error initializing JWT")
		return
	}

	database, err := common.GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Error establishing database connection: %v", err)
		return
	}

	queries := db.New(database)
	service := &ServiceImpl{
		Queries: queries,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Minute))

	r.Post("/user", CreateUserHandler(service))
	r.Get("/user", GetUserHandler(service))
	r.Get("/user/all", GetUsersHandler(service))
	r.Patch(fmt.Sprintf(api.VerifyUserEndpoint, "{id}"), VerifyUserHandler(service))
	r.Group(
		func(r chi.Router) {
			r.Use(jwtauth.Verifier(common.TokenAuth))
			r.Use(jwtauth.Authenticator(common.TokenAuth))
			r.Use(common.AuthMiddleware())

			r.Get("/user/me", GetCurrentUserHandler(service))
			r.Patch("/user/{id}", UpdateUserHandler(service))
			r.Delete("/user/{id}", DeleteUserHandler(service))
		},
	)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
		return
	}

	fmt.Println("Listening on port " + port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
