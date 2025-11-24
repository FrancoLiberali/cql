==============================
Logger
==============================

When connecting to the database, i.e. when creating the `cql.DB` object, 
it is possible to configure the type of logger to use, the logging level, among others. 
As explained in the :ref:`connection section <cql/connecting_to_a_database:Connection>`, 
this can be done by using the `cql.Open` method:

.. code-block:: go

  db, err = cql.Open(
    dialector,
    &cql.Config{
      Logger: logger.Default,
    },
  )

Any logger that complies with `logger.Interface` can be configured.

Log levels
------------------------------

The log levels provided by cql are the same as those of gorm:

- `logger.Error`: To only view error messages in case they occur during the execution of a sql query.
- `logger.Warn`: The previous level plus warnings for execution of queries and transactions that take 
  longer than a certain time 
  (configurable with SlowQueryThreshold and SlowTransactionThreshold respectively, 200ms by default).
- `logger.Info`: The previous level plus information messages for each query and transaction executed.

Transactions
------------------

For the logs corresponding to transactions 
(slow transactions and transaction execution) 
to be performed, it is necessary to use the cql.Transaction method.

Default logger
-------------------------------

cql provides a default logger that will print Slow SQL and happening errors. 

You can create one with the default configuration using 
(take into account that logger is github.com/FrancoLiberali/cql/logger 
and gormLogger is gorm.io/gorm/logger):

.. code-block:: go

  logger.Default

or use `logger.New` to customize it:

.. code-block:: go

  logger.New(logger.Config{
    LogLevel:                  gormLogger.Warn,
    SlowQueryThreshold:        200 * time.Millisecond,
    SlowTransactionThreshold:  200 * time.Millisecond,
    IgnoreRecordNotFoundError: false,
    ParameterizedQueries:      false,
    Colorful:                  true,
  })

The LogLevel is also configurable via the `ToLogMode` method. 

**Example**

.. code-block:: bash

  example.go:30 [10.392ms] [rows:1] INSERT INTO "products" ("id","created_at","updated_at","deleted_at","string","int","float","bool") VALUES ('4e6d837b-5641-45c9-a028-e5251e1a18b1','2023-07-21 17:19:59.563','2023-07-21 17:19:59.563',NULL,'',1,0.000000,false)

Zap logger
------------------------------

cql provides the possibility to use `zap <https://github.com/uber-go/zap>`_ as logger. 
For this, there is a package called `gormzap`. 
The information displayed by the zap logger will be the same as if we were using the default logger 
but in a structured form, with the following information:

* level: ERROR, WARN or DEBUG
* message:

  * query_error for errors during the execution of a query (ERROR)
  * query_slow for slow queries (WARN)
  * transaction_slow for slow transactions (WARN)
  * query_exec for query execution (DEBUG)
  * transaction_exec for transaction execution (DEBUG)
* error: <error_message> (for errors only)
* elapsed_time: query or transaction execution time
* rows_affected: number of rows affected by the query
* sql: query executed

You can create one with the default configuration using:

.. code-block:: go

  gormzap.NewDefault(zapLogger)

where `zapLogger` is a zap logger, or use `gormzap.New` to customize it:

.. code-block:: go

  gormzap.New(zapLogger, logger.Config{
    LogLevel:                  logger.Warn,
    SlowQueryThreshold:        200 * time.Millisecond,
    SlowTransactionThreshold:  200 * time.Millisecond,
    IgnoreRecordNotFoundError: false,
    ParameterizedQueries:      false,
  })

The LogLevel is also configurable via the `ToLogMode` method. 
Any configuration of the zap logger is done directly during its creation following the 
`zap documentation <https://pkg.go.dev/go.uber.org/zap#hdr-Configuring_Zap>`_. 
Note that the zap logger has its own level setting, so the lower of the two settings 
will be the one finally used.

**Example**

.. code-block:: bash

  DEBUG	example.go:107	query_exec	{"elapsed_time": "3.291981ms", "rows_affected": "1", "sql": "SELECT products.* FROM \"products\" WHERE products.int = 1 AND \"products\".\"deleted_at\" IS NULL"}
