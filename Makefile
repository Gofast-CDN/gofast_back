# Default values
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
TARGET := $(if $(ARGS),$(word 1,$(ARGS)),./...)
DB ?= false

.PHONY: dev
dev:
	docker compose --profile dev up api mongodb

.PHONY: dev-build
dev-build:
	docker compose --profile dev up --build api mongodb

.PHONY: lint
lint:
	docker compose --profile dev run --rm --no-deps dev golangci-lint run

.PHONY: lint-fix
lint-fix:
	docker compose --profile dev run --rm --no-deps dev golangci-lint run --fix

.PHONY: test
test:
ifeq ($(DB),true)
	docker compose --profile test up -d mongodb
	docker compose --profile test run --rm dev go test -v ./tests/...
	docker compose --profile test down
else
	docker compose --profile test run --rm --no-deps dev go test -v ./tests/...
endif

# Add specific test targets
.PHONY: test-services
test-services:
	docker compose --profile test run --rm --no-deps dev go test -v ./tests/services

.PHONY: test-routes
test-routes:
	docker compose --profile test run --rm --no-deps dev go test -v ./tests/routes

.PHONY: test-auth
test-auth:
	docker compose --profile test run --rm --no-deps dev go test -v ./tests/utils/auth

# This prevents make from trying to process the arguments as targets
%:
	@: