package common

import (
	"github.com/lestrrat-go/jwx/v2/jwa"
)

const (
	UserClaimsCtxKey  = "user"
	UserIdClaimsKey   = "user_id"
	UsernameClaimsKey = "username"
	EmailClaimsKey    = "email"
	JWTCtxKey         = "JWT"
	JWTAlg            = jwa.HS256
)
