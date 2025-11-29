==============================
cqllint
==============================

`cqllint` is a Go linter that checks that cql queries will not generate run-time errors. 

While, in most cases, queries created using cql are checked at compile time, 
there are still some cases that can generate run-time errors (see :doc:`/cql/type_safety`).

cqllint analyses the Go code written to detect these cases and fix them without the need to execute the query. 
It also adds other detections that would not generate runtime errors but are possible misuses of cql.

.. note::

    At the moment, the errors cql.ErrFieldModelNotConcerned, cql.ErrFieldIsRepeated, 
    cql.ErrAppearanceMustBeSelected and cql.ErrAppearanceOutOfRange are detected.

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

Detections
-------------------------------

cqllint has two types of detections: errors and misuses. 
Errors are those that would generate an error at runtime, 
while misuses would not generate an error but are an indication that the code is incorrect.

An example of an error is the detection of :ref:`cql.ErrFieldModelNotConcerned <cql/query_type_safety:Dynamic operators and functions>` in cql.Query.

On the contrary, an example of misuse is the use of :ref:`cql/update:Repeated sets` in cql.Update.

The list of each of the detections performed by cqllint can be found at:

- Query: :ref:`cql/query_type_safety:Type safety limitations and cqllint`.
- Select: :ref:`cql/select:Type safety limitations and cqllint`.
- Insert: :ref:`cql/insert:Type safety limitations and cqllint`.
- Update: :ref:`cql/update:Type safety limitations and cqllint`.
- Delete: :ref:`cql/delete:Type safety`.

Scope and limitations
-------------------------

cqllint analyzes the entire scope of a CQL method call, 
so detection works both outside and inside the function call:

.. code-block:: go
    :class: with-errors
    :caption: Inside function call
    :emphasize-lines: 5
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Concat(
            conditions.Country.Name,
        ).Is().Eq(cql.String("error")),
    ).Find()

.. code-block:: go
    :class: with-errors
    :caption: In variable
    :emphasize-lines: 7
    :linenos:

    countryName := conditions.Country.Name

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Concat(
            countryName,
        ).Is().Eq(cql.String("error")),
    ).Find()

In these cases the detection will be:

.. code-block:: none

    $ cqllint ./...
    example.go:5: models.Country is not joined by the query

It also works within lists:

.. code-block:: go
    :class: with-errors
    :caption: In list
    :emphasize-lines: 3
    :linenos:

    conditions := []condition.Condition[models.City]{
        conditions.City.Name.Is().Eq(
            conditions.Country.Name,
        ),
    }

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions...,
    ).Find()

.. code-block:: none

    $ cqllint ./...
    example.go:3: models.Country is not joined by the query

On the contrary, it cannot go beyond the current scope, so, for example, 
it will not be able to detect parameters that a function receives.

.. code-block:: go
    :class: with-errors
    :caption: In parameter
    :emphasize-lines: 7
    :linenos:

    countryName := conditions.Country.Name

    func doQuery(conditions []condition.Condition[models.City]) error {
        _, err := cql.Query[models.City](
            context.Background(),
            db,
            conditions...,
        ).Find()

        return err
    }

It also cannot analyze the conditions that our code has on the conditions to be used, 
considering that every condition present in the code will be used, for example:

.. code-block:: go
    :class: with-errors
    :caption: In list
    :emphasize-lines: 3
    :linenos:

    conditions := []condition.Condition[models.City]{}

    joinCountry := false

    if joinCountry {
        conditions := append(conditions, conditions.City.Country())
    }

    conditions := append(conditions, conditions.City.Name.Is().Eq(
        conditions.Country.Name,
    ))

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions...,
    ).Find()

In this case, since joinCountry is false, at runtime the join with country will not be performed and will result in an error, 
but cqllint will consider the join to be present.
