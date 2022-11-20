run:
	docker-compose down --remove-orphans -v
	docker-compose build
	docker-compose up -d db
	sleep 10
	docker-compose run --rm migrate
	docker-compose up -d app

lint:
	which golangci-lint || go install github.com/golangci/golangci-lint/...@v1.50.1
	golangci-lint run -v ./...