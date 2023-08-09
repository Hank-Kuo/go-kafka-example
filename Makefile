run-producer:
	go run ./cmd/producer/main.go
run-consumer:
	poetry run python -m kafka_py.consumer
py-consumer:
	poetry run python -m kafka_py.consumer

pb-complier:
	@echo Compiling $$APP proto...
	@protoc --go_out=./ --go-grpc_out=require_unimplemented_servers=false:. ./pb/$$APP/*.proto
	
local:
	docker-compose -f docker-compose.local.yaml up -d

docker-up:
	docker-compose up -d
docker-down:
	docker-compose down
docker-restart:
	docker-compose restart

migrate:
	migrate create -ext sql -dir migrations -seq ${NAME}
migrate-up:	
	migrate --verbose -database "postgres://postgres:postgres@localhost:5432/db?sslmode=disable" -path migrations up
migrate-down:	
	migrate --verbose -database "postgres://postgres:postgres@localhost:5432/db?sslmode=disable" -path migrations down
migrate-force:	
	migrate --verbose -database "postgres://postgres:postgres@localhost:5432/db?sslmode=disable" -path migrations force ${V}

poetry-export:
	poetry export -f requirements.txt --without-hashes > requirements.txt


linter:
	@echo Starting linters
	golangci-lint run ./...