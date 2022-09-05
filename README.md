# badaas

Backend and Distribution as a Service

## Build

To build the project: 

- [Install go](https://go.dev/dl/#go1.18.4) v1.18
- Install project dependencies
```
go get
```
- Run build command
```
go build
```

Once all is done, you have a binary `badaas` at the root of the project.

## Development

### Directory structure

This is the default directory structure we use for the project:

```
badaas
├ commands             ⇨ Contains all the CLI commands.
├ configuration        ⇨ Contains configuration holders.
├ controllers          ⇨ Contains all the web controllers.
├ features             ⇨ Contains all the e2e tests.
├ resources            ⇨ Contains all the applications resources, like constants or other.
├ router               ⇨ Contains the route definitions for application.
├ scripts              ⇨ Contains shell scripts, the Docker files and docker-compose files for e2e test.
```

### E2E testing

We use [godog](https://github.com/cucumber/godog) to run all e2e tests.

To execute E2E tests :

```bash
# Build containers and launch db and api .
docker compose -f "scripts/e2e/docker-compose.yml" up --build

# In another process, run the test
go test
```

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

### Git: Default branch

The default branch is main. Direct commit on it is forbidden. The only way to update the application is through pull request.

Release tag are only done on the `main` branch.

### Git: Branch naming policy

`[BRANCH_TYPE]/[BRANCH_NAME]`

* `BRANCH_TYPE` is a prefix to describe the purpose of the branch. Accepted prefixes are:
    * `feature`, used for feature development
    * `bugfix`, used for bug fix
    * `improvement`, used for refacto
    * `library`, used for updating library
    * `prerelease`, used for preparing the branch for the release
    * `release`, used for releasing project
    * `hotfix`, used for applying a hotfix on main
    * `poc`, used for proof of concept 
* `BRANCH_NAME` is managed by this regex: `[a-z0-9._-]` (`_` is used as space character).
