==============================
Select
==============================

If you only want certain attributes from the models as query results, cql provides cql.Select.

This function allows us to take the results of a query received by parameter and select only certain attributes,
both from the main query model and from the joined models.

To perform this selection, the ``Into`` method of the selected field (or aggregation) is used. It receives a
function that, given a pointer to a result, returns the address of the attribute where the selected value will
be stored.

Example 1: Select only one field from the main model

.. code-block:: go
    :caption: Model

    type MyModel struct {
        model.UUIDModel

        Value1 int64
        Value2 string
    }

    type Results struct {
        Value1 int64
    }

.. code-block:: go

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
            conditions.MyModel.Value1.Is().Eq(cql.Int64(4)),
        ),
        conditions.MyModel.Value1.Into(func(result *Results) *int64 { return &result.Value1 }),
    )

Example 2: Select more than one field from the main model

.. code-block:: go
    :caption: Model

    type MyModel struct {
        model.UUIDModel

        Value1 int64
        Value2 string
    }

    type Results struct {
        Value1 int64
        Value2 string
    }

.. code-block:: go

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
            conditions.MyModel.Value1.Is().Eq(cql.Int64(4)),
        ),
        conditions.MyModel.Value1.Into(func(result *Results) *int64 { return &result.Value1 }),
        conditions.MyModel.Value2.Into(func(result *Results) *string { return &result.Value2 }),
    )

Joins
-----------------------

It is possible to select different attributes from the different entities joined in the queries:

.. code-block:: go
    :caption: Model

     type MyOtherModel struct {
        model.UUIDModel

        Name string
    }

    type MyModel struct {
        model.UUIDModel

        Value1 int64

        Related   MyOtherModel
        RelatedID model.UUID
    }

    type Results struct {
        Value1 int64
        Name   string
    }

.. code-block:: go

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
            conditions.MyModel.Value1.Is().Eq(cql.Int64(4)),
            conditions.MyModel.Related(),
        ),
        conditions.MyModel.Value1.Into(func(result *Results) *int64 { return &result.Value1 }),
        conditions.MyOtherModel.Name.Into(func(result *Results) *string { return &result.Name }),
    )

Functions
-----------------------

cql supports applying functions to selected values before retrieving them, either with static values or with other attributes.
For more details on the available functions, please consult :ref:`functions <cql/query:functions>`.

.. note::

    A numeric arithmetic expression follows SQL's type promotion, so its result is bound as ``float64``
    (see :ref:`Type safety <cql/select:type safety>` below). The destination attribute must therefore be a ``float64``.

Example 1: Function with static value

In this case, we will add 2 to the values obtained.

.. code-block:: go
    :caption: Model

    type MyModel struct {
        model.UUIDModel

        Value1 int64
    }

    type Results struct {
        Value1 float64
    }

.. code-block:: go

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
            conditions.MyModel.Value1.Is().Eq(cql.Int64(4)),
        ),
        conditions.MyModel.Value1.Plus(cql.Int64(2)).Into(func(result *Results) *float64 { return &result.Value1 }),
    )

Example 2: Function with other attribute

In this case, we will add two attributes.

.. code-block:: go
    :caption: Model

    type MyModel struct {
        model.UUIDModel

        Value1 int64
        Value2 int64
    }

    type Results struct {
        Value1PlusValue2 float64
    }

.. code-block:: go

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
            conditions.MyModel.Value1.Is().Eq(cql.Int64(4)),
        ),
        conditions.MyModel.Value1.Plus(conditions.MyModel.Value2).Into(func(result *Results) *float64 { return &result.Value1PlusValue2 }),
    )

Aggregations
-----------------------

When selecting, it is also possible to perform aggregations on the values. The available aggregations depend on the type of attribute.

The aggregations available for all types are:

- Count: returns the number of values that are not null.
- Min: returns the minimum value of all values.
- Max: returns the maximum value of all values.

For numeric attributes, the following aggregations are also available:

- Sum: calculates the summation of all values.
- Average: calculates the average (arithmetic mean) of all values.
- And: calculates the bitwise AND of all non-null values (null values are ignored). Not available for: sqlite, sqlserver.
- Or: calculates the bitwise OR of all non-null values (null values are ignored). Not available for: sqlite, sqlserver.

