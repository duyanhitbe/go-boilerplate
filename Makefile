include .env

build:
	@go build -o bin/golang-boilerplate cmd/api/main.go
run: build
	@./bin/golang-boilerplate
unit-test:
	@godotenv -f .env.test go test -v ./...
sqlc:
	@sqlc generate
migrate:
	@migrate create -dir internal/database/migrations -seq -ext sql ${name}
migrate-up:
	@migrate -path ./internal/database/migrations -database "${DB_URL}" -verbose up
migrate-down:
	@migrate -path ./internal/database/migrations -database "${DB_URL}" -verbose down
migrate-force:
	@migrate -path ./internal/database/migrations -database "${DB_URL}" force ${step}
mock:
	@mockgen -package mockdb -destination internal/database/mock/store.go github.com/duyanhitbe/go-boilerplate/internal/database/generated Store
	@mockgen -package mockhash -destination internal/hash/mock/hash.go github.com/duyanhitbe/go-boilerplate/internal/hash Hash
	@mockgen -package mocktoken -destination internal/token/mock/token.go github.com/duyanhitbe/go-boilerplate/internal/token Token