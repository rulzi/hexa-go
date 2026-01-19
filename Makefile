GO111MODULE=on

build:
	export GO111MODULE on; \
	go build ./...

build-generate:
	export GO111MODULE on; \
	go build -o hexa-go cmd/api/main.go

docker-build:
	docker build -t hexa-go:latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-restart:
	docker-compose restart

docker-ps:
	docker-compose ps

docker-clean:
	docker-compose down -v
	docker system prune -f

docker-db-backup:
	docker-compose exec mysql mysqldump -u hexa_user -phexapassword123 hexa_go > backup.sql

docker-db-restore:
	docker-compose exec -T mysql mysql -u hexa_user -phexapassword123 hexa_go < backup.sql

run:
	go run cmd/api/main.go
	
lint: build
	golint -set_exit_status ./...
	golangci-lint run ./...

test: lint
	go test ./... -v -covermode=count -coverprofile=coverage.out