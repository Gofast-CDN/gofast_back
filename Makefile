# Default values
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
TARGET := $(if $(ARGS),$(word 1,$(ARGS)),./...)

.PHONY: dev
dev:
	docker compose up

.PHONY: dev-build
dev-build:
	docker compose up --build

.PHONY: down
down:
	docker compose down -v

.PHONY: lint
lint:
	docker compose run --rm dev golangci-lint run

.PHONY: lint-fix
lint-fix:
	docker compose run --rm dev golangci-lint run --fix

.PHONY: test
test:
	docker compose run --rm dev go test -v $(TARGET)