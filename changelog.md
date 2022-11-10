# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)

## [Unreleased]

### Added

- Setup project (ci and sonar)
- Setup e2e test solution (cucumber + docker).
- Setup Docker based build system.
- Add default api endpoint `info`
- Setup command based pattern using verdeter
- Add an http error handling mecanism
- Add a json controller
- Add a dto package
- The tasks in the CI are ran in parallel.
- Add a Generic CRUD Repository
- Add a configuration structure containing all the configuration holder.
- Refactor codebase to use the DI framework uber-go/fx. Now all services and controllers relies on interfaces.
- Add an generic ID to the repository interface

[unreleased]: https://github.com/ditrit/badaas/blob/main/changelog.md#unreleased