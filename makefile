VERSION=$(shell git describe --tags 2>/dev/null || echo "snapshot")

.PHONY: build
build: bin/blast-validator

.PHONY: release
release: bin/blast-validator.$(VERSION).tar.gz
	@rm bin/blast-validator

.PHONY: gen-docs
gen-docs: docs/index.html

docs/index.html: openapi.yml node_modules/.bin/redoc-cli
	node_modules/.bin/redoc-cli bundle -o docs/index.html openapi.yml

bin/blast-validator.$(VERSION).tar.gz: bin/blast-validator
	@cd bin && tar -czf blast-validator.$(VERSION).tar.gz blast-validator

bin/blast-validator: $(shell find v1 -type f -name '*.go')
	CGO_ENABLED=0 go build -o bin/blast-validator ./v1/cmd/server/main.go

node_modules/.bin/redoc-cli:
	npm install