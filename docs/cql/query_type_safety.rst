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

Type safety limitations and cqllint
------------------------------------------------

Dynamic operators and functions
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

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

    field's model is not concerned by the query (not joined); not concerned model: 
    models.Country; operator: Eq; model: models.City, field: Name

Now, if we run :doc:`/cql/cqllint` we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:5: models.Country is not joined by the query

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

    field's model is not concerned by the query (not joined); not concerned model: models.Seller; method: Descending

Now, if we run :doc:`/cql/cqllint` we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:6: models.Country is not joined by the query

Group by
^^^^^^^^^^^^^^^^^^^^^^^^^^

Similarly, the GroupBy and Having methods are not safe at compile 
time to terminate if the fields used are part of the query, as in the following examples:

.. code-block:: go
    :caption: Correct
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ).GroupBy(
            conditions.MyModel.Name,
        ),
        cql.ValueInto(conditions.MyModel.Name, func(value string, result *Result) {
            result.Name = value
        }),
        cql.ValueInto(conditions.MyModel.Status.Aggregate().Sum(), func(value float64, result *Result) {
            result.SumStatus = int(value)
        }),
    )

.. code-block:: go
    :class: with-errors
    :caption: Not joined model in group by
    :emphasize-lines: 6
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ).GroupBy(
            conditions.MyOtherModel.Name,
        ),
        cql.ValueInto(conditions.MyModel.Name, func(value string, result *Result) {
            result.Name = value
        }),
        cql.ValueInto(conditions.MyModel.Status.Aggregate().Sum(), func(value float64, result *Result) {
            result.SumStatus = int(value)
        }),
    )

Which would generate the following error of type cql.ErrFieldModelNotConcerned at runtime:

.. code-block:: none

    field's model is not concerned by the query (not joined); not concerned model: MyOtherModel; method: GroupBy

.. code-block:: go
    :class: with-errors
    :caption: Not joined model in having
    :emphasize-lines: 8
    :linenos:

    results, err := cql.Select(
        cql.Query[MyModel](
            context.Background(),
            db,
        ).GroupBy(
            conditions.MyModel.Name,
        ).Having(
            conditions.MyOtherModel.Status.Aggregate().Count().Gt(cql.Int(2)),
        ),
        cql.ValueInto(conditions.MyModel.Name, func(value string, result *Result) {
            result.Name = value
        }),
        cql.ValueInto(conditions.MyModel.Status.Aggregate().Sum(), func(value float64, result *Result) {
            result.SumStatus = int(value)
        }),
    )

Which would generate the following error of type cql.ErrFieldModelNotConcerned at runtime:

.. code-block:: none

    field's model is not concerned by the query (not joined); not concerned model: MyOtherModel; method: Having

Now, if we run :doc:`/cql/cqllint` we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:6: MyOtherModel is not joined by the query

Appearance
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The selection of the :ref:`cql/advanced_query:appearance` can generate two runtime errors:

- cql.ErrAppearanceMustBeSelected: generated when you try to use a model that appears 
  (is joined) more than once in the query without selecting which one you want to use.
- cql.ErrAppearanceOutOfRange: generated when you try select an appearance number (with the Appearance method) 
  greater than the number of appearances of a model.

Both errors can be determined before runtime using :doc:`/cql/cqllint`.

ErrAppearanceMustBeSelected
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

To generate this error we must join the same model more than once and not select the appearance number:

.. code-block:: go
    :caption: example.go
    :class: with-errors
    :emphasize-lines: 10
    :linenos:

    _, err := cql.Query[models.Child](
        context.Background(),
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

Now, if we run :doc:`/cql/cqllint` we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:10: models.ParentParent appears more than once, select which one you want to use with Appearance

ErrAppearanceOutOfRange
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

To generate this error we must use the Appearance method with a value greater than the number of appearances of a model:

.. code-block:: go
    :caption: example.go
    :class: with-errors
    :emphasize-lines: 5
    :linenos:

    _, err := cql.Query[models.Phone](
        context.Background(),
        db,
        conditions.Phone.Brand(
            conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Appearance(1)),
        ),
    ).Find()

If we execute this query we will obtain an error of type `cql.ErrAppearanceOutOfRange` with the following message:

.. code-block:: none

    selected appearance is bigger than field's model number of appearances; model: models.Phone; operator: Eq; model: models.Brand, field: Name

Now, if we run :doc:`/cql/cqllint` we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:5: selected appearance is bigger than models.Phone's number of appearances

Unnecessary Appearance selection
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This is the case when the Appearance method is used without being necessary, 
i.e. when the model appears only once:

.. code-block:: go
    :caption: example.go
    :class: with-errors
    :emphasize-lines: 5
    :linenos:

    _, err := cql.Query[models.Phone](
        context.Background(),
        db,
        conditions.Phone.Brand(
            conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Appearance(0)),
        ),
    ).Find()

If we run :doc:`/cql/cqllint` we will see the following report:

.. code-block:: none

    $ cqllint ./...
    example.go:5: Appearance call not necessary, models.Phone appears only once

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
