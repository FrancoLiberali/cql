# cql-gen <!-- omit in toc -->

`cql-gen` is the command line tool that makes it possible to use cql in your project.

- [Install with go install](#install-with-go-install)
- [Build from sources](#build-from-sources)
- [Execute](#execute)
- [Contributing](#contributing)
- [License](#license)

## Install with go install

For simply installing it, use:

```bash
go install github.com/FrancoLiberali/cql-gen@latest
```

Or you can build it from sources.

## Build from sources

Get the sources of the project, either by visiting the [releases](https://github.com/FrancoLiberali/cql/releases) page and downloading an archive or clone the main branch (please be aware that is it not a stable version).

To build the project:

- [Install go](https://go.dev/dl/#go1.18.4) v1.18
- `cd cql-gen`
- Install project dependencies

    ```bash
    go get
    ```

- Run build command

    ```bash
    go build .
    ```

Well done, you have a binary `cql-gen` at the root of the project.

## Execute

cql-gen is used to generate conditions to query your objects using cql. For each cql Model found in the input packages a file containing all possible Conditions on that object will be generated, allowing you to use cql.

```bash
Usage:
  cql-gen [flags]

Flags:
  -d, --dest_package string   Destination package (not used if ran with go generate)
  -h, --help                  help for cql-gen
  -v, --verbose               Verbose logging
      --version               version for cql-gen
```

Its use is recommended through `go generate`. To see an example of how to do it visit the [quickstart](https://github.com/FrancoLiberali/cql-quickstart/blob/main/conditions/cql.go).

## Contributing

See [this section](./CONTRIBUTING.md).

## License

cql-gen is Licensed under the [Mozilla Public License Version 2.0](../LICENSE).
