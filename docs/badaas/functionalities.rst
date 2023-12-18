==============================
Functionalities
==============================

InfoModule
-------------------------------

`InfoModule` adds the path `/info`, where the api version will be answered. 
To set the version you want to be responded you must provide a function that returns it:

.. code-block:: go

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

AuthModule
-------------------------------

`AuthModule` adds `/login` and `/logout`, which allow us to add authentication to our 
application in a simple way:

.. code-block:: go

  func main() {
    badaas.BaDaaS.AddModules(
      badaas.AuthModule,
    ).Start()
  }
