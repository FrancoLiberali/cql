==============================
Configuration
==============================

Methods
-------------------------------

Badaas use `verdeter <https://github.com/ditrit/verdeter>`_ to manage it's configuration, 
so Badaas is POSIX compliant by default.

Badaas-cli automatically generates a default configuration in `badaas/config/badaas.yml`, 
but you are free to modify it if you need to.

This can be done using environment variables, configuration files or CLI flags.
CLI flags take priority on the environment variables and the environment variables take 
priority on the content of the configuration file.

As an example we will define the `database.port` configuration key using the 3 methods:

- Using a CLI flag: :code:`--database.port=1222`
- Using an environment variable: :code:`export BADAAS_DATABASE_PORT=1222` (*dots are replaced by underscores*)
- Using a config file (in YAML here)

.. code-block:: yaml

  # /etc/badaas/badaas.yml
  database:
    port: 1222

The config file can be placed at `/etc/badaas/badaas.yml` or `$HOME/.config/badaas/badaas.yml` 
or in the same folder as the badaas binary `./badaas.yml`.

If needed, the location can be overridden using the config key `config_path`.

Config file
----------------------------

In this section, we will focus our attention on config files but 
we won't forget that we can use environment variables and CLI flags to change Badaas' config.

The config file can be formatted in any syntax that 
`viper <https://github.com/spf13/viper>`_ supports but we will only use YAML syntax in our docs.

To see a complete example of a working config file: head to 
`badaas.example.yml <https://github.com/ditrit/badaas/blob/main/badaas.example.yml>`_.

Database
^^^^^^^^^^^^^^^^^^^^^^^^
.. code-block:: yaml

  # The settings for the database.
  database:
    # The host of the database server. 
    # (mandatory)
    host: e2e-db-1

    # The port of the database server. 
    # (mandatory)
    port: 26257

    # The name of the database to use. 
    name: badaas_db

    # The sslmode of the connection to the database server. 
    # (mandatory)
    sslmode: disable

    # The username of the account on the database server. 
    # (mandatory)
    username: root

    # The password of the account on the database server.
    # (mandatory)
    password: postgres

    # The settings for the initialization of the database server. 
    init:
      # Number of time badaas will try to establish a connection to the database server.
      # default (10)
      retry: 10

      # Waiting time between connection, in seconds.
      # default (5)
      retryTime: 5

Please note that the init section `init:` is not mandatory. 
Badaas is suited with a simple but effective retry mechanism that will retry 
`database.init.retry` time to establish a connection with the database. 
Badaas will wait `database.init.retryTime` seconds between each retry.

Logger
^^^^^^^^^^^^^^^^^^^^^^^^

Badaas use a structured logger that can output json logs in 
production and user adapted logs for debug using the `logger.mode` key.

Badaas offers the possibility to change the log message of the 
Middleware Logger but provides a sane default. It is formatted using the Jinja syntax. 
The values available are `method`, `url` and `protocol`.

.. code-block:: yaml

  # The settings for the logger.
  logger:
    # Either `dev` or `prod`
    # default (`prod`)
    mode: prod

    # Disable error stacktrace from logs
    # default (true)
    disableStacktrace: true

    # Threshold for the slow query warning in milliseconds
    # default (200)
    # use 0 to disable slow query warnings
    slowQueryThreshold: 200

    # Threshold for the slow transaction warning in milliseconds
    # default (200)
    # use 0 to disable slow transaction warnings
    slowTransactionThreshold: 200

    # If true, ignore gorm.ErrRecordNotFound error for logger
    # default (false)
    ignoreRecordNotFoundError: false

    # If true, don't include params in the query execution logs
    # default (false)
    parameterizedQueries: false

    request:
      # Change the log emitted when badaas receives a request on a valid endpoint.
      template: "Receive {{method}} request on {{url}}"

HTTP Server
^^^^^^^^^^^^^^^^^^^^^^^^

You can change the host Badaas will bind to, the port and the timeout in seconds.

Additionally you can change the number of elements returned by default for a paginated response

.. code-block:: yaml

  # The settings for the http server.
  server:
    # The address to bind badaas to.
    # default ("0.0.0.0")
    host: "" 

    # The port badaas should use.
    # default (8000)
    port: 8000

    # The maximum timeout for the http server in seconds.
    # default (15)
    timeout: 15 

    # The settings for the pagination.
    pagination:
      page:
        # The maximum number of record per page 
        # default (100)
        max: 100


Default values
^^^^^^^^^^^^^^^^^^^^^^^^

The section allow to change some settings for the first run.

.. code-block:: yaml

  # The settings for the first run.
  default:
    # The admin settings for the first run
    admin:
      # The admin password for the first run. Won't change is the admin user already exists.
      password: admin

Session management
^^^^^^^^^^^^^^^^^^^^^^^^

You can change the way the session service handle user sessions.
Session are extended if the user made a request to badaas in the "roll duration". 
The session duration and the refresh interval of the cache can be changed. 
They contains some good defaults.

Please see the diagram below to see what is the roll duration relative to the session duration.
::

        |   session duration                        |
        |<----------------------------------------->|
    ----|-------------------------|-----------------|----> time
        |                         |                 |
                                  |<--------------->|
                                     roll duration

.. code-block:: yaml

  # The settings for session service
  # This section contains some good defaults, don't change those value unless you need to.
  session:
    # The duration of a user session, in seconds
    # Default (14400) equal to 4 hours
    duration: 14400
    # The refresh interval in seconds. Badaas refresh it's internal session cache periodically.
    # Default (30)
    pullInterval: 30
    # The duration in which the user can renew it's session by making a request.
    # Default (3600) equal to 1 hour
    rollDuration: 3600