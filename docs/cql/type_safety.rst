==============================
Type safety
==============================

Compile time safety
-------------------------------

One of the most important features of the CQL is

.. code-block:: none

    Is compile-time safe:
        queries are validated at compile time to avoid errors 
        such as comparing attributes that are of different types, 
        trying to use attributes or navigate relationships that do not exist, 
        using information from tables that are not included in the query, etc.; 
        ensuring that a runtime error will not be raised.

While there are other libraries that provide an API type safety 
(`gorm-gen <https://gorm.io/gen/>`_, `jooq <https://www.jooq.org/>`_ (Java), 
`diesel <https://diesel.rs/>`_ (Rust)), CQL is the only one that allows us to be sure 
that the generated query is correct, (almost) avoiding runtime errors 
(to understand why "almost" see :ref:`cql/type_safety:runtime errors`)

Conditions of the model
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

cql will only allow us to add conditions on the model we are querying, 
prohibiting the use of conditions from other models in the wrong place:

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
    ).Find()

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 3
    :linenos:

    _, err := cql.Query[models.City](
        db,
        conditions.Country.Name.Is().Eq("Paris"),
    ).Find()

In this case, the compilation error will be:

.. code-block:: none

    cannot use conditions.Country.Name.Is().Eq("Paris")
    (value of type condition.WhereCondition[models.Country]) as condition.Condition[models.City]...

Similarly, conditions are checked when making joins:

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        db,
        conditions.City.Country(
            conditions.Country.Name.Is().Eq("France"),
        ),
    ).Find()

.. code-block:: go
    :caption: Incorrect
    :class: with-errors
    :emphasize-lines: 4
    :linenos:

    _, err := cql.Query[models.City](
        db,
        conditions.City.Country(
            conditions.City.Name.Is().Eq("France"),
        ),
    ).Find()

Name of an attribute or operator
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
Since the conditions are made using the auto-generated code, 
the attributes and methods used on it will only allow us to use attributes and operators that exist:


.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
    ).Find()

.. code-block:: go
    :caption: Incorrect
    :class: with-errors
    :emphasize-lines: 3
    :linenos:

    _, err := cql.Query[models.City](
        db,
        conditions.City.Namee.Is().Eq("Paris"),
    ).Find()

In this case, the compilation error will be:

.. code-block:: none

    conditions.City.Namee undefined (type conditions.cityConditions has no field or method Namee)

Type of an attribute
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

cql not only verifies that the attribute used exists but also verifies that 
the value compared to the attribute is of the correct type:

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
    ).Find()

.. code-block:: go
    :caption: Incorrect
    :class: with-errors
    :emphasize-lines: 3
    :linenos:

    _, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq(100),
    ).Find()

In this case, the compilation error will be:

.. code-block:: none

    cannot use 100 (untyped int constant) as string value in argument to conditions.City.Name.Is().Eq

Type of an attribute (dynamic operator)
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

cql also checks that the type of the attributes is correct when using dynamic operators. 
In this case, the type of the two attributes being compared must be the same: 

.. code-block:: go
    :caption: Correct
    :linenos:

    _, err := cql.Query[models.City](
        db,
        conditions.City.Country(
            conditions.Country.Name.IsDynamic().Eq(conditions.City.Name.Value()),
        ),
    ).Find()

.. code-block:: go
    :caption: Incorrect
    :class: with-errors
    :emphasize-lines: 4
    :linenos:

    _, err := cql.Query[models.City](
        db,
        conditions.City.Country(
            conditions.Country.Name.IsDynamic().Eq(conditions.City.Population.Value()),
        ),
    ).Find()

In this case, the compilation error will be:

.. code-block:: none

    cannot use conditions.City.Population (variable of type condition.UpdatableField[models.City, int]) as condition.FieldOfType[string] value in argument to conditions.Country.Name.IsDynamic().Eq...

Runtime errors
-------------------------------

Although all the above checks are at compile-time, 
there are still some possible cases that generate the following run-time errors:

- cql.ErrFieldModelNotConcerned **(1)**: generated when trying to use a model that is not related 
  to the rest of the query (not joined).
- cql.ErrJoinMustBeSelected: generated when you try to use a model that is included 
  (joined) more than once in the query without selecting which one you want to use (see :ref:`cql/advanced_query:select join`).
- cql.ErrFieldIsRepeated: generated when a field is repeated inside a Set call (see :doc:`/cql/update`).
- cql.ErrOnlyPreloadsAllowed: generated when trying to use conditions within a preload of collections (see :ref:`cql/advanced_query:collections`).
- cql.ErrUnsupportedByDatabase: generated when an attempt is made to use a method or function that is not supported by the database engine used.
- cql.ErrOrderByMustBeCalled: generated when in MySQL you try to do a delete/update with Limit but without using OrderBy.

.. note::

    **(1)** errors avoided with :doc:`/cql/cqllint`.

However, these errors are discovered by CQL before the query is executed. 
In addition, CQL will add to the error clear information about the problem so that it is easy to fix, for example:

.. code-block:: go
    :caption: Query
    :class: with-errors
    :emphasize-lines: 4
    :linenos:

    _, err := cql.Query[models.Product](
        ts.db,
        conditions.Product.Int.Is().Eq(1),
    ).Descending(conditions.Seller.ID).Find()

    fmt.Println(err)

.. code-block:: none
    :caption: Result

    field's model is not concerned by the query (not joined); not concerned model: models.Seller; method: Descending