For boolean attributes, the following aggregations are also available:

- All: returns true if all the values are true.
- Any: returns true if at least one value is true.
- None: returns true if all values are false.

.. note::

    Numeric aggregations (Sum, Average, Min, Max, Count, ...) are bound as ``float64``, so the destination
    attribute must be a ``float64``.

Example:

.. code-block:: go
    :caption: Model

    type MyModel struct {
        model.UUIDModel

        Value1 int64
        Value2 int64
    }

    type Results struct {
        Value1Sum float64
        Value2Max float64
    }

.. code-block:: go

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
            conditions.MyModel.Value1.Is().Eq(cql.Int64(4)),
        ),
        conditions.MyModel.Value1.Aggregate().Sum().Into(func(result *Results) *float64 { return &result.Value1Sum }),
        conditions.MyModel.Value2.Aggregate().Max().Into(func(result *Results) *float64 { return &result.Value2Max }),
    )

.. warning::

    Aggregations and non-aggregations cannot be combined within the same select.

Type safety
-----------------------

Select, in addition to inheriting the type safety of Query, adds a new layer of type safety to selections.

Selection type
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

``Into`` ensures at compile time that the selected value is stored in an attribute of the correct type: the
function you pass must return the address of an attribute whose type matches the selected value.

A plain column binds its own Go type (an ``int64`` column into an ``int64`` attribute, a ``string`` column into a
``string`` attribute, and so on). Only numeric expressions (``Plus``, ``Minus``, ...) and numeric aggregations widen
to ``float64``; for those the destination attribute must be a ``float64``.

.. code-block:: go
    :caption: Model

    type MyModel struct {
        model.UUIDModel

        Value1 int64
        Value2 string
    }

    type Results struct {
        Value1 int64
    }

.. code-block:: go
    :caption: Correct
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ),
        conditions.MyModel.Value1.Into(func(result *Results) *int64 { return &result.Value1 }),
    )

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 6
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ),
        conditions.MyModel.Value2.Into(func(result *Results) *int64 { return &result.Value1 }),
    )

In this case, ``Value2`` is a ``string`` column, so ``Into`` expects a function returning ``*string``. The
compilation error will be:

.. code-block:: none

    type func(result *Results) *int64 of (func(result *Results) *int64 literal)
    does not match inferred type func(*Results) *string for func(*TResults) *string

Into functions
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

As in cql.Query, the functions applied to the selected values are type-safe.

.. code-block:: go
    :caption: Model

    type MyModel struct {
        model.UUIDModel

        Value1 int64
        Value2 string
        Value3 int64
    }

    type Results struct {
        Value1 float64
    }

.. code-block:: go
    :caption: Correct
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ),
        conditions.MyModel.Value1.Plus(conditions.MyModel.Value3).Into(func(result *Results) *float64 { return &result.Value1 }),
    )

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 6
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ),
        conditions.MyModel.Value1.Plus(conditions.MyModel.Value2).Into(func(result *Results) *float64 { return &result.Value1 }),
    )

In this case, the compilation error will be:

.. code-block:: none

    cannot use conditions.MyModel.Value2 (variable of struct type condition.StringField[MyModel])
    as condition.ValueOfType[float64] value in argument to conditions.MyModel.Value1.Plus: condition.StringField[MyModel]
    does not implement condition.ValueOfType[float64] (wrong type for method GetValue)

Type safety limitations and cqllint
------------------------------------------------

Once again, similar to cql.Query, ``Into`` is not safe at compile time to determine whether
the values selected or used in functions are joined in the query, as in the following examples:

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 6
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ),
        conditions.MyOtherModel.Value1.Into(func(result *Results) *float64 { return &result.Value1 }),
    )

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 6
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ),
        conditions.MyModel.Value1.Plus(conditions.MyOtherModel.Value2).Into(func(result *Results) *float64 { return &result.Value1 }),
    )

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 6
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ),
        conditions.MyOtherModel.Value1.Aggregate().Sum().Into(func(result *Results) *float64 { return &result.Value1 }),
    )

Which would generate the following error at runtime:

.. code-block:: none

    field's model is not concerned by the query (not joined); not concerned model: models.MyOtherModel

Now, if we run :doc:`/cql/cqllint` we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:6: models.MyOtherModel is not joined by the query
