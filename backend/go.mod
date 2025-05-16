module email-marketing-service

go 1.23.0

toolchain go1.23.9

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/shopspring/decimal v1.4.0
	github.com/sqlc-dev/pqtype v0.3.0
	golang.org/x/oauth2 v0.30.0
)

require cloud.google.com/go/compute/metadata v0.3.0 // indirect
