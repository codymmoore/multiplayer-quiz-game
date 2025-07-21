package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"user/db/generated"
)

// AuthMiddleware Validates JWT and extracts claims
func AuthMiddleware(queries db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				_, claimsMap, err := jwtauth.FromContext(ctx)
				if err != nil {
					http.Error(
						w,
						fmt.Sprintf("unable to get claims map from context: %v", err),
						http.StatusUnauthorized,
					)
					return
				}

				claims := UserClaims{
					ID:       claimsMap["user_id"].(int),
					Username: claimsMap["username"].(string),
					Email:    claimsMap["email"].(string),
				}

				statusCode, user, err := getUser(ctx, claims.ID)
				if err != nil {
					http.Error(w, fmt.Sprintf("unable to get user: %v", err), http.StatusInternalServerError)
					return
				} else if statusCode != http.StatusOK {
					http.Error(w, "invalid user", http.StatusUnauthorized)
					return
				}

				if user.Username != claims.Username {
					http.Error(w, "invalid username", http.StatusUnauthorized)
					return
				} else if user.Email != claims.Email {
					http.Error(w, "invalid email", http.StatusUnauthorized)
					return
				}

				ctx = context.WithValue(ctx, UsersClaimKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

// getUser Retrieve the user with the specified ID
func getUser(ctx context.Context, userId int) (int, *getUserResponse, error) {
	getUserUrl, err := GetBaseUrl()
	if err != nil {
		return -1, nil, fmt.Errorf("unable to get base url: %w", err)
	}

	getUserResp, err := http.Get(fmt.Sprintf("%s/user?id=%d", getUserUrl, userId))
	if err != nil {
		return -1, nil, fmt.Errorf("unable to get get user response: %w", err)
	}
	defer getUserResp.Body.Close()

	statusCode := getUserResp.StatusCode
	if statusCode != http.StatusOK {
		return statusCode, nil, nil
	}

	var getUserResponse getUserResponse
	if err := json.NewDecoder(getUserResp.Body).Decode(&getUserResponse); err != nil {
		return statusCode, nil, fmt.Errorf("unable to parse get user response: %w", err)
	}

	return statusCode, &getUserResponse, nil
}
