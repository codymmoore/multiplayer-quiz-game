package user

import (
    "database/sql"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/jwtauth/v5"
    "log"
    "os"
    "user/db"
)

var tokenAuth *jwtauth.JWTAuth

func main() {
    conn, err := sql.Open(os.Getenv("DATABASE_DRIVER"), os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
        return
    }

    queries := db.New(conn)
    service := &ServiceImpl{
        queries: *queries,
    }

    router := chi.NewRouter()

    router.Use(middleware.Logger)
    // TODO router.Use(jwtauth.Verifier(tokenAuth))
    // TODO router.Use(jwtauth.Authenticator(tokenAuth))
    router.Use(AuthMiddleware(*queries))

    router.Post("/user", CreateUserHandler(service))
    router.Get("/user/me", GetCurrentUserHandler(service))
    router.Get("/user", GetUserHandler(service))
    router.Get("/user/all", GetUsersHandler(service))
    router.Put("/user", UpdateUserHandler(service))
    router.Delete("/user/{id}", DeleteUserHandler(service))
}
