NAME := report-creator

PKG := `go list -mod=mod -f {{.Dir}} ./...`
MAIN := app/cmd/main.go

ifeq ($(RACE),1)
	GOFLAGS=-race
endif

run:
	@docker-compose build && docker-compose up app
test:
	@docker-compose build && docker-compose up test


install-lint: install-goimports install-golangci install-looppointer


install-goimports:
	@go install golang.org/x/tools/cmd/goimports@latest

install-golangci:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

install-looppointer:
	@go install github.com/kyoh86/looppointer/cmd/looppointer@latest

install-swagger:
	@go install github.com/swaggo/swag/cmd/swag

.PHONY: build
build:
	@echo build $(VERSION)
	@CGO_ENABLED=0 \
	GOPRIVATE=$(GOPRIVATE) \
		go build \
		-mod=mod \
		$(LDFLAGS) \
		$(GOFLAGS) \
		-o ${NAME} \
		$(MAIN)

swag:
	@swag init -g ./app/cmd/main.go

fmt:
	@goimports -local ${NAME} -l -w $(PKG)

lint:
	@golangci-lint run -c .golangci.yml
	@looppointer -c 3 ./...

mod-download:
	go mod download all

mod-tidy:
	go mod tidy

mod: mod-tidy mod-download install-looppointer install-swagger

pre-commit: fmt lint test
