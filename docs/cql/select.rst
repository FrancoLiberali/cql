==============================
Select
==============================

If you only want certain attributes from the models as query results, cql provides cql.Select.

This function allows us to take the results of a query received by parameter and select only certain attributes, 
both from the main query model and from the joined models.

To perform this selection, the cql.ValueInto function is used, which receives the field to be selected and 
a function to save that field in the results list.

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
        cql.ValueInto(conditions.MyModel.Value1, func(value float64, result *Results) {
            result.Value1 = int64(value)
        }),
    )

Example 2: Select more than one field from the main model

.. code-block:: go
    :caption: Model

    type MyModel struct {
        model.UUIDModel

        Value1 int64
        Value2 String
    }

    type Results struct {
        Value1 int64
        Value2 String
    }

.. code-block:: go

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
            conditions.MyModel.Value1.Is().Eq(cql.Int64(4)),
        ),
        cql.ValueInto(conditions.MyModel.Value1, func(value float64, result *Results) {
            result.Value1 = int64(value)
        }),
        cql.ValueInto(conditions.MyModel.Value2, func(value string, result *Results) {
            result.Value2 = value
        }),
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
        Name String
    }

.. code-block:: go

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
            conditions.MyModel.Value1.Is().Eq(cql.Int64(4)),
            conditions.MyModel.Related(),
        ),
        cql.ValueInto(conditions.MyModel.Value1, func(value float64, result *Results) {
            result.Value1 = int64(value)
        }),
        cql.ValueInto(conditions.MyOtherModel.Name, func(value string, result *Results) {
            result.Name = value
        }),
    )

Functions
-----------------------

cql supports applying functions to selected values before retrieving them, either with static values or with other attributes.
For more details on the available functions, please consult :ref:`functions <cql/query:functions>`.

Example 1: Function with static value

In this case, we will add 2 to the values obtained.

.. code-block:: go
    :caption: Model

    type MyModel struct {
        model.UUIDModel

        Value1 int64
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
        cql.ValueInto(conditions.MyModel.Value1.Plus(cql.Int64(2)), func(value float64, result *Results) {
            result.Value1 = int64(value)
        }),
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
        Value1PlusValue2 int64
    }

.. code-block:: go

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
            conditions.MyModel.Value1.Is().Eq(cql.Int64(4)),
        ),
        cql.ValueInto(conditions.MyModel.Value1.Plus(conditions.MyModel.Value2), func(value float64, result *Results) {
            result.Value1PlusValue2 = int64(value)
        }),
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

Example:

.. code-block:: go
    :caption: Model

    type MyModel struct {
        model.UUIDModel

        Value1 int64
        Value2 int64
    }

    type Results struct {
        Value1Sum int64
        Value2Max int64
    }

.. code-block:: go

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
            conditions.MyModel.Value1.Is().Eq(cql.Int64(4)),
        ),
        cql.ValueInto(conditions.MyModel.Value1.Aggregate().Sum(), func(value float64, result *Results) {
            result.Value1Sum = int64(value)
        }),
        cql.ValueInto(conditions.MyModel.Value2.Aggregate().Max(), func(value float64, result *Results) {
            result.Value2Max = int64(value)
        }),
    )

Type safety
-----------------------

Select, in addition to inheriting the type safety of Query, adds a new layer of type safety to selections.

Selection type
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

cql.ValueInto ensures that the type of the selected attribute is stored in a result of the correct type.

Note that selections of all numeric types are of type float64, which must then be cast to the desired type if it is different.

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
        cql.ValueInto(conditions.MyModel.Value1, func(value float64, result *Results) {
            result.Value1 = int64(value)
        }),
    )

.. code-block:: go
    :class: with-errors
    :caption: Incorrect 1
    :emphasize-lines: 6
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ),
        cql.ValueInto(conditions.MyModel.Value2, func(value float64, result *Results) {
            result.Value1 = value
        }),
    )

In this case, the compilation error will be:

.. code-block:: none

    in call to cql.ValueInto, type func(value float64, result *Results) of (func(value float64, result *ResultInt) literal)
    does not match inferred type func(string, *TResults) for func(TValue, *TResults)

.. code-block:: go
    :class: with-errors
    :caption: Incorrect 2
    :emphasize-lines: 7
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ),
        cql.ValueInto(conditions.MyModel.Value2, func(value string, result *Results) {
            result.Value1 = value
        }),
    )

In this case, the compilation error will be:

.. code-block:: none

    cannot use value (variable of type string) as int64 value in assignment

Functions
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
        cql.ValueInto(conditions.MyModel.Value1.Plus(conditions.MyModel.Value3), func(value float64, result *Results) {
            result.Value1 = int64(value)
        }),
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
        cql.ValueInto(conditions.MyModel.Value1.Plus(conditions.MyModel.Value2), func(value float64, result *Results) {
            result.Value1 = value
        }),
    )

In this case, the compilation error will be:

.. code-block:: none

    cannot use conditions.MyModel.Value2 (variable of struct type condition.StringField[MyModel]) 
    as condition.ValueOfType[float64] value in argument to conditions.MyModel.Value1.Plus: condition.StringField[MyModel] 
    does not implement condition.ValueOfType[float64] (wrong type for method GetValue)

Limitations
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Once again, similar to cql.Query, cql.ValueInto is not safe at compile time to determine whether 
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
        cql.ValueInto(conditions.MyOtherModel.Value1, func(value float64, result *Results) {
            result.Value1 = value
        }),
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
        cql.ValueInto(conditions.MyModel.Value1.Plus(conditions.MyOtherModel.Value2), func(value float64, result *Results) {
            result.Value1 = value
        }),
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
        cql.ValueInto(conditions.MyOtherModel.Value1.Aggregate().Sum(), func(value float64, result *Results) {
            result.Value1 = value
        }),
    )

Which would generate the following error at runtime:

.. code-block:: none

    field's model is not concerned by the query (not joined); not concerned model: models.MyOtherModel

.. TODO link a la seccion correcta
These errors can be determined before runtime using :doc:`/cql/cqllint`.