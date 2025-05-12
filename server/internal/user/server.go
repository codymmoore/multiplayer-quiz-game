package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // registers "postgres" driver
	"log"
	"net/http"
	"os"
	"user/db"
)

var baseUrl string

// RunServer Start the user service and listen for requests
func RunServer() {
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

	baseUrl, err = getBaseUrl()
	if err != nil {
		log.Fatalf("Error getting base url: %v", err)
		return
	}

	queries := db.New(database)
	service := &ServiceImpl{
		Queries: *queries,
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	// TODO router.Use(jwtauth.Verifier(tokenAuth))
	// TODO router.Use(jwtauth.Authenticator(tokenAuth))
	// TODO router.Use(user.authMiddleware(*queries))

	router.Post("/user", CreateUserHandler(service))
	router.Get("/user/me", GetCurrentUserHandler(service))
	router.Get("/user", GetUserHandler(service))
	router.Get("/user/all", GetUsersHandler(service))
	router.Put("/user/{id}", UpdateUserHandler(service))
	router.Delete("/user/{id}", DeleteUserHandler(service))

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

// getRouteUrl Get the base URL + route pattern for specified context
func getRouteUrl(context context.Context) string {
	return baseUrl + getRoutePattern(context)
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

// getRoutePattern Get the route from the specified context
func getRoutePattern(context context.Context) string {
	routeContext := chi.RouteContext(context)
	return routeContext.RoutePattern()
}
