# Docker compose commands out of the container
dev: up bash

prod:
	@docker compose up -d app-prod

up:
	@docker compose up -d app-dev

ps:
	@docker compose ps

logs:
	@docker compose logs app-dev --follow

build:
	@docker compose down && docker compose build --no-cache app-dev && docker compose up -d app-dev && docker compose exec app-dev bash

down:
	@docker compose down

stop:
	@docker compose stop

bash:
	@docker compose exec app-dev bash

# Dev commands to be used in the container
mysql:
	@mysql -h mysql-service -u default -p

test:
	@clear;go test -v ./...

cover:
	@go test -coverprofile=test/coverage.out ./... && go tool cover -html=test/coverage.out -o test/coverage.html && go run test/cover.go

sqlc:
	@sqlc generate -f sqlc.mysql.yaml

dbup:
	@mysql -h mysql-service -u default -p < databases/sql/mysql/schema/tasks_management.sql

dbdown:
	@mysql -h mysql-service -u default -p < databases/sql/down.sql

dbdump:
	@mysqldump -h mysql-service -u root -p dbdev > databases/sql/backup/$$(date +"%Y-%m-%-d").sql

air:
	@clear && air

swag:
	@swag init -g ./cmd/server/server.go -o api --parseDependency

# Generate JWT token with specified roles
jwt:
	@go run test/auth/generate_jwt.go --roles $(roles)

.PHONY: deploy-all deploy-app deploy-mysql deploy-rabbitmq clean

# Deploy all services
deploy-all: deploy-mysql deploy-rabbitmq deploy-app

# Deploy the main application
deploy-app:
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/app.yaml

# Deploy MySQL
deploy-mysql:
	kubectl apply -f k8s/mysql.yaml

# Deploy RabbitMQ
deploy-rabbitmq:
	kubectl apply -f k8s/rabbitmq.yaml

# Clean up all resources
clean:
	kubectl delete -f k8s/app.yaml
	kubectl delete -f k8s/mysql.yaml
	kubectl delete -f k8s/rabbitmq.yaml
	kubectl delete -f k8s/configmap.yaml
