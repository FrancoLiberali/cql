==============================
Query type safety
==============================

The cql.Query method and its query system provide type safety in multiple cases.

Conditions of the model
-------------------------------

cql will only allow us to add conditions on the model we are querying, 
prohibiting the use of conditions from other models in the wrong place:

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
    ).Find()

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 4
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.Country.Name.Is().Eq(cql.String("Paris")),
    ).Find()

In this case, the compilation error will be:

.. code-block:: none

    cannot use conditions.Country.Name.Is().Eq(cql.String("Paris"))
    (value of interface type condition.WhereCondition[models.Country]) as condition.Condition[models.City]...

Similarly, conditions are checked when making joins:

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Country(
            conditions.Country.Name.Is().Eq(cql.String("France")),
        ),
    ).Find()

.. code-block:: go
    :caption: Incorrect
    :class: with-errors
    :emphasize-lines: 5
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Country(
            conditions.City.Name.Is().Eq(cql.String("France")),
        ),
    ).Find()

Name of an attribute or operator
-------------------------------

Since the conditions are made using the auto-generated code, 
the attributes and methods used on it will only allow us to use attributes and operators that exist:

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
    ).Find()

.. code-block:: go
    :caption: Incorrect
    :class: with-errors
    :emphasize-lines: 4
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Namee.Is().Eq(cql.String("Paris")),
    ).Find()

In this case, the compilation error will be:

.. code-block:: none

    conditions.City.Namee undefined (type conditions.cityConditions has no field or method Namee)

Type of an attribute
-------------------------------

cql not only verifies that the attribute used exists but also verifies that 
the value compared to the attribute is of the correct type:

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
    ).Find()

.. code-block:: go
    :caption: Incorrect
    :class: with-errors
    :emphasize-lines: 4
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.Int64(100)),
    ).Find()

In this case, the compilation error will be:

.. code-block:: none

    cannot use cql.Int64(100) (value of struct type condition.NumericValue[int64]) as condition.ValueOfType[string] 
    value in argument to conditions.City.Name.Is().Eq: condition.NumericValue[int64] does not implement 
    condition.ValueOfType[string] (wrong type for method GetValue)

Type of an attribute (dynamic operator)
-------------------------------

cql also checks that the type of the attributes is correct when using dynamic operators. 
In this case, the type of the two attributes being compared must be the same: 

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Country(
            conditions.Country.Name.Is().Eq(conditions.City.Name),
        ),
    ).Find()

.. code-block:: go
    :caption: Incorrect
    :class: with-errors
    :emphasize-lines: 5
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Country(
            conditions.Country.Name.Is().Eq(conditions.City.Population),
        ),
    ).Find()

In this case, the compilation error will be:

.. code-block:: none

    cannot use conditions.City.Population (variable of type condition.UpdatableField[models.City, int]) as 
    condition.FieldOfType[string] value in argument to conditions.Country.Name.Is().Eq...

Limitations
-------------------------------

Dynamic operators and functions
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

cql.Query is not safe at compile time to determine whether 
the values used in dynamic operators or functions, as in the following examples:

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Country(
            conditions.Country.Name.Is().Eq(conditions.City.Name),
        ),
    ).Find()

.. code-block:: go
    :class: with-errors
    :caption: Not joined model in dynamic function
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
    :caption: Not joined model in dynamic operator
    :emphasize-lines: 5
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(
            conditions.Country.Name,
        ),
    ).Find()

.. code-block:: go
    :class: with-errors
    :caption: Not joined model in dynamic function in dynamic operator
    :emphasize-lines: 6
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(
            conditions.City.Name.Concat(
                conditions.Country.Name,
            ),
        ),
    ).Find()

Which would generate the following error of type cql.ErrFieldModelNotConcerned at runtime:

.. code-block:: none

    field's model is not concerned by the query (not joined); not concerned model: models.Country

.. TODO link a la seccion correcta
These errors can be determined before runtime using :doc:`/cql/cqllint`.

Modifier methods
^^^^^^^^^^^^^^^^^^^^^^^^^^

Similarly, the Descending and Ascending sorting modifier methods are not safe at compile 
time to terminate if the fields used are part of the query, as in the following examples:

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
    ).Descending(
        conditions.City.Name,
    ).Find()

.. code-block:: go
    :class: with-errors
    :caption: Not joined model in order
    :emphasize-lines: 6
    :linenos:

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
    ).Descending(
        conditions.Country.Name,
    ).Find()

Which would generate the following error of type cql.ErrFieldModelNotConcerned at runtime:

.. code-block:: none
    :caption: Result

    field's model is not concerned by the query (not joined); not concerned model: models.Seller; method: Descending

.. TODO link a la seccion correcta
This error can be determined before runtime using :doc:`/cql/cqllint`.

Appearance
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The selection of the :ref:`cql/advanced_query:appearance` can generate two runtime errors:

- cql.ErrAppearanceMustBeSelected: generated when you try to use a model that appears 
  (is joined) more than once in the query without selecting which one you want to use.
- cql.ErrAppearanceOutOfRange: generated when you try select an appearance number (with the Appearance method) 
  greater than the number of appearances of a model.

Both errors can be determined before runtime using :doc:`/cql/cqllint`.

Collections preloads
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Another possible runtime error is cql.ErrOnlyPreloadsAllowed, generated when trying to use conditions 
within a preload of :ref:`cql/advanced_query:collections`:

.. code-block:: go
    :caption: Model

    type Seller struct {
        model.UUIDModel

        Name string

        Company   *Company
        CompanyID *model.UUID // Company HasMany Seller (Company 0..1 -> 0..* Seller)
    }

    type Company struct {
        model.UUIDModel

        Sellers *[]Seller // Company HasMany Seller (Company 0..1 -> 0..* Seller)
    }

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.Company](
        context.Background(),
        db,
        conditions.Company.Sellers.Preload(),
    ).Find()

.. code-block:: go
    :class: with-errors
    :caption: Conditions inside preload
    :emphasize-lines: 5
    :linenos:

    _, err := cql.Query[models.Company](
        context.Background(),
        db,
        conditions.Company.Sellers.Preload(
            conditions.Seller.ID.Is().Eq(cql.String("Franco")),
        ),
    ).Find()

This error has not yet been determined by :doc:`/cql/cqllint`.
