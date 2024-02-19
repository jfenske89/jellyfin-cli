PROJECT_NAME=jellyfin-cli

all: install fmt vet lint test compile

install:
	go mod tidy

fmt:
	goimports -w --local codeberg.org/jfenske ./

vet:
	go vet ./...

test:
	go test -v ./...

compile:
	go build -o ./bin/jellyfin-cli internal/cmd/jellyfin-cli/main.go

lint:
	golangci-lint run --verbose

upgrade:
	go get -u ./...
	go mod tidy
