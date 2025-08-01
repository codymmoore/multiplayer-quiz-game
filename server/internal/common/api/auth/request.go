package auth

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RenewRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type SendVerificationEmailRequest struct {
	UserId int    `json:"userId"`
	Email  string `json:"email"`
}

type VerifyEmailRequest struct {
	VerificationCode string `json:"verificationCode"`
}
