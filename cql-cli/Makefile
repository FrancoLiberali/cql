test: clean_test_results
	go test ./tests -v

clean_test_results:
	rm -f tests/*_conditions.go
	rm -f tests/**/cql.go

.PHONY: test