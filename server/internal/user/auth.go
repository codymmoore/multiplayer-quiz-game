package user

import (
	"common"
	"context"
	"database/sql"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"user/db"
)

const USER_CLAIMS_KEY string = "user"

// AuthMiddleware Validates JWT and extracts claims
func AuthMiddleware(queries db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				_, claimsMap, err := jwtauth.FromContext(ctx)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
				}

				claims := common.UserClaims{
					ID:       claimsMap["user_id"].(int),
					Username: claimsMap["username"].(string),
					Email:    claimsMap["email"].(string),
				}

				params := db.GetUserParams{
					ID:       sql.NullInt32{Int32: int32(claims.ID), Valid: true},
					Username: sql.NullString{String: claims.Username, Valid: true},
					Email:    sql.NullString{String: claims.Email, Valid: true},
				}
				if _, err = queries.GetUser(ctx, params); err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
				}

				ctx = context.WithValue(ctx, USER_CLAIMS_KEY, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}
