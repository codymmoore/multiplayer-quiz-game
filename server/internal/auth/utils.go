package auth

import (
	"common"
	"crypto/sha256"
	"fmt"
	"github.com/keighl/postmark"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"time"
)

// GenerateJWT generates a JWT with the specified user info
func GenerateJWT(userId int, username string, email string) (string, error) {
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
		return "", err
	}

	signedToken, err := jwt.Sign(token, jwt.WithKey(common.JWTAlg, common.JWTSecret))
	if err != nil {
		return "", err
	}

	return string(signedToken), nil
}

// SecureHash creates a SHA256 hash for the specified
func SecureHash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

// SendEmail sends the specified email
func SendEmail(email *postmark.Email) error {
	client := postmark.NewClient("[SERVER-TOKEN]", "[ACCOUNT-TOKEN]")
	_, err := client.SendEmail(*email)
	if err != nil {
		return err
	}
	return nil
}
