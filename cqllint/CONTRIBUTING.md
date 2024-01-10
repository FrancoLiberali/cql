# Contribute to the development of cqllint <!-- omit in toc -->

- [Local compilation](#local-compilation)
- [Tests](#tests)

## Local compilation

You can make modifications to the cqllint source code and compile it locally with:

```bash
go build .
```

You can then run the cqllint executable directly or add it to your $GOPATH to run it from a project:

```bash
go install .
```

## Tests

We use the standard test suite in combination with [github.com/stretchr/testify](https://github.com/stretchr/testify) to do our unit testing.

To run them, please run:

```sh
go test ./...
```
