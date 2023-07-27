run:
	go run ./cmd/producer/main.go
run-consumer:
	go run ./cmd/consumer/main.go

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

test:
	curl -X POST -H "Content-Type: application/json" -d '{"name": "test", "email": "test6@gamil.com", "password": "1f23"}' localhost:8000/api/register


py-run:
	python -m kafka_py.run