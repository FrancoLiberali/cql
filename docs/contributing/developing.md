# Developing

This document provides the information you need to know before developing code for a pull request.

## Environment

- Install [go](https://go.dev/doc/install) >= v1.20
- Install project dependencies: `go get`
- Install [docker](https://docs.docker.com/engine/install/) and [compose plugin](https://docs.docker.com/compose/install/)

## Directory structure

This is the directory structure we use for the project:

- `docker/` : Contains the docker, docker-compose and configuration files for different environments.
- `docs/`: Contains the documentation showed for readthedocs.io.
- `orm/` *(Go code)*: Contains the code of the orm used by badaas.
- `test/`: Contains all the tests.

At the root of the project, you will find:

- The README.
- The changelog.
- The LICENSE file.

## Tests

### Dependencies

Running tests have some dependencies as: `gotestsum`, etc.. Install them with `make install_dependencies`.

### Linting

We use `golangci-lint` for linting our code. You can test it with `make lint`. The configuration file is in the default path (`.golangci.yml`). The file `.vscode.settings.json.template` is a template for your `.vscode/settings.json` that formats the code according to our configuration.

### Tests

We use the standard test suite in combination with [github.com/stretchr/testify](https://github.com/stretchr/testify) to do our testing. Tests have a database. Badaas-orm is tested on multiple databases. By default, the database used will be postgresql:

```sh
make test
```

To run the tests on another database you can use: `make test_postgresql`, `make test_cockroachdb`, `make test_mysql`, `make test_sqlite`, `make test_sqlserver`. All of them will be verified by our continuous integration system.

## Requirements

To be acceptable, contributions must:

- Have a good quality of code, based on <https://go.dev/doc/effective_go>.
- Have at least 80 percent new code coverage (although a higher percentage may be required depending on the importance of the feature). The tests that contribute to coverage are unit tests and integration tests.
- The features defined in the PR base issue must be explicitly tested by an e2e test or by integration tests in case it is not possible (for badaas-orm features for example).

## Use of Third-party code

Third-party code must include licenses.
