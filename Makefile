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
	@test -f .env.docker || (echo "Create .env.docker from .env.docker.example first" && exit 1)
	@set -a && . ./.env.docker && set +a && \
	docker-compose exec mysql mysqldump -u "$$DB_USER" -p"$$DB_PASSWORD" "$$DB_NAME" > backup.sql

docker-db-restore:
	@test -f .env.docker || (echo "Create .env.docker from .env.docker.example first" && exit 1)
	@set -a && . ./.env.docker && set +a && \
	docker-compose exec -T mysql mysql -u "$$DB_USER" -p"$$DB_PASSWORD" "$$DB_NAME" < backup.sql

run:
	go run cmd/api/main.go
	
lint: build
	golint -set_exit_status ./...
	golangci-lint run ./...

test: lint
	go test ./... -v -covermode=count -coverprofile=coverage.out