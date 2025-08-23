package common

// UserClaims Stores JWT information for User
type UserClaims struct {
	ID       int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
