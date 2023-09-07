PATHS = $(shell go list ./... | grep -v testintegration)

install_dependencies:
	go install gotest.tools/gotestsum@latest
	go install github.com/vektra/mockery/v2@v2.20.0
	go install github.com/ditrit/badaas-cli@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run
	cd test_e2e && golangci-lint run --config ../.golangci.yml

test_unit:
	gotestsum --format pkgname $(PATHS)

rmdb:
	docker stop badaas-test-db && docker rm badaas-test-db

postgresql:
	docker compose -f "docker/postgresql/docker-compose.yml" up -d

cockroachdb:
	docker compose -f "docker/cockroachdb/docker-compose.yml" up -d

mysql:
	docker compose -f "docker/mysql/docker-compose.yml" up -d

sqlserver:
	docker compose -f "docker/sqlserver/docker-compose.yml" up -d --build

test_integration_postgresql: postgresql
	DB=postgres gotestsum --format testname ./testintegration

test_integration_cockroachdb: cockroachdb
	DB=postgres gotestsum --format testname ./testintegration -tags=cockroachdb

test_integration_mysql: mysql
	DB=mysql gotestsum --format testname ./testintegration -tags=mysql

test_integration_sqlite:
	DB=sqlite gotestsum --format testname ./testintegration

test_integration_sqlserver: sqlserver
	DB=sqlserver gotestsum --format testname ./testintegration

test_integration: test_integration_postgresql

test_e2e:
	docker compose -f "docker/cockroachdb/docker-compose.yml" -f "docker/test_api/docker-compose.yml" up -d
	./docker/wait_for_api.sh 8000/info
	go test ./test_e2e -v

test_generate_mocks:
	mockery --all --keeptree

.PHONY: test_unit test_integration test_e2e

