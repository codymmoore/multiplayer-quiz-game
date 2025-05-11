module user

go 1.24

replace user/dto => ./dto

replace common => ../common

require (
	common v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-chi/jwtauth/v5 v5.3.3
	golang.org/x/crypto v0.38.0
	user/dto v0.0.0-00010101000000-000000000000
)

require (
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.6 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx/v2 v2.1.3 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)
