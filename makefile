include .env

migration:
	go run cmd/migration/main.go

seed:
	go run cmd/seed/main.go

run.app:
	go build -o tmp/main cmd/app/main.go
	./tmp/main


install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Proto
.PHONY: inject-tag
inject-tag:
	protoc-go-inject-tag -input="mtgrpc/*.pb.go"

.PHONY: generate-proto
generate-proto:
	protoc \
		--proto_path=. \
		--go_out=. \
		--go_grpc_out=. \
		proto/*.proto

.PHONY: proto
proto: generate-proto
# proto: generate-proto inject-tag

migrate-up:
	migrate -path db/migrations -database "postgres://${POSTGRESQL_USERNAME}:${POSTGRESQL_PASSWORD}@localhost:${POSTGRESQL_PORT}/${POSTGRESQL_DATABASE}?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgres://${POSTGRESQL_USERNAME}:${POSTGRESQL_PASSWORD}@localhost:${POSTGRESQL_PORT}/${POSTGRESQL_DATABASE}?sslmode=disable" down
