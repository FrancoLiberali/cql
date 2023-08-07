==============================
Quickstart
==============================

To quickly get badaas up and running, you can head to the 
`example <https://github.com/ditrit/badaas-example>`_. 
By following its README.md, you will see how to use badaas and it will be util 
as a template to start your own project.

Step-by-step instructions
-----------------------------------

Once you have started your project with `go init`, you must add the dependency to badaas:

.. code-block:: bash

  go get -u github.com/ditrit/badaas

Then, you can use the following structure to configure and start your application

.. code-block:: go

  func main() {
    badaas.BaDaaS.AddModules(
      // add badaas modules
    ).Provide(
      // provide constructors
    ).Invoke(
      // invoke functions
    ).Start()
  }

Config badaas functionalities
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

You are free to choose which badaas functionalities you wish to use. 
To add them, you must add the corresponding module, for example:

.. code-block:: go

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

Add your own functionalities
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

With the "Provide" and "Invoke" functions you will be able to add your own functionalities to the application. 
For example, to add a route you must first provide the constructor of the controller and 
then invoke the function that adds this route:

.. code-block:: go

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

For details visit :doc:`functionalities`.

Run it
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Once you have defined the functionalities of your project (the `/hello` route in this case), 
you can run the application using the steps described in the example README.md
