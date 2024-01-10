# Contribute to the development of cql-gen <!-- omit in toc -->

- [Local compilation](#local-compilation)
- [Tests](#tests)

## Local compilation

You can make modifications to the cql-gen source code and compile it locally with:

```bash
go build .
```

You can then run the cql-gen executable directly or add a link in your $GOPATH to run it from a project:

```bash
go install .
```

## Tests

We use the standard test suite in combination with [github.com/stretchr/testify](https://github.com/stretchr/testify) to do our unit testing.

To run them, please run:

```sh
make test
```
