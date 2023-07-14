all: gofmt vet build

gofmt:
	scripts/check_gofmt.sh

vet:
	go vet ./...

test:
	go test ./... -coverprofile=coverage.txt -covermode=atomic

build:
	go build ./...
