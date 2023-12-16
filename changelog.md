# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)

## [Unreleased]

### Added

- Setup project (ci and sonar).
- Setup e2e test solution (cucumber + docker).
- Setup Docker based build system.
- Add default api endpoint `info`.
- Setup command based pattern using verdeter.
- Add an http error handling mechanism.
- Add a json controller.
- Add a dto package.
- The tasks in the CI are ran in parallel.
- Add a Generic CRUD Repository.
- Add a configuration structure containing all the configuration holder.
- Refactor codebase to use the DI framework uber-go/fx. Now all services and controllers relies on interfaces.
- Add an generic ID to the repository interface.
- Add a retry mechanism for the database connection.
- Add `init` flag to migrate database and create admin user.
- Add a CONTRIBUTING.md and a documentation file for configuration (configuration.md)
- Add a session services.
- Add a basic authentication controller.
- Now config keys are only declared once with constants in the `configuration/` package.
- Add a dto that is returned on a successful login.
- Update verdeter to version v0.4.0.
- Transform BadAas into a library.

[unreleased]: https://github.com/ditrit/badaas/blob/main/changelog.md#unreleased
