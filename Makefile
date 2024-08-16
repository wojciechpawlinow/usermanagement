.PHONY: run test help default

default: help

help:
	@echo 'Available commands:'
	@echo
	@echo 'Usage:'
	@echo '    make build       compile the project and build docker image'
	@echo '    make run         run application inside container'
	@echo '    make test        run tests on a compiled project'
	@echo '    make test-cc     run tests with coverage report on a compiled project'
	@echo '    make lint        source code linting'
	@echo '    make reset-db    down and up migrations'
	@echo

build:
	@echo "\n! ensure creating config.yaml from a dist file before building the image if you're overriding envs\n"
	@docker-compose build --no-cache app

run:
	@docker-compose up -d app && docker logs -f users_app

test:
	DB_READ_HOST=localhost DB_WRITE_HOST=localhost go test -race ./...

test-cc:
	DB_READ_HOST=localhost DB_WRITE_HOST=localhost go test -race -failfast -coverprofile cover.out ./...
	go tool cover -html=cover.out

lint:
	go mod tidy
	go vet ./...
	go fmt ./...
	goimports -w -local github.com/wojciechpawlinow .

DB_USER ?= user
DB_PASSWORD ?= pass
DB_HOST ?= mysql
DB_PORT ?= 3306
DB_NAME ?= users
reset-db:
	@docker exec -it users_app migrate -path /app/migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" down
	@docker exec -it users_app migrate -path /app/migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" up