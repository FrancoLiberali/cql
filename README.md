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
├ controllers          ⇨ Contains all the web controllers.
├ features             ⇨ Contains all the e2e tests.
├ resources            ⇨ Contains all the applications resources, like constants or other.
├ router               ⇨ Contains the route definitions for application.
```

### E2E testing

We use [godog](https://github.com/cucumber/godog) to run all e2e tests.

To execute E2E tests :

```
# Build the project
go build

# Run the application
./badaas

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
