package common

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/lib/pq" // registers "postgres" driver
	"net/http"
	"os"
	"strings"
)

// InitJWT Initialize the global JWTAuth instance using the JWT_SECRET environment variable
func InitJWT() error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return errors.New("JWT_SECRET environment variable not set")
	}
	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
	return nil
}

// GetDatabaseConnection Establishes a database connection and returns the database object
func GetDatabaseConnection() (*sql.DB, error) {
	databaseDriver := os.Getenv("DATABASE_DRIVER")
	if databaseDriver == "" {
		return nil, errors.New("DATABASE_DRIVER environment variable not set")
	}

	databaseHost := os.Getenv("DATABASE_URL")
	if databaseHost == "" {
		return nil, errors.New("DATABASE_URL environment variable not set")
	}

	database, err := sql.Open(databaseDriver, databaseHost)
	if err != nil {
		return nil, err
	}

	if err = database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}

// GetBaseUrl Get the base URL for the current service
func GetBaseUrl() (string, error) {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		return "", errors.New("BASE_URL environment variable not set")
	}
	return baseUrl, nil
}

// GetRouteUrl Get the base URL + route pattern for specified context
func GetRouteUrl(context context.Context) (string, error) {
	baseUrl, err := GetBaseUrl()
	if err != nil {
		return "", err
	}
	return baseUrl + GetRoutePattern(context), nil
}

// GetRoutePattern Get the route from the specified context
func GetRoutePattern(context context.Context) string {
	routeContext := chi.RouteContext(context)
	return routeContext.RoutePattern()
}

// GetJWT Gets the JWT string from the request header
func GetJWT(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)
	return tokenString
}
