test: clean_test_results
	go test ./... -v

clean_test_results:
	rm -f cmd/gen/conditions/*_conditions.go
	rm -f cmd/gen/conditions/tests/**/cql.go

.PHONY: test