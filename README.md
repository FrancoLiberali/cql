# BADAAS: Backend And Distribution As A Service

Badaas enables the effortless construction of ***distributed, resilient, highly available and secure applications by design***, while ensuring very simple deployment and management (NoOps).

> **Warning**
> BaDaaS is still under development. Each of its components can have a different state of evolution that you can consult in [Features and components](#features-and-components)

- [Features and components](#features-and-components)
- [Quickstart](#quickstart)
  - [Example](#example)
  - [Step-by-step instructions](#step-by-step-instructions)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Features and components

Badaas provides several key features, each provided by a component that can be used independently and has a different state of evolution:

- **Authentication**(unstable): Badaas can authenticate users using its internal authentication scheme or externally by using protocols such as OIDC, SAML, Oauth2...
- **Authorization**(wip_unstable): On resource access, Badaas will check if the user is authorized using a RBAC model.
- **Distribution**(todo): Badaas is built to run in clusters by default. Communications between nodes are TLS encrypted using [shoset](https://github.com/ditrit/shoset).
- **Persistence**(wip_unstable): Applicative objects are persisted as well as user files. Those resources are shared across the clusters to increase resiliency. To achieve this, BaDaaS uses the [BaDaaS ORM](https://github.com/ditrit/badaas/orm) component.
- **Querying Resources**(unstable): Resources are accessible via a REST API.
- **Posix compliant**(stable): Badaas strives towards being a good unix citizen and respecting commonly accepted norms. (see [Configuration](#configuration))
- **Advanced logs management**(todo): Badaas provides an interface to interact with the logs produced by the clusters. Logs are formatted in json by default.

## Quickstart

### Example

To quickly get badaas up and running, you can head to the [example](https://github.com/ditrit/badaas-example). This example will help you to see how to use badaas and as a template to start your own project

### Step-by-step instructions

Once you have started your project with `go init`, you must add the dependency to badaas:

```bash
go get -u github.com/ditrit/badaas
```

Then, you can use the following structure to configure and start your application

```go
func main() {
  badaas.BaDaaS.AddModules(
    // add badaas modules
  ).Provide(
    // provide constructors
  ).Invoke(
    // invoke functions
  ).Start()
}
```

#### Config badaas functionalities

You are free to choose which badaas functionalities you wish to use. To add them, you must add the corresponding module, for example:

```go
func main() {
  badaas.BaDaaS.AddModules(
    badaas.InfoModule,
    badaas.AuthModule,
  ).Provide(
    NewAPIVersion,
  ).Start()
}

func NewAPIVersion() *semver.Version {
  return semver.MustParse("0.0.0-unreleased")
}
```

#### Add your own functionalities

With the "Provide" and "Invoke" functions you will be able to add your own functionalities to the application. For example, to add a route you must first provide the constructor of the controller and then invoke the function that adds this route:

```go
func main() {
  badaas.BaDaaS.Provide(
    NewHelloController,
  ).Invoke(
    AddExampleRoutes,
  ).Start()
}

type HelloController interface {
  SayHello(http.ResponseWriter, *http.Request) (any, httperrors.HTTPError)
}

type helloControllerImpl struct{}

func NewHelloController() HelloController {
  return &helloControllerImpl{}
}

func (*helloControllerImpl) SayHello(response http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
  return "hello world", nil
}

func AddExampleRoutes(
  router *mux.Router,
  jsonController middlewares.JSONController,
  helloController HelloController,
) {
  router.HandleFunc(
    "/hello",
    jsonController.Wrap(helloController.SayHello),
  ).Methods("GET")
}
```

#### Run it

Once you have defined the functionalities of your project (the `/hello` route in this case), you can run the application using the steps described in the example README.md

### Provided functionalities

#### InfoModule

`InfoModule` adds the path `/info`, where the api version will be answered. To set the version you want to be responded you must provide a function that returns it:

```go
func main() {
  badaas.BaDaaS.AddModules(
    badaas.InfoModule,
  ).Provide(
    NewAPIVersion,
  ).Start()
}

func NewAPIVersion() *semver.Version {
  return semver.MustParse("0.0.0-unreleased")
}
```

#### AuthModule

`AuthModule` adds `/login` and `/logout`, which allow us to add authentication to our application in a simple way:

```go
func main() {
  badaas.BaDaaS.AddModules(
    badaas.AuthModule,
  ).Start()
}
```

### Configuration

Badaas use [verdeter](https://github.com/ditrit/verdeter) to manage it's configuration, so Badaas is POSIX compliant by default.

Badgen automatically generates a default configuration in `badaas/config/badaas.yml`, but you are free to modify it if you need to.

This can be done using environment variables, configuration files or CLI flags.
CLI flags take priority on the environment variables and the environment variables take priority on the content of the configuration file.

As an example we will define the `database.port` configuration key using the 3 methods:

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

## License

Badaas is Licensed under the [Mozilla Public License Version 2.0](./LICENSE).
