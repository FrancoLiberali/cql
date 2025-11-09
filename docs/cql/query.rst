==============================
Query
==============================

Read (query) operations are provided by cql via its compiled query system.

Query creation
-----------------------

To create a query you must use the cql.Query[models.MyModel] method,
where models.MyModel is the model you expect this query to answer. 
This function takes as parameters the db or the :ref:`transaction <cql/query:transactions>` 
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

- Count: returns the amount of models that fulfill the conditions.
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
They are generated from the definition of your models using :ref:`cql-gen <cql/cqlgen:Conditions generation>`.

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
        context.Background(),
        db,
        conditions.MyModel.Name.Is().Eq(cql.String("a_string")),
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
        context.Background(),
        db,
        conditions.MyModel.Related(
            conditions.MyOtherModel.Name.Is().Eq(cql.String("a_string")),
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
        context.Background(),
        db,
        conditions.MyModel.Code.Is().Eq(cql.Int64(4)),
        conditions.MyModel.Related(
            conditions.MyOtherModel.Name.Is().Eq(cql.String("a_string")),
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

Static values
------------------------

As can be seen in the previous examples, operators can receive another value to perform the comparison.
These values can be static or :ref:`dynamic <cql/advanced_query:Dynamic operators>`.

For static values, it is necessary to define their type using one of the functions provided by cql:

- Int(value int)
- Int8(value int8)
- Int16(value int16)
- Int32(value int32)
- Int64(value int64)
- UInt(value uint)
- UInt8(value uint8)
- UInt16(value uint16)
- UInt32(value uint32)
- UInt64(value uint64)
- Float32(value float32)
- Float64(value float64)
- Bool(value bool)
- String(value string)
- ByteArray(value []byte)
- Time(value time.Time)
- UUID(value model.UUID)

This ensures that operations are only performed between compatible types.

Custom types
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

In addition to these static values, it is possible to define your own types to be supported by CQL.

For this, the type must implement the ValueOfType[T any] interface, which consists of two methods:

- ToSQL(query *CQLQuery) (string, []any, error):
    Allows to define how the type is translated to SQL,
    allowing you to define the SQL statement to be used,
    the parameters for this statement, and an error.
- GetValue() T: Allows to define the type with which this type is comparable.


Functions
------------------------

It is also possible to apply functions on the values to be used in the operations.

For example, we can query all MyModels for which half of their Code is less than 10.

.. code-block:: go

    type MyModel struct {
        model.UUIDModel

        Code int
    }

    myModels, err := cql.Query[MyModel](
        context.Background(),
        db,
        conditions.MyModel.Code.Divided(cql.Int64(2)).Is().Lt(cql.Int64(10)),
    ).Find()

The functions that are applicable depend on the type of attribute.

For numeric attributes:

- Plus(other)
- Minus(other)
- Times(other)
- Divided(other)
- Modulo(other)
- Power(other)
- SquareRoot()
- Absolute()
- And(other)
- Or(other)
- Xor(other)
- Not()
- ShiftLeft(other)
- ShiftRight(other)

For string values:

- Concat(other)
