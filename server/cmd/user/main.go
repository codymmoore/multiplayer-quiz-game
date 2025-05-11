package main

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // registers "postgres" driver
	"log"
	"net/http"
	"os"
	"user"
	"user/db"
)

// main Initializes the user service
func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
		return
	}

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

	baseUrl, err := getBaseUrl()
	if err != nil {
		log.Fatalf("Error getting base url: %v", err)
		return
	}

	queries := db.New(database)
	service := &user.ServiceImpl{
		Queries: *queries,
		BaseUrl: baseUrl,
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	// TODO router.Use(jwtauth.Verifier(tokenAuth))
	// TODO router.Use(jwtauth.Authenticator(tokenAuth))
	// router.Use(user.AuthMiddleware(*queries))

	router.Post("/user", user.CreateUserHandler(service))
	router.Get("/user/me", user.GetCurrentUserHandler(service))
	router.Get("/user", user.GetUserHandler(service))
	router.Get("/user/all", user.GetUsersHandler(service))
	router.Put("/user/{id}", user.UpdateUserHandler(service))
	router.Delete("/user/{id}", user.DeleteUserHandler(service))

	portStr := ":" + os.Getenv("PORT")
	log.Println("Listening on port " + portStr)

	err = http.ListenAndServe(portStr, router)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
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

	database, err := sql.Open(databaseDriver, databaseHost)
	if err != nil {
		return nil, err
	}

	return database, nil
}

// getBaseUrl Generates the base URL using the host and port specified in the environment file
func getBaseUrl() (string, error) {
	host := os.Getenv("HOST")
	if host == "" {
		return "", errors.New("HOST environment variable not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return "", errors.New("PORT environment variable not set")
	}

	return host + ":" + port, nil
}
