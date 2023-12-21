==============================
Quickstart
==============================

To integrate cql into your project, you can head to the 
`quickstart <https://github.com/FrancoLiberali/cql-quickstart>`_.

Refer to its README.md for running it.

Understand it
----------------------------------

Once you have started your project with `go init`, you must add the dependency to cql:

.. code-block:: bash

    go get -u github.com/FrancoLiberali/cql gorm.io/gorm

Create a package for your :ref:`models <cql/concepts:model>`, for example:

.. code-block:: go

  package models

  import (
    "github.com/FrancoLiberali/cql/model"
  )

  type MyModel struct {
    model.UUIDModel

    Name string
  }

Once done, you can :ref:`generate the conditions <cql/concepts:conditions generation>` 
to perform queries on them. 
Create a new package named conditions and add a file with the following content:

.. code-block:: go

  package conditions

  //go:generate cql-gen ../models

Then, you can generate the conditions using `cql-gen` as described in the `README.md <https://github.com/FrancoLiberali/cql-quickstart/blob/main/README.md>`_.

In main.go create a main function that creates a :ref:`gorm.DB <cql/concepts:GormDB>`
that allows connection with the database and call the :ref:`AutoMigrate <cql/concepts:auto migration>` 
method with the models you want to be persisted:

.. code-block:: go

  func main() {
    gormDB, err := NewDBConnection()
    if err != nil {
      panic(err)
    }

    err = gormDB.AutoMigrate(
      models.MyModel{},
    )
    if err != nil {
      panic(err)
    }

    // You are ready to do queries with cql.Query[models.MyModel]
  }

  func NewDBConnection() (*gorm.DB, error) {
    return cql.Open(
      postgres.Open(
        fmt.Sprintf(
          "user=%s password=%s host=%s port=%d sslmode=%s dbname=%s",
          "root", "postgres", "localhost", 26257, "disable", "cql_db",
        ),
      ),
      &gorm.Config{
        Logger: logger.Default.ToLogMode(logger.Info),
      },
    )
  }

Use it
----------------------

Now that you know how to integrate cql into your project, 
you can learn how to use it by following the :doc:`tutorial`.