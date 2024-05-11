DEV_VOLUMES = --mount type=bind,src=$$(pwd)/.,dst=/app \
	--mount type=volume,dst=/app/bin
VOLUMES = --mount type=bind,src=$$(pwd)/secrets.yaml,dst=/run/secrets/secrets

GOBIN ?= $$(pwd)/bin/

.PHONY: dev
dev: 
	docker run $(VOLUMES) $(DEV_VOLUMES) -e ENV=dev --entrypoint=/bin/sh \
	fidelity2ynab -c "make bin/sync && ./bin/sync"

.PHONY: prod
prod:
	make build && \
	docker run $(VOLUMES) fidelity2ynab

.PHONY: build
build:
	docker build -t fidelity2ynab .

.PHONY: test
test:
	docker run $(VOLUMES) $(DEV_VOLUMES) --entrypoint=/bin/sh fidelity2ynab -c \
	"go test -v ./..."

.PHONY: test-all
test-all:
	docker run $(VOLUMES) $(DEV_VOLUMES) --entrypoint=/bin/sh fidelity2ynab -c \
	"go test -tags=browser -v ./..."

bin/sync: $(filter-out %_test.go, $(shell find pkg -type f) $(shell find cmd -type f))
	GOBIN=$(GOBIN) go install ./cmd/sync
