==============================
cqllint
==============================

`cqllint` is a Go linter that checks that cql queries will not generate run-time errors. 

While, in most cases, queries created using cql are checked at compile time, 
there are still some cases that can generate run-time errors (see :ref:`cql/type_safety:Runtime errors`).

cqllint analyses the Go code written to detect these cases and fix them without the need to execute the query.

.. note::

    At the moment, only the error cql.ErrFieldModelNotConcerned is detected.

We recommend integrating cqllint into your CI so that the use of cql ensures 100% that your queries will be executed correctly.

Installation
----------------------------

For simply installing it, use:

.. code-block:: bash

    go install github.com/FrancoLiberali/cql/cqllint@latest

.. note::

    At the moment, only the error cql.ErrFieldModelNotConcerned is detected.

.. warning::

    The version of cqllint used must be the same as the version of cql. 
    You can install a specific version using `go install github.com/FrancoLiberali/cql/cqllint@vX.Y.Z`, 
    where X.Y.Z is the version number.

Execution
----------------------------

cqllint can be used independently by running:

.. code-block:: bash

    cqllint ./...

o using `go vet`:

.. code-block:: bash

    go vet -vettool=$(which cqllint) ./...

Example
----------------------------

The simplest example of an error case is trying to make a comparison 
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