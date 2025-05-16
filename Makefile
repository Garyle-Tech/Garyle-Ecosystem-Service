migrate-up:
	migrate -path ./migrations -database "postgres://postgres:admin@localhost:5432/garyle_ecosystem_db?sslmode=disable" up 

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:admin@localhost:5432/garyle_ecosystem_db?sslmode=disable" down 


# use it with "make migrate-create name=migration_name"
migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name) 

run:
	go run cmd/api/main.go