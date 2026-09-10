# Maintaining

This document is intended for CQL maintainers only.

## How to release

Release tag are only done on the `main` branch. We use [Semantic Versioning](https://semver.org/spec/v2.0.0.html) as guideline for the version management.

Steps to release:

- Create a new branch labeled `release/vX.Y.Z` from the latest `main`.
- Improve the version number in `cql-gen/version/version.go` and `cqllint/version/version.go`.
- If this release bumps the gorm fork version, update it in **both** `go.mod` and `docs/cql/quickstart.rst` in the same commit (see [The gorm fork](#the-gorm-fork)).
- Commit the modifications with the label `Release version X.Y.Z`.
- Create a pull request on github for this branch into `main`.
- Once the pull request validated and merged, tag the `main` branch using `./create_tag.sh X.Y.Z`.

## The gorm fork

CQL depends on a fork of gorm ([`github.com/FrancoLiberali/gorm`](https://github.com/FrancoLiberali/gorm)) that
adds features gorm's maintainers did not accept upstream because they are CQL-specific. The dependency is wired
through a `replace` directive in CQL's `go.mod`:

```
replace gorm.io/gorm => github.com/FrancoLiberali/gorm vX.Y.Z
```

Go applies `replace` directives **only in the main module** — they are ignored when CQL is pulled in as a
dependency. That is why the quickstart instructs consumers to add the same directive to their own `go.mod`
(`docs/cql/quickstart.rst`). The version there is not tracked automatically; **it must be kept in sync with
`go.mod` by hand**.

When bumping the fork version:

- Tag the new fork version in the `github.com/FrancoLiberali/gorm` repo first.
- Update the version in CQL's `go.mod` `replace` line.
- Update the same version in the `go mod edit -replace ...` command in `docs/cql/quickstart.rst`.

Do all three together — a mismatch between `go.mod` and the quickstart command will make consumers' projects
fail to compile against the fork API.
