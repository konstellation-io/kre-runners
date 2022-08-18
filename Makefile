.DEFAULT_GOAL := help

# AutoDoc
# -------------------------------------------------------------------------
.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help


.PHONY: build-py
build-py: ## Build Python runner docker image.
	docker build -t konstellation/kre-py:latest kre-py

.PHONY: build-go
build-go: ## Build GoLang runner docker image.
	docker build -t konstellation/kre-go:latest kre-go

.PHONY: build-entrypoint
build-entrypoint: ## Build Entrypoint docker image.
	docker build -t konstellation/kre-entrypoint:latest kre-entrypoint

.PHONY: build
build: build-py build-go build-entrypoint ## Build all docker images.
	echo "done"

.PHONY: protos
protos: ## Generate proto files
	./generate_protobuf_code.sh

		