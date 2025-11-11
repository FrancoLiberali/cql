==============================
Update
==============================

While update operations can still be performed using gorm's Save method 
(see `gorm documentation <https://gorm.io/docs/update.html>`_), 
this is useful only if the model(s) to be updated have already been loaded from the database.

On the contrary, cql's Update method allows the update of all the models that meet 
the conditions entered without the need to load the information (via the direct execution of an UPDATE statement).

Update methods
------------------------

Update operations are divided into two parts: the Update method and the Set method. 
In the first one, we must define the conditions that will determine which models will be updated. 
Here, the whole system of compilable queries is valid (for details visit :doc:`/cql/query`). 
In the second one, we define the updates to be performed.

The object obtained using `cql.Update` has different methods that 
will allow you to modify the query:

Modifier methods
^^^^^^^^^^^^^^^^^^^^^^^^^^

Modifier methods are those that modify the query in a certain way, affecting the models updated:

- Limit: specifies the number of models to be updated.
- Ascending: specifies an ascending order when updating models.
- Descending: specifies a descending order when updating models.
- Returning: specifies that the updated models must be fetched from the database after being updated (not supported by MySQL). Preload of related data is also possible (not supported by SQLite). 

Finishing methods
^^^^^^^^^^^^^^^^^^^^^^^

Finishing methods are those that cause the query to be executed:

- Set: defines the updates to be performed.
- SetMultiple: (only supported by MySQL) allows updates to be made to different tables at the same time.

Example
------------------------

.. code-block:: go

    type MyModel struct {
        model.UUIDModel

        Name string
    }

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.Name.Is().Eq(cql.String("a_string")),
    ).Set(
        conditions.MyModel.Name.Set().Eq(cql.String("a_string_2")),
    )

As you can see, the syntax for the Set method is similar to the queries system with 
the difference that the Set method must be used instead of Is.

For attributes that allow null (nullable values, pointers, nullable relations) the .Set().Null() method will also be available.

Joins
------------------------

It is also possible to perform joins in the first part of the update (Update method):

.. code-block:: go

    type MyOtherModel struct {
        model.UUIDModel

        Name string
    }

    type MyModel struct {
        model.UUIDModel

        Name string

        Related   *MyOtherModel
        RelatedID *model.UUID
    }

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.Related(
            conditions.MyOtherModel.Name.Is().Eq(cql.String("a_string")),
        ),
    ).Set(
        conditions.MyModel.Name.Set().Eq(cql.String("a_string_2")),
    )

Here the only limitation is that in the Set part, only the values of the initial model can be updated 
(not of the joined models). 

This limitation is imposed by the database engines, with the exception of MySQL, 
which allows multiple tables to be updated at the same time. To do this, you use the SetMultiple method:

.. code-block:: go

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.Related(
            conditions.MyOtherModel.Name.Is().Eq(cql.String("a_string")),
        ),
    ).SetMultiple(
        conditions.MyModel.Name.Set().Eq(cql.String("a_string_2")),
        conditions.MyOtherModel.Name.Set().Eq(cql.String("a_string_2")),
    )

Dynamic updates
------------------------

Updates can also be dynamic, meaning that the set can be a value from the same entity or another entity. 
:ref:`Functions <cql/query:functions>`  can also be used on the values.

For example:

.. code-block:: go

    type MyModel struct {
        model.UUIDModel

        Value1 int
        Value2 int
    }

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.Value1.Is().Eq(cql.Int64(2)),
    ).Set(
        conditions.MyModel.Value1.Set().Eq(conditions.MyModel.Value2.Divided(cql.Int64(2))),
    )

Updated at
------------------------

If your model contains a base model with timestamps (model.UUIDModelWithTimestamps or model.UIntModelWithTimestamps), 
cql will automatically add ``updated_at = now()`` to entities that are updated.
