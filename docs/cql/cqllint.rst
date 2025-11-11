==============================
cqllint
==============================

`cqllint` is a Go linter that checks that cql queries will not generate run-time errors. 

While, in most cases, queries created using cql are checked at compile time, 
there are still some cases that can generate run-time errors (see :ref:`cql/type_safety:Runtime errors`).

cqllint analyses the Go code written to detect these cases and fix them without the need to execute the query. 
It also adds other detections that would not generate runtime errors but are possible misuses of cql.

.. note::

    At the moment, only the errors cql.ErrFieldModelNotConcerned, cql.ErrFieldIsRepeated, 
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
        conditions.Brand.Name.Is().Eq(conditions.City.Name),
    ).Find()

If we execute this query we will obtain an error of type `cql.ErrFieldModelNotConcerned` with the following message:

.. code-block:: none

    field's model is not concerned by the query (not joined); not concerned model: models.City; operator: Eq; model: models.Brand, field: Name

Now, if we run cqllint we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:3: models.City is not joined by the query

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

ErrAppearanceMustBeSelected
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

To generate this error we must join the same model more than once and not select the appearance number:

.. code-block:: go
    :caption: example.go
    :class: with-errors
    :emphasize-lines: 9
    :linenos:

    _, err := cql.Query[models.Child](
        db,
        conditions.Child.Parent1(
            conditions.Parent1.ParentParent(),
        ),
        conditions.Child.Parent2(
            conditions.Parent2.ParentParent(),
        ),
        conditions.Child.ID.Is().Eq(conditions.ParentParent.ID),
    ).Find()

If we execute this query we will obtain an error of type `cql.ErrAppearanceMustBeSelected` with the following message:

.. code-block:: none

    field's model appears more than once, select which one you want to use with Appearance; model: models.ParentParent; operator: Eq; model: models.Child, field: ID

Now, if we run cqllint we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:9: models.ParentParent appears more than once, select which one you want to use with Appearance

ErrAppearanceOutOfRange
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

To generate this error we must use the Appearance method with a value greater than the number of appearances of a model:

.. code-block:: go
    :caption: example.go
    :class: with-errors
    :emphasize-lines: 4
    :linenos:

    _, err := cql.Query[models.Phone](
        db,
        conditions.Phone.Brand(
            conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Appearance(1)),
        ),
    ).Find()

If we execute this query we will obtain an error of type `cql.ErrAppearanceOutOfRange` with the following message:

.. code-block:: none

    selected appearance is bigger than field's model number of appearances; model: models.Phone; operator: Eq; model: models.Brand, field: Name

Now, if we run cqllint we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:4: selected appearance is bigger than models.Phone's number of appearances

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
        conditions.Brand.Name.Set().Dynamic(conditions.Brand.Name),
    )

If we run cqllint we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:5: conditions.Brand.Name is set to itself

Unnecessary Appearance selection
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

This is the case when the Appearance method is used without being necessary, 
i.e. when the model appears only once:

.. code-block:: go
    :caption: example.go
    :class: with-errors
    :emphasize-lines: 4
    :linenos:

    _, err := cql.Query[models.Phone](
        db,
        conditions.Phone.Brand(
            conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Appearance(0)),
        ),
    ).Find()

If we run cqllint we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:4: Appearance call not necessary, models.Phone appears only once