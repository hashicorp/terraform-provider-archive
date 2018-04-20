default: test

test: lint
	TF_SCHEMA_PANIC_ON_ERROR=1 go test -v -cover $(TESTARGS) -timeout=30s ./...

testacc: lint
	TF_SCHEMA_PANIC_ON_ERROR=1 TF_ACC=1 go test -v -cover $(TESTARGS) -timeout 120m ./...

lint:
	@bash scripts/lint.sh

vendor-status:
	@govendor status

.PHONY: build test testacc lint vendor-status
