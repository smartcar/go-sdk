
all: gofmt vet lint build

gofmt:
	scripts/check_gofmt.sh

lint:
	golint ./...

vet:
	go vet ./...

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic

build:
	go build ./...
