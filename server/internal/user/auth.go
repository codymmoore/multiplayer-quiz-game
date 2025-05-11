package user

import (
	"common"
	"context"
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
					ID:       int32(claims.ID),
					Username: claims.Username,
					Email:    claims.Email,
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
