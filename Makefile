PATHS = $(shell go list ./... | grep -v testintegration)

install_dependencies:
	go install gotest.tools/gotestsum@latest
	go install github.com/vektra/mockery/v2@v2.20.0
	go install github.com/ditrit/badaas-cli@latest

lint:
	golangci-lint run

test_unit:
	gotestsum --format pkgname $(PATHS)

test_integration:
	docker compose -f "docker/test_db/docker-compose.yml" up -d
	./docker/wait_for_db.sh
	gotestsum --format testname ./testintegration

test_e2e:
	docker compose -f "docker/test_db/docker-compose.yml" -f "docker/test_api/docker-compose.yml" up -d
	./docker/wait_for_api.sh 8000/info
	go test ./test_e2e -v

test_generate_mocks:
	mockery --all --keeptree

.PHONY: test_unit test_integration test_e2e

