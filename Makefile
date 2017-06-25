.PHONY: install build test lint vet

build:
	@CGO_ENABLED=0 go build -o ./main -a -ldflags '-s' -installsuffix cgo main.go

install:
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/stretchr/testify
	@go get -v ./

test:
	@go test -i
	@go test -race -v `go list ./... | grep -v /vendor/`
	@go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status
	@go vet `go list ./... | grep -v /vendor/`