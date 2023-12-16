# badaas-cli <!-- omit in toc -->

`badaas-cli` is the command line tool that makes it possible to configure and run a badaas application.

- [Install with go install](#install-with-go-install)
- [Build from sources](#build-from-sources)
- [Commands](#commands)
  - [badaas-cli gen](#badaas-cli-gen)
- [Contributing](#contributing)
- [License](#license)

## Install with go install

For simply installing it, use:

```bash
go install github.com/ditrit/badaas-cli
```

Or you can build it from sources.

## Build from sources

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

Well done, you have a binary `badaas-cli` at the root of the project.

## Commands

You can see the available commands by running:

```bash
badaas-cli help
```

For more information about the functionality provided and how to use each command use:

```bash
badaas-cli help [command]
```

### badaas-cli gen

gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.

`gen` will generate the docker and configuration files needed to run the application in the `badaas/docker` and `badaas/config` folders respectively.

All these files can be modified in case you need different values than those provided by default. For more information about the configuration head to [configuration docs](../../configuration.md).

A Makefile will be generated for the execution of a badaas server, with the command:

```bash
make badaas_run
```

## Contributing

See [this section](./CONTRIBUTING.md).

## License

badaas-cli is Licensed under the [Mozilla Public License Version 2.0](./LICENSE).
