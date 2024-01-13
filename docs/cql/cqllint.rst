==============================
cqllint
==============================

`cqllint` is a Go linter that checks that cql queries will not generate run-time errors. 

While, in most cases, queries created using cql are checked at compile time, 
there are still some cases that can generate run-time errors (see :ref:`cql/type_safety:Runtime errors`).

cqllint analyses the Go code written to detect these cases and fix them without the need to execute the query. 
It also adds other detections that would not generate runtime errors but are possible misuses of cql.

.. note::

    At the moment, only the errors cql.ErrFieldModelNotConcerned and cql.ErrFieldIsRepeated are detected.

We recommend integrating cqllint into your CI so that the use of cql ensures 100% that your queries will be executed correctly.

Installation
----------------------------

For simply installing it, use:

.. code-block:: bash

    go install github.com/FrancoLiberali/cql/cqllint@latest

.. warning::

    The version of cqllint used must be the same as the version of cql. 
    You can install a specific version using `go install github.com/FrancoLiberali/cql/cqllint@vX.Y.Z`, 
    where X.Y.Z is the version number.

Execution
----------------------------

cqllint can be used independently by running:

.. code-block:: bash

    cqllint ./...

or using `go vet`:

.. code-block:: bash

    go vet -vettool=$(which cqllint) ./...

Errors
-------------------------------

ErrFieldModelNotConcerned
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The simplest example this error case is trying to make a comparison 
with an attribute of a model that is not joined by the query:

.. code-block:: go
    :caption: example.go
    :class: with-errors
    :emphasize-lines: 3
    :linenos:

    _, err := cql.Query[models.Brand](
        db,
        conditions.Brand.Name.IsDynamic().Eq(conditions.City.Name.Value()),
    ).Find()

If we execute this query we will obtain an error of type `cql.ErrFieldModelNotConcerned` with the following message:

.. code-block:: none

    field's model is not concerned by the query (not joined); not concerned model: models.City; operator: Eq; model: models.Brand, field: Name

Now, if we run cqllint we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:3: models.City is not joined by the query

In this way, we will be able to correct this error without having to execute the query.

ErrFieldIsRepeated
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The simplest example this error case is trying to set the value of an attribute twice:

.. code-block:: go
    :caption: example.go
    :class: with-errors
    :emphasize-lines: 5,6
    :linenos:

    _, err := cql.Update[models.Brand](
        db,
        conditions.Brand.Name.Is().Eq("nike"),
    ).Set(
        conditions.Brand.Name.Set().Eq("adidas"),
        conditions.Brand.Name.Set().Eq("puma"),
    )

If we execute this query we will obtain an error of type `cql.ErrFieldIsRepeated` with the following message:

.. code-block:: none

    field is repeated; field: models.Brand.Name; method: Set

Now, if we run cqllint we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:5: conditions.Brand.Name is repeated
    example.go:6: conditions.Brand.Name is repeated

In this way, we will be able to correct this error without having to execute the query.

Misuses
-------------------------

Although some cases would not generate runtime errors, cqllint will detect them as they are possible misuses of cql.

Set the same value
^^^^^^^^^^^^^^^^^^^^^^^^^

This case occurs when making a Set of exactly the same value:

.. code-block:: go
    :caption: example.go
    :class: with-errors
    :emphasize-lines: 5
    :linenos:

    _, err := cql.Update[models.Brand](
        db,
        conditions.Brand.Name.Is().Eq("nike"),
    ).Set(
        conditions.Brand.Name.Set().Dynamic(conditions.Brand.Name.Value()),
    )

If we run cqllint we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:5: conditions.Brand.Name is set to itself