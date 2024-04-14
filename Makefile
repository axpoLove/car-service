generate:
	go generate ./...

docker:
	docker compose up -d

run-car-info-api:
	go run cmd/car-info-api/main.go

run-car-api:
	go run cmd/console/main.go up
	go run cmd/car-api/main.go


