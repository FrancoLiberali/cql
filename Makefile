install_dependencies:
	go install gotest.tools/gotestsum@latest
	go install github.com/FrancoLiberali/cql/cql-gen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run
	cd cql-gen && golangci-lint run --config ../.golangci.yml

rmdb:
	docker stop cql-test-db && docker rm cql-test-db

postgresql:
	docker compose -f "docker/postgresql/docker-compose.yml" up -d

cockroachdb:
	docker compose -f "docker/cockroachdb/docker-compose.yml" up -d

mysql:
	docker compose -f "docker/mysql/docker-compose.yml" up -d

sqlserver:
	docker compose -f "docker/sqlserver/docker-compose.yml" up -d --build

test_postgresql: postgresql
	DB=postgres gotestsum --format testname ./...

test_cockroachdb: cockroachdb
	DB=postgres gotestsum --format testname ./... -tags=cockroachdb

test_mysql: mysql
	DB=mysql gotestsum --format testname ./... -tags=mysql

test_sqlite:
	DB=sqlite gotestsum --format testname ./...

test_sqlserver: sqlserver
	DB=sqlserver gotestsum --format testname ./...

test: test_sqlite

.PHONY: test test_postgresql test_sqlserver test_sqlite test_mysql mysql postgresql cockroachdb sqlserver

