GO111MODULE=on

build:
	export GO111MODULE on; \
	go build ./...

build-generate:
	export GO111MODULE on; \
	go build -o hexa-go cmd/api/main.go

docker-build:
	docker build hexa-go:latest .

run:
	go run cmd/api/main.go
	
lint: build
	golint -set_exit_status ./...

test: lint
	go test ./... -v -covermode=count -coverprofile=coverage.out