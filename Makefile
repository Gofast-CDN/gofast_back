.PHONY: lint
lint:
	docker compose run --rm dev golangci-lint run

.PHONY: lint-fix
lint-fix:
	docker compose run --rm dev golangci-lint run --fix

.PHONY: test
test:
	docker compose run --rm dev go test -v ./...
