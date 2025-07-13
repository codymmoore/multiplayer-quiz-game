module server

go 1.24

replace user => ./internal/user

require github.com/go-chi/chi/v5 v5.2.1

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
)
