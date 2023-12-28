==============================
Compile-type safety
==============================

One of the most important features of the CQL is::

    Is compile-time safe:
        its query system is validated at compile time to avoid errors 
        such as comparing attributes that are of different types, 
        trying to use attributes or navigate relationships that do not exist, 
        using information from tables that are not included in the query, etc.

This allows us to be sure that the code written with CQL will generate a correct query, avoiding runtime errors.

Conditions of the model
-------------------------------

cql will only allow us to add conditions on the model we are querying, 
prohibiting the use of conditions from other models in the wrong place:

.. code-block:: go
    :caption: Correct

    _, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
    ).Find()

.. code-block:: go
    :caption: Incorrect

    _, err := cql.Query[models.City](
        db,
        :compilation-error:`conditions.Country.Name.Is().Eq("Paris"),`
    ).Find()

In this case, the compilation error will be:

.. code-block:: bash

    cannot use conditions.Country.Name.Is().Eq("Paris")
    (value of type condition.WhereCondition[models.Country]) as condition.Condition[models.City]...

Similarly, conditions are checked when making joins:

.. code-block:: go
    :caption: Correct

    _, err := cql.Query[models.City](
        db,
        conditions.City.Country(
            conditions.Country.Name.Is().Eq("France"),
        ),
    ).Find()

.. code-block:: go
    :caption: Incorrect

    _, err := cql.Query[models.City](
        db,
        conditions.City.Country(
            :compilation-error:`conditions.City.Name.Is().Eq("France"),`
        ),
    ).Find()

Name of an attribute or operator
--------------------------------------

Since the conditions are made using the auto-generated code, 
the attributes and methods used on it will only allow us to use attributes and operators that exist:


.. code-block:: go
    :caption: Correct

    _, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
    ).Find()

.. code-block:: go
    :caption: Incorrect

    _, err := cql.Query[models.City](
        db,
        conditions.City.:compilation-error:`Namee`.Is().Eq("Paris"),
    ).Find()

In this case, the compilation error will be:

.. code-block:: bash

    conditions.City.Namee undefined (type conditions.cityConditions has no field or method Namee)

Type of an attribute
--------------------------------------

cql not only verifies that the attribute used exists but also verifies that 
the value compared to the attribute is of the correct type:

.. code-block:: go
    :caption: Correct

    _, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
    ).Find()

.. code-block:: go
    :caption: Incorrect

    _, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq(:compilation-error:`100`),
    ).Find()

In this case, the compilation error will be:

.. code-block:: bash

    cannot use 100 (untyped int constant) as string value in argument to conditions.City.Name.Is().Eq
