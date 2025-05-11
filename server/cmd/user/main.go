package main

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"os"
	"user"
	"user/db"
)

// main Initializes the user service
func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
		return
	}
	// tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	database, err := getDatabaseConnection()
	if err != nil {
		log.Fatalf("Error getting database connection: %v", err)
		return
	}

	queries := db.New(database)
	service := &user.ServiceImpl{
		Queries: *queries,
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	// TODO router.Use(jwtauth.Verifier(tokenAuth))
	// TODO router.Use(jwtauth.Authenticator(tokenAuth))
	router.Use(user.AuthMiddleware(*queries))

	router.Post("/user", user.CreateUserHandler(service))
	router.Get("/user/me", user.GetCurrentUserHandler(service))
	router.Get("/user", user.GetUserHandler(service))
	router.Get("/user/all", user.GetUsersHandler(service))
	router.Put("/user", user.UpdateUserHandler(service))
	router.Delete("/user/{id}", user.DeleteUserHandler(service))
}

// getDatabaseConnection Establishes a database connection and returns the database object
func getDatabaseConnection() (*sql.DB, error) {
	databaseDriver := os.Getenv("DATABASE_DRIVER")
	if databaseDriver == "" {
		return nil, errors.New("DATABASE_DRIVER environment variable not set")
	}

	databaseHost := os.Getenv("DATABASE_HOST")
	if databaseHost == "" {
		return nil, errors.New("DATABASE_HOST environment variable not set")
	}

	database, err := sql.Open(os.Getenv("DATABASE_DRIVER"), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	return database, nil
}
