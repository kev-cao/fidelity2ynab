GOBIN ?= $$(pwd)/bin/

.PHONY: dev
dev: 
	docker compose --profile dev up

.PHONY: build-dev
build-dev:
	docker compose --profile dev build

.PHONY: prod
prod:
	docker compose --profile prod up

.PHONY: build-prod
build-prod:
	docker compose --profile prod build

.PHONY: test
test:
	go test -v ./...

bin/sync: $(shell find pkg -type f) $(shell find cmd -type f)
	GOBIN=$(GOBIN) go install ./cmd/sync
