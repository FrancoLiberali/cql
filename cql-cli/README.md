# cql-cli <!-- omit in toc -->

`cql-cli` is the command line tool that makes it possible to use cql in your project.

- [Install with go install](#install-with-go-install)
- [Build from sources](#build-from-sources)
- [Commands](#commands)
  - [cql-cli gen conditions](#cql-cli-gen-conditions)
- [Contributing](#contributing)
- [License](#license)

## Install with go install

For simply installing it, use:

```bash
go install github.com/FrancoLiberali/cql/cql-cli
```

Or you can build it from sources.

## Build from sources

Get the sources of the project, either by visiting the [releases](https://github.com/FrancoLiberali/cql/releases) page and downloading an archive or clone the main branch (please be aware that is it not a stable version).

To build the project:

- [Install go](https://go.dev/dl/#go1.18.4) v1.18
- `cd cql-cli`
- Install project dependencies

    ```bash
    go get
    ```

- Run build command

    ```bash
    go build .
    ```

Well done, you have a binary `cql-cli` at the root of the project.

## Commands

You can see the available commands by running:

```bash
cql-cli help
```

For more information about the functionality provided and how to use each command use:

```bash
cql-cli help [command]
```

### cql-cli gen conditions

gen conditions is the command you can use to generate conditions to query your objects using cql. For each cql.Model found in the input packages a file containing all possible Conditions on that object will be generated, allowing you to use cql.

Its use is recommended through `go generate`. To see an example of how to do it click [here](https://github.com/ditrit/badaa-orm-example/blob/main/standalone/conditions/orm.go).

## Contributing

See [this section](./CONTRIBUTING.md).

## License

cql-cli is Licensed under the [Mozilla Public License Version 2.0](../LICENSE).
