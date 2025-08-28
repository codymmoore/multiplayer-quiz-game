package user

import (
	"common"
	api "common/api/user"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(time.Minute))

	router.Post("/user", CreateUserHandler(service))
	router.Group(
		func(router chi.Router) {
			// TODO r.Use(jwtauth.Verifier(tokenAuth))
			// TODO r.Use(jwtauth.Authenticator(tokenAuth))
			// TODO r.Use(user.authMiddleware(*queries))

			router.Get("/user/me", GetCurrentUserHandler(service))
			router.Get("/user", GetUserHandler(service))
			router.Get("/user/all", GetUsersHandler(service))
			router.Patch("/user/{id}", UpdateUserHandler(service))
			router.Delete("/user/{id}", DeleteUserHandler(service))
			router.Patch(fmt.Sprintf(api.VerifyUserEndpoint, "{id}"), VerifyUserHandler(service))
		},
	)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
		return
	}

	fmt.Println("Listening on port " + port)

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
