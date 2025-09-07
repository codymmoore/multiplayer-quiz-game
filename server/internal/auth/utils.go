package auth

import (
	"common"
	"crypto/sha256"
	"fmt"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"time"
)

// GenerateJWTToken generates a jwt.Token with the specified user info
func GenerateJWTToken(userId int, username string, email string) (jwt.Token, error) {
	token, err := jwt.NewBuilder().
		Issuer("quizchief-auth").
		Subject(username).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(30*time.Minute)).
		Claim(common.UserIdClaimsKey, userId).
		Claim(common.UsernameClaimsKey, username).
		Claim(common.EmailClaimsKey, email).
		Build()
	if err != nil {
		return nil, err
	}
	return token, nil
}

// SecureHash creates a SHA256 hash for the specified
func SecureHash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
