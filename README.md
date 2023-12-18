# BADAAS: Backend And Distribution As A Service

[![Build Status](https://github.com/ditrit/badaas/actions/workflows/CI.yml/badge.svg?branch=main)](https://github.com/ditrit/badaas/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ditrit/badaas)](https://goreportcard.com/report/github.com/ditrit/badaas)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ditrit_badaas&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ditrit_badaas)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ditrit_badaas&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ditrit_badaas)
[![OpenSSF Best Practices](https://bestpractices.coreinfrastructure.org/projects/7624/badge)](https://bestpractices.coreinfrastructure.org/projects/7624)

[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/ditrit/badaas)

[![Discord DitRit](https://dcbadge.vercel.app/api/server/zkKfj9gj2C?style=flat&theme=default-inverted)](https://discord.gg/zkKfj9gj2C)

BaDaaS enables the effortless construction of ***distributed, resilient, highly available and secure applications by design***, while ensuring very simple deployment and management (NoOps).

> **Warning**
> BaDaaS is still under development and each of its components can have a different state of evolution

## Features and components

Badaas provides several key features, each provided by a component that can be used independently and has a different state of evolution:

- **Authentication**(unstable): Badaas can authenticate users using its internal authentication scheme or externally by using protocols such as OIDC, SAML, Oauth2...
- **Authorization**(wip_unstable): On resource access, Badaas will check if the user is authorized using a RBAC model.
- **Distribution**(todo): Badaas is built to run in clusters by default. Communications between nodes are TLS encrypted using [shoset](https://github.com/ditrit/shoset).
- **Persistence**(wip_unstable): Applicative objects are persisted as well as user files. Those resources are shared across the clusters to increase resiliency. To achieve this, BaDaaS uses the [badaas-orm](https://github.com/ditrit/badaas/orm) component.
- **Querying Resources**(unstable): Resources are accessible via a REST API.
- **Posix compliant**(stable): Badaas strives towards being a good unix citizen and respecting commonly accepted norms.
- **Advanced logs management**(todo): Badaas provides an interface to interact with the logs produced by the clusters. Logs are formatted in json by default.

## Documentation

<!-- TODO add link to docs -->

## Contributing

See [this section](./docs/contributing/contributing.md).

## Code of Conduct

This project has adopted the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md)

## License

Badaas is Licensed under the [Mozilla Public License Version 2.0](./LICENSE).
