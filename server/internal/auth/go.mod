module auth

go 1.24

replace common => ../common

require (
	common v0.0.0-00010101000000-000000000000
	github.com/keighl/postmark v0.0.0-20190821160221-28358b1a94e3
	github.com/lestrrat-go/jwx/v2 v2.1.3
	golang.org/x/crypto v0.31.0
)

require (
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/go-chi/chi/v5 v5.2.1 // indirect
	github.com/go-chi/jwtauth/v5 v5.3.3 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.6 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	goji.io v2.0.2+incompatible // indirect
	golang.org/x/sys v0.28.0 // indirect
)
