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
:ref:`Functions <cql/query:functions>` can also be used on the values.

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

Type safety
------------------------

Update uses the same system of compilable conditions as cql.Query, 
so it shares its features and limitations in terms of type safety at compile time. 

.. TODO actualizar si se mueve
For more details, see :doc:`/cql/type_safety`.

Set
^^^^^^^^^^^^^^^^^^^^^^^

In addition, Update also provides the same type safety in Set methods, 
ensuring that the value to be set is of the same type as the attribute to be modified:

.. code-block:: go
    :caption: Model
    :linenos:

    type MyModel struct {
        model.UUIDModel

        ValueInt    int
        ValueString string
    }

.. code-block:: go
    :caption: Correct
    :linenos:

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.ValueInt.Is().Eq(cql.Int64(2)),
    ).Set(
        conditions.MyModel.ValueInt.Set().Eq(conditions.MyModel.ValueInt.Divided(cql.Int64(2))),
    )

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 6
    :linenos:

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.ValueInt.Is().Eq(cql.Int64(2)),
    ).Set(
        conditions.MyModel.ValueInt.Set().Eq(conditions.MyModel.ValueString),
    )

In this case, the compilation error will be:

.. code-block:: none

    cannot use conditions.MyModel.ValueString (variable of struct type condition.StringField[MyModel]) as 
    condition.ValueOfType[float64] value in argument to conditions.MyModel.ValueInt.Set().Eq: 
    condition.StringField[MyModel] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)

Returning
^^^^^^^^^^^^^^^^^^^^^^^

In cql.Update, the Returning method is also safe at compile time, 
allowing you to only obtain results in a list of the correct type:

.. code-block:: go
    :caption: Correct
    :linenos:

    myModelsUpdated := []MyModel{}

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.ValueInt.Is().Eq(cql.Int64(2)),
    ).Returning(&myModelsUpdated).Set(
        conditions.MyModel.ValueInt.Set().Eq(cql.Int64(3)),
    )

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 1,7
    :linenos:

    myModelsUpdated := []MyOtherModel{}

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.ValueInt.Is().Eq(cql.Int64(2)),
    ).Returning(&myModelsUpdated).Set(
        conditions.MyModel.ValueInt.Set().Eq(cql.Int64(3)),
    )

In this case, the compilation error will be:

.. code-block:: none

    cannot use &myModelsUpdated (value of type *[]MyOtherModel) as *[]MyModel value in argument to 
    cql.Update[MyModel](context.Background(), db, conditions.MyModel.ValueInt.Is().Eq(cql.Int64(2))).Returning

Null update
^^^^^^^^^^^^^^^^^^^^^^^

For fields that are nullable, such as pointers or null.* types, cql.Update will 
allow you to safely set their value to null at compile time, i.e., giving a compile-time error 
if you try to update a non-nullable attribute to null:

.. code-block:: go
    :caption: Model
    :linenos:

    type MyModel struct {
        model.UUIDModel

        ValueInt        int
        ValueIntPointer *int
    }

.. code-block:: go
    :caption: Correct
    :linenos:

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.ValueInt.Is().Eq(cql.Int64(2)),
    ).Set(
        conditions.MyModel.ValueIntPointer.Set().Null(),
    )

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 6
    :linenos:

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.ValueInt.Is().Eq(cql.Int64(2)),
    ).Set(
        conditions.MyModel.ValueInt.Set().Null(),
    )

In this case, the compilation error will be:

.. code-block:: none

    conditions.MyModel.ValueInt.Set().Null undefined 
    (type condition.FieldSet[MyModel, int] has no field or method Null)

Limitations
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Once again, similar to cql.Query, Set is not safe at compile time to determine whether 
the values used in Eq or used in functions in Eq are joined in the query, as in the following examples:

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 6
    :linenos:

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.ValueInt.Is().Eq(cql.Int64(2)),
    ).Set(
        conditions.MyModel.ValueInt.Set().Eq(conditions.MyOtherModel.Value1),
    )

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 6
    :linenos:

    updatedCount, err := cql.Update[MyModel](
        context.Background(),
        db,
        conditions.MyModel.ValueInt.Is().Eq(cql.Int64(2)),
    ).Set(
        conditions.MyModel.ValueInt.Set().Eq(conditions.MyModel.ValueInt.Plus(conditions.MyOtherModel.Value1)),
    )

Which would generate the following error at runtime:

.. code-block:: none

    field's model is not concerned by the query (not joined); not concerned model: models.MyOtherModel

.. TODO link a la seccion correcta
These errors can be determined before runtime using :doc:`/cql/cqllint`.