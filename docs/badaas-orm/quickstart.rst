==============================
Quickstart
==============================

To integrate badaas-orm into your project, you can head to the 
`quickstart <https://github.com/ditrit/badaas-orm-quickstart>`_.

Refer to its README.md for running it.

Understand it
----------------------------------

Once you have started your project with `go init`, you must add the dependency to BaDaaS:

.. code-block:: bash

    go get -u github.com/ditrit/badaas

Create a package for your :ref:`models <badaas-orm/concepts:model>`, for example:

.. code-block:: go

  package models

  import (
    "github.com/ditrit/badaas/orm/model"
  )

  type MyModel struct {
    model.UUIDModel

    Name string
  }

Once done, you can :ref:`generate the conditions <badaas-orm/concepts:conditions generation>` 
to perform queries on them. 
Create a new package named conditions and add a file with the following content:

.. code-block:: go

  package conditions

  //go:generate badaas-cli gen conditions ../models

Then, you can generate the conditions using `badaas-cli` as described in the README.md.

In main.go create a main function that creates a :ref:`gorm.DB <badaas-orm/concepts:GormDB>`
that allows connection with the database and call the :ref:`AutoMigrate <badaas-orm/concepts:auto migration>` 
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

    // You are ready to do queries with orm.NewQuery[models.MyModel]
  }

  func NewDBConnection() (*gorm.DB, error) {
    return orm.Open(
      postgres.Open(orm.CreatePostgreSQLDSN("localhost", "root", "postgres", "disable", "badaas_db", 26257)),
      &gorm.Config{
        Logger: logger.Default.ToLogMode(logger.Info),
      },
    )
  }

Use it
----------------------

Now that you know how to integrate badaas-orm into your project, 
you can learn how to use it by following the :doc:`tutorial`.