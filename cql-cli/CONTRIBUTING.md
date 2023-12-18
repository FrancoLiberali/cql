# Contribute to the development of badaas

- [Contribute to the development of badaas](#contribute-to-the-development-of-badaas)
  - [Local compilation](#local-compilation)
  - [Tests](#tests)
  - [Git](#git)
    - [Branch naming policy](#branch-naming-policy)
    - [Default branch](#default-branch)
    - [How to release](#how-to-release)

## Local compilation

You can make modifications to the badaas-cli source code and compile it locally with:

```bash
go build .
```

You can then run the badaas-cli executable directly or add a link in your $GOPATH to run it from a project:

```bash
ln -sf badaas-cli $GOPATH/bin/badaas-cli
```

## Tests

We use the standard test suite in combination with [github.com/stretchr/testify](https://github.com/stretchr/testify) to do our unit testing.

To run them, please run:

```sh
make test
```

## Git

### Branch naming policy

`[BRANCH_TYPE]/[BRANCH_NAME]`

- `BRANCH_TYPE` is a prefix to describe the purpose of the branch.
  Accepted prefixes are:
  - `feature`, used for feature development
  - `bugfix`, used for bug fix
  - `improvement`, used for refactor
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
- Improve the version number in `changelog.md`.
- Verify the content of the `changelog.md`.
- Commit the modifications with the label `Release version X.Y.Z`.
- Create a pull request on github for this branch into `main`.
- Once the pull request validated and merged, tag the `main` branch with `vX.Y.Z`.
- After the tag is pushed, make the release on the tag in GitHub.
