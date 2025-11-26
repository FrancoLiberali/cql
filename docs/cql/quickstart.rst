==============================
Quickstart
==============================

To integrate cql into your project, you can head to the 
`quickstart <https://github.com/FrancoLiberali/cql-quickstart>`_.

Run it
----------------------------------

Refer to its `README.md <https://github.com/FrancoLiberali/cql-quickstart/blob/main/README.md>`_ for running it.

Understand it
----------------------------------

Once you have started your project with `go init`, you must add the dependency to cql:

.. code-block:: bash

    go get github.com/FrancoLiberali/cql

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
to perform queries on them. In this case, the file `conditions/cql.go` has the following content:

.. code-block:: go

  package conditions

  //go:generate cql-gen ../models

Then, you can generate the conditions running:

.. code-block:: bash

  go generate ./...

In `main.go` there is a main function that creates a :ref:`cql.DB <cql/concepts:cqlDB>`
that allows connection with the database and calls the :ref:`AutoMigrate <cql/concepts:auto migration>` 
method with the models you want to be persisted.

After this, you are ready to query your objects using cql.Query.

Now that you know how to integrate cql into your project, 
you can learn how to use it by following the :doc:`tutorial`.