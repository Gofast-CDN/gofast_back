.PHONY: lint
lint:
	docker compose run --rm dev golangci-lint run

.PHONY: lint-fix
lint-fix:
	docker compose run --rm dev golangci-lint run --fix