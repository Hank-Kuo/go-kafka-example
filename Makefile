run-producer:
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
	curl -X POST -H "Content-Type: application/json" -d '{"email": "hank_kuo@trendmicro.com", "password": "1f23"}' localhost:8000/api/login
test1:
	curl -X POST -H "Content-Type: application/json" -d '{"name": "xxx", "email": "hank_kuo@trendmicro2.com", "password": "1f23"}' localhost:8000/api/register
test2:
	curl -X POST -H "Content-Type: application/json" -d '{"email": "hank_kuo@trendmicro2.com", "password": "1f23"}' localhost:8000/api/login
py-consumer:
	poetry run python -m kafka_py.consumer
py-producer:
	poetry run python -m kafka_py.producer
poetry-export:
	poetry export -f requirements.txt --without-hashes > requirements.txt