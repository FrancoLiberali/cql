# badaas-cli <!-- omit in toc -->

`badaas-cli` is the command line tool that makes it possible to configure and run a badaas application.

- [Install with go install](#install-with-go-install)
- [Build from sources](#build-from-sources)
- [Commands](#commands)
  - [badaas-cli gen docker](#badaas-cli-gen-docker)
  - [badaas-cli gen conditions](#badaas-cli-gen-conditions)
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

### badaas-cli gen docker

gen docker is the command you can use to generate the files and configurations necessary for your project to use badaas in a simple way.

`gen docker` will generate the docker and configuration files needed to run the application in the `badaas/docker` and `badaas/config` folders respectively.

All these files can be modified in case you need different values than those provided by default. For more information about the configuration head to [configuration docs](github.com/ditrit/badaas/configuration.md).

A Makefile will be generated for the execution of a badaas server, with the command:

```bash
make badaas_run
```

### badaas-cli gen conditions

gen conditions is the command you can use to generate conditions to query your objects using badaas-orm. For each BaDaaS Model found in the input packages a file containing all possible Conditions on that object will be generated, allowing you to use badaas-orm.

Its use is recommended through `go generate`. To see an example of how to do it click [here](https://github.com/ditrit/badaa-orm-example/blob/main/standalone/conditions/orm.go).

## Contributing

See [this section](./CONTRIBUTING.md).

## License

badaas-cli is Licensed under the [Mozilla Public License Version 2.0](./LICENSE).
