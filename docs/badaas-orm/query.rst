==============================
Query
==============================

Create, Save and Delete methods are done directly with gormDB object using the corresponding methods. 
For details visit 
<https://gorm.io/docs/create.html>, <https://gorm.io/docs/update.html> and <https://gorm.io/docs/delete.html>. 
On the other hand, read (query) operations are provided by badaas-orm via its compilable query system.

Query creation
-----------------------

To create a query you must use the orm.NewQuery[models.MyModel] method,
where models.MyModel is the model you expect this query to answer. 
This function takes as parameters the :ref:`transaction <badaas-orm/query:transactions>` 
on which to execute the query and the query :ref:`badaas-orm/query:conditions`.

Transactions
--------------------

To execute transactions badaas-orm provides the function orm.Transaction. 
The function passed by parameter will be executed inside a gorm transaction 
(for more information visit https://gorm.io/docs/transactions.html). 
Using this method will also allow the transaction execution time to be logged.

Query methods
------------------------

The `orm.Query` object obtained using `orm.NewQuery` has different methods that 
will allow you to obtain the results of the query:

- FindOne: will allow you to obtain the only one model that meets the conditions received by parameter
  or an error will be returned if none or more than one model comply with them.
- Find: will allow you to obtain the list of models that meet the conditions received by parameter.

Conditions
------------------------

The set of conditions that are received by the `orm.NewQuery` method 
form the badaas-orm compilable query system. 
It is so named because the conditions will verify at compile time that the query to be executed is correct.

These conditions are objects of type Condition that contain the 
necessary information to perform the queries in a safe way. 
They are generated from the definition of your models using badaas-cli.

Conditions generation
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The generation of conditions is done with badaas-cli. For this, we need to install badaas-cli:

.. code-block:: bash

    go install github.com/ditrit/badaas-cli

Then, inside our project we will have to create a package called conditions 
(or another name if you wish) and inside it a file with the following content:

.. code-block:: go

    package conditions

    //go:generate badaas-cli gen conditions ../models_path_1 ../models_path_2

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
  of type FieldIdentifier that allows to use that attribute in queries 
  (for :ref:`dynamic conditions <badaas-orm/advanced_query:dynamic operators>` for example).
- A method for each attribute of your original model with the same name + Is 
  (if models.MyModel.Name exists, then conditions.MyModel.NameIs() is generated), 
  which will allow you to create operations for that attribute in your queries.
- A method for each relation of your original model with the same name 
  (if models.MyModel.MyOtherModel exists, then conditions.MyModel.MyOtherModel() is generated), 
  which will allow you to perform joins in your queries.
- Methods for :doc:`/badaas-orm/preloading`.

Then, combining these conditions, the Connection Conditions (orm.And, orm.Or, orm.Not) 
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

    myModels, err := orm.NewQuery[MyModel](
        gormDB,
        conditions.MyModel.NameIs().Eq("a_string"),
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

    myModels, err := orm.NewQuery[MyModel](
        gormDB,
        conditions.MyModel.Related(
            conditions.MyOtherModel.NameIs().Eq("a_string"),
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

    myModels, err := orm.NewQuery[MyModel](
        gormDB,
        conditions.MyModel.CodeIs().Eq(4),
        conditions.MyModel.Related(
            conditions.MyOtherModel.NameIs().Eq("a_string"),
        ),
    ).Find()

Operators
------------------------

The different operators to use inside your queries are defined by 
the methods of the FieldIs type, which is returned when using, for example, 
the conditions.MyModel.CodeIs() method. 
Below you will find the complete list of available operators:

- Eq(value): EqualTo
- NotEq(value): NotEqualTo
- Lt(value): LessThan
- LtOrEq(value): LessThanOrEqualTo
- Gt(value): GreaterThan
- GtOrEq(value): GreaterThanOrEqualTo
- Null()
- NotNull()
- Between(v1, v2): Equivalent to v1 < attribute < v2
- NotBetween(v1, v2): Equivalent to NOT (v1 < attribute < v2)
- True() (Not supported by: sqlserver)
- NotTrue() (Not supported by: sqlserver)
- False() (Not supported by: sqlserver)
- NotFalse() (Not supported by: sqlserver)
- Unknown() (Not supported by: sqlserver, sqlite)
- NotUnknown() (Not supported by: sqlserver, sqlite)
- Distinct(value) (Not supported by: mysql)
- NotDistinct(value) (Not supported by: mysql)
- Like(pattern)
- In(values)
- NotIn(values)

In addition to these, badaas-orm gives the possibility to use operators 
that are only supported by a certain database (outside the standard). 
These operators can be found in <https://pkg.go.dev/github.com/ditrit/badaas/orm/mysql>, 
<https://pkg.go.dev/github.com/ditrit/badaas/orm/sqlserver>, 
<https://pkg.go.dev/github.com/ditrit/badaas/orm/psql> 
and <https://pkg.go.dev/github.com/ditrit/badaas/orm/sqlite>. 
To use them, use the Custom method of FieldIs type.