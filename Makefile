run:
	go run ./cmd/main.go
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
