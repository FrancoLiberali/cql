# BADAAS: Backend And Distribution As A Service

Badaas enables the effortless construction of ***distributed, resilient, highly available and secure applications by design***, while ensuring very simple deployment and management (NoOps). 

Badaas provides several key features:

- **Authentification**: Badaas can authentify users using its internal authentification scheme or externally by using protocols such as OIDC, SAML, Oauth2...
- **Habilitation**: On a resource access, Badaas will check if the user is authorized using a RBAC model.
- **Distribution**: Badaas is built to run in clusters by default. Communications between nodes are TLS encrypted using [shoset](https://github.com/ditrit/shoset).
- **Persistence**: Applicative objects are persisted as well as user files. Those resources are shared accross the clusters to increase resiliency.
- **Querying Resources**: Resources are accessible via a REST API.
- **Posix complient**: Badaas strives towards being a good unix citizen and respecting commonly accepted norms. (see [Configuration](#configuration))
- **Advanced logs management**: Badaas provides an interface to interact with the logs produced by the clusters. Logs are formated in json by default.

To quickly get badaas up and running, please head to the [miniblog tutorial](<!-- TODO: link the miniblog tutorial here -->)

- [Quickstart](#quickstart)
- [Docker install](#docker-install)
- [Install from sources](#install-from-sources)
  - [Prerequisites](#prerequisites)
  - [Configuration](#configuration)
- [Contributing](#contributing)
- [Licence](#licence)

## Quickstart

You can either use the [Docker Install](#docker-install) or build it from source .

## Docker install

You can build the image using `docker build -t badaas .` since we don't have an official docker image yet.

## Install from sources

### Prerequisites

Get the sources of the project, either by visiting the [releases](https://github.com/ditrit/badaas/releases) page and downloading an archive or clone the main branch (please be aware that is it not a stable version).

To build the project:

- [Install go](https://go.dev/dl/#go1.18.4) v1.18
- Install project dependencies

```bash
go get
```

- Run build command

```bash
go build .
```

Well done, you have a binary `badaas` at the root of the project.

Then you can launch Badaas directly with:

```bash
export BADAAS_DATABASE_PORT=<complete>
export BADAAS_DATABASE_HOST=<complete>
export BADAAS_DATABASE_DBNAME=<complete>
export BADAAS_DATABASE_SSLMODE=<complete>
export BADAAS_DATABASE_USERNAME=<complete>
export BADAAS_DATABASE_PASSWORD=<complete>
./badaas 
```

### Configuration

Badaas use [verdeter](https://github.com/ditrit/verdeter) to manage it's configuration. So Badaas is POSIX complient by default.

Badaas can be configured using environment variables, configuration files or CLI flags.
CLI flags take priority on the environment variables and the environment variables take priority on the content of the configuration file.

As an exemple we will define the `database.port` configuration key using the 3 methods:

- Using a CLI flag: `--database.port=1222`
- Using an environment variable: `export BADAAS_DATABASE_PORT=1222` (*dots are replaced by underscores*)
- Using a config file (in YAML here):

    ```yml
    # /etc/badaas/badaas.yml
    database:
        port: 1222
    ```

The config file can be placed at `/etc/badaas/badaas.yml` or `$HOME/.config/badaas/badaas.yml` or in the same folder as the badaas binary `./badaas.yml`.

If needed, the location can be overridden using the config key `config_path`.

***For a full overview of the configuration keys: please head to the [configuration documentation](./configuration.md).***

## Contributing

See [this section](./CONTRIBUTING.md).

## Licence

Badaas is Licenced under the [Mozilla Public License Version 2.0](./LICENSE).
