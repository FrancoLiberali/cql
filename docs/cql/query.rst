==============================
Query
==============================

Read (query) operations are provided by cql via its compiled query system.

Query creation
-----------------------

To create a query you must use the cql.Query[models.MyModel] method,
where models.MyModel is the model you expect this query to answer. 
This function takes as parameters the :ref:`transaction <cql/query:transactions>` 
on which to execute the query and the :ref:`cql/query:conditions`.

Transactions
--------------------

To execute transactions, cql provides the function cql.Transaction. 
The function passed by parameter will be executed inside a gorm transaction 
(for more information visit https://gorm.io/docs/transactions.html). 
Using this method will also allow the transaction execution time to be logged.

Query methods
------------------------

The object obtained using `cql.Query` has different methods that 
will allow you to obtain the results of the query:

Modifier methods
^^^^^^^^^^^^^^^^^^^^^^^^^^

Modifier methods are those that modify the query in a certain way, affecting the results obtained:
- Limit: specifies the number of models to be retrieved.
- Offset: specifies the number of models to skip before starting to return the results.
- Ascending: specifies an ascending order when retrieving models.
- Descending: specifies a descending order when retrieving models from database.

Finishing methods
^^^^^^^^^^^^^^^^^^^^^^^

Finishing methods are those that cause the query to be executed and the result(s) of the query to be returned:

- First: finds the first model ordered by primary key.
- Take: finds the first model returned by the database in no specified order.
- Last: finds the last model ordered by primary key.
- FindOne: finds the only one model that matches given conditions or returns error if 0 or more than 1 are found.
- Find: finds list of models that meet the conditions.

Conditions
------------------------

The set of conditions that are received by the `cql.Query` method 
form the cql compiled query system. 
It is so named because the conditions will verify at compile time that the query to be executed is correct.

These conditions are objects of type Condition that contain the 
necessary information to perform the queries in a safe way. 
They are generated from the definition of your models using cql-cli.

Conditions generation
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The generation of conditions is done with cql-cli. For this, we need to install cql-cli:

.. code-block:: bash

    go install github.com/FrancoLiberali/cql/cql-cli

Then, inside our project we will have to create a package called conditions 
(or another name if you wish) and inside it a file with the following content:

.. code-block:: go

    package conditions

    //go:generate cql-cli gen conditions ../models_path_1 ../models_path_2

where ../models_path_1 ../models_path_2 are the relative paths between the package conditions 
and the packages containing the definition of your models (can be only one).

Now, from the root of your project you can execute:

.. code-block:: bash

  go generate ./...

and the conditions for each of your models will be created in the conditions package.

Use of the conditions
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

After performing the conditions generation, 
your conditions package will have a replica of your models package, 
i.e. if, for example, the type models.MyModel is part of your models, 
the variable conditions.MyModel will be in the conditions package. 
This variable is called the condition model and it has:

- An attribute for each attribute of your original model with the same name 
  (if models.MyModel.Name exists, then conditions.MyModel.Name is generated), 
  that allows to use that attribute in queries creating operations for that attribute in your queries.
- A method for each relation of your original model with the same name 
  (if models.MyModel.MyOtherModel exists, then conditions.MyModel.MyOtherModel() is generated), 
  which will allow you to perform joins in your queries.
- Methods for :doc:`/cql/preloading`.

Then, combining these conditions, the Connection Conditions (cql.And, cql.Or, cql.Not) 
you will be able to make all the queries you need in a safe way.

Examples
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**Filter by an attribute**

In this example we query all MyModel that has "a_string" in the Name attribute.

.. code-block:: go

    type MyModel struct {
        model.UUIDModel

        Name string
    }

    myModels, err := cql.Query[MyModel](
        gormDB,
        conditions.MyModel.Name.Is().Eq("a_string"),
    ).Find()

**Filter by an attribute of a related model**

In this example we query all MyModels whose related MyOtherModel has "a_string" in its Name attribute.

.. code-block:: go

    type MyOtherModel struct {
        model.UUIDModel

        Name string
    }

    type MyModel struct {
        model.UUIDModel

        Related   MyOtherModel
        RelatedID model.UUID
    }

    myModels, err := cql.Query[MyModel](
        gormDB,
        conditions.MyModel.Related(
            conditions.MyOtherModel.Name.Is().Eq("a_string"),
        ),
    ).Find()

**Multiple conditions**

In this example we query all MyModels that has a 4 in the Code attribute and 
whose related MyOtherModel has "a_string" in its Name attribute.

.. code-block:: go

    type MyOtherModel struct {
        model.UUIDModel

        Name string
    }

    type MyModel struct {
        model.UUIDModel

        Code int

        Related   MyOtherModel
        RelatedID model.UUID
    }

    myModels, err := cql.Query[MyModel](
        gormDB,
        conditions.MyModel.Code.Is().Eq(4),
        conditions.MyModel.Related(
            conditions.MyOtherModel.Name.Is().Eq("a_string"),
        ),
    ).Find()

Operators
------------------------

The different operators to use inside your queries are defined by 
the methods of the FieldIs type, which is returned when calling the Is() method. 
Below you will find the complete list of available operators:

- Eq(value): Equal to
- NotEq(value): Not equal to
- Lt(value): Less than
- LtOrEq(value): Less than or equal to
- Gt(value): Greater than
- GtOrEq(value): Greater than or equal to
- Null()
- NotNull()
- Between(v1, v2): Equivalent to v1 < attribute < v2
- NotBetween(v1, v2): Equivalent to NOT (v1 < attribute < v2)
- Distinct(value)
- NotDistinct(value)
- In(values)
- NotIn(values)

For boolean attributes:

- True()
- NotTrue()
- False()
- NotFalse()
- Unknown(): unknown is null for booleans
- NotUnknown(): unknown is null for booleans

For string attributes:

- Like(pattern)

In addition to these, cql gives the possibility to use operators 
that are only supported by a certain database (outside the standard). 
For doing it, you must use the Custom method and give the operator as argument, for example:

.. code-block:: go

    conditions.MyModel.Code.Is().Custom(psql.ILike("_a%")),

These operators can be found in <https://pkg.go.dev/github.com/FrancoLiberali/cql/mysql>, 
<https://pkg.go.dev/github.com/FrancoLiberali/cql/sqlserver>, 
<https://pkg.go.dev/github.com/FrancoLiberali/cql/psql> 
and <https://pkg.go.dev/github.com/FrancoLiberali/cql/sqlite>. 

You can also define your own operators following the condition.Operator interface.
