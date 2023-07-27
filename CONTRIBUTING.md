# Contribute to the development of badaas

- [Tests](#tests)
  - [Unit tests](#unit-tests)
  - [Feature tests (of end to end tests)](#feature-tests-of-end-to-end-tests)
- [Logger](#logger)
- [Directory structure](#directory-structure)
- [Git](#git)
  - [Branch naming policy](#branch-naming-policy)
  - [Default branch](#default-branch)
  - [How to release](#how-to-release)

## Tests

### Unit tests

We use the standard test suite in combination with [github.com/stretchr/testify](https://github.com/stretchr/testify) to do our unit testing. Mocks are generated using [mockery](https://github.com/vektra/mockery) a mock generator using this command `mockery --all --keeptree`.

To run them, please run:

```sh
go test $(go list ./... | sed 1d) -v
```

### Feature tests (of end to end tests)

We use docker to run a Badaas instance in combination with one node of CockroachDB.

Run:

```sh
docker compose -f "scripts/e2e/docker-compose.yml" up -d --build
```

Then in an another shell:

```sh
go test -v
```

The feature files can be found in the `feature` folder.

## Logger

We use ubber's [zap](https://pkg.go.dev/go.uber.org/zap) to log stuff, please take `zap.Logger` as an argument for your services constructors. [fx](https://github.com/uber-go/fx) will provide your service with an instance.

## Directory structure

This is the directory structure we use for the project:

- `commands/` *(Go code)*: Contains all the CLI commands. This package relies heavily on github.com/ditrit/verdeter.
- `configuration/` *(Go code)*: Contains all the configuration holders. Please only use the interfaces, they are all mocked for easy testing
- `controllers/` *(Go code)*: Contains all the http controllers, they handle http requests and consume services.
- `docs/`: Contains the documentation.
- `features/`: Contains all the feature tests (or end to end tests).
- `logger/` *(Go code)*: Contains the logger creation logic. Please don't call it from your own services and code, use the dependency injection system.
- `persistance/` *(Go code)*: 
  - `/gormdatabase/` *(Go code)*: Contains the logic to create a <https://gorm.io> database. Also contains a go package named `gormzap`: it is a compatibility layer between *gorm.io/gorm* and *github.com/uber-go/zap*.
  - `/models/` *(Go code)*: Contains the models. (For a structure to me considered a valid model, it has to embed `models.BaseModel` and satisfy the `models.Tabler` interface. This interface returns the name of the sql table.)
    - `/dto/` *(Go code)*: Contains the Data Transfert Objects. They are used mainly to decode json payloads.
  - `/pagination/` *(Go code)*: Contains the pagination logic.
  - `/repository/` *(Go code)*: Contains the repository interface and implementation. Use uint as ID when using gorm models.
- `resources/` *(Go code)*: Contains the resources shared with the rest of the codebase (ex: API version).
- `router/` *(Go code)*: Contains http router of badaas.
  - `/middlewares/` *(Go code)*: Contains the various http middlewares that we use.
- `scripts/e2e/` : Contains the docker-compose file for end-to-end test.
  - `/api/` : Contains the Dockerfile to build badaas with a dedicated config file.
  - `/db/` : Contains the Dockerfile to build a developpement version of CockroachDB.
- `services/` *(Go code)*: Contains the Dockerfile to build a developpement version of CockroachDB.
  - `/auth/protocols/`: Contains the implementations of authentication clients for differents protocols. 
    - `/basicauth/` *(Go code)*: Handle the authentication using email/password.
    - `/oidc/` *(Go code)*: Handle the authentication via Open-ID Connect.
  - `/sessionservice/` *(Go code)*: Handle sessions and their lifecycle.
  - `/userservice/` *(Go code)*: Handle users.
  - `validators/` *(Go code)*: Contains validators such as an email validator.

At the root of the project, you will find:

- The README.
- The changelog.
- The files for the E2E test http support.
- The LICENCE file.

## Git

### Branch naming policy

`[BRANCH_TYPE]/[BRANCH_NAME]`

- `BRANCH_TYPE` is a prefix to describe the purpose of the branch. 
  
  Accepted prefixes are:
  - `feature`, used for feature development
  - `bugfix`, used for bug fix
  - `improvement`, used for refacto
  - `library`, used for updating library
  - `prerelease`, used for preparing the branch for the release
  - `release`, used for releasing project
  - `hotfix`, used for applying a hotfix on main
  - `poc`, used for proof of concept 
- `BRANCH_NAME` is managed by this regex: `[a-z0-9._-]` (`_` is used as space character).

### Default branch

The default branch is `main`. Direct commit on it is forbidden. The only way to update the application is through pull request.

Release tag are only done on the `main` branch.

### How to release

We use [Semantic Versioning](https://semver.org/spec/v2.0.0.html) as guideline for the version management.

Steps to release:

- Create a new branch labeled `release/vX.Y.Z` from the latest `main`.
- Improve the version number in `changelog.md` and `resources/api.go`.
- Verify the content of the `changelog.md`.
- Commit the modifications with the label `Release version X.Y.Z`.
- Create a pull request on github for this branch into `main`.
- Once the pull request validated and merged, tag the `main` branch with `vX.Y.Z`
- After the tag is pushed, make the release on the tag in GitHub
