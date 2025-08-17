package common

import (
    "common/api/user"
    "context"
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

                baseUrl, err := GetBaseUrl()
                if err != nil {
                    http.Error(w, fmt.Sprintf("unable to get base url: %v", err), http.StatusInternalServerError)
                }

                userClient := user.ClientImpl{
                    BaseUrl:    baseUrl,
                    HttpClient: http.DefaultClient,
                }

                getUserRequest := &user.GetUserRequest{
                    UserId:   &claims.ID,
                    Username: &claims.Username,
                    Email:    &claims.Email,
                }
                if _, err := userClient.GetUser(getUserRequest, GetJWT(r)); err != nil {
                    http.Error(w, fmt.Sprintf("unable to get user: %v", err), http.StatusInternalServerError)
                    return
                }

                ctx = context.WithValue(ctx, UsersClaimKey, claims)
                next.ServeHTTP(w, r.WithContext(ctx))
            },
        )
    }
}
