lint:
	golangci-lint run

test_unit:
	go test ./... -v

test_e2e:
	docker compose -f "docker/test_db/docker-compose.yml" -f "docker/test_api/docker-compose.yml" up -d
	./docker/wait_for_api.sh 8000/info
	go test ./test_e2e -v

test_generate_mocks:
	mockery --all --keeptree

.PHONY: test_unit test_e2e

