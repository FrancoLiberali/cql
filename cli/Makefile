lint:
	golangci-lint run

test_unit: clean_test_unit_results
	go test ./... -v

clean_test_unit_results:
	rm -f cmd/gen/conditions/*_conditions.go
	rm -f cmd/gen/conditions/tests/**/badaas-orm.go

.PHONY: test_unit