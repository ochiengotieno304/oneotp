protoc:
	protoc -I ./pkg/proto \
  --go_out ./pkg/pb --go_opt paths=source_relative \
  --go-grpc_out ./pkg/pb --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./pkg/pb --grpc-gateway_opt paths=source_relative \
  ./pkg/proto/*.proto

serve:
	air -c .air.toml

docker:
	docker compose up

docker-build:
	docker compose up --build

docker-down:
	docker compose down

docker-watch:
	docker compose watch

migrate:
	go run ./pkg/db/migrate.go

test:
	go run ./tests/main.go