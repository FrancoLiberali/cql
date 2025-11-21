==============================
Delete
==============================

While delete operations can still be performed using gorm's Delete method 
(see `gorm documentation <https://gorm.io/docs/delete.html>`_), 
this is useful only if the model(s) to be delete have already been loaded from the database.

On the contrary, cql's Delete method allows the deletion of all the models that meet 
the conditions entered without the need to load the information (via the direct execution of a DELETE statement).

Delete methods
------------------------

Delete operations are divided into two parts: the Delete method and the Exec method. 
In the first one, we must define the conditions that will determine which models will be deleted. 
Here, the whole system of compilable queries is valid (for details visit :doc:`/cql/query`). 

The object obtained using `cql.Delete` has different methods that 
will allow you to modify the query:

Modifier methods
^^^^^^^^^^^^^^^^^^^^^^^^^^

Modifier methods are those that modify the query in a certain way, affecting the models delete:

- Limit: specifies the number of models to be deleted. (only supported by MySQL)
- Ascending: specifies an ascending order when deleted models. (only supported by MySQL)
- Descending: specifies a descending order when deleted models. (only supported by MySQL)
- Returning: specifies that the models must be fetched from the database after being deleted (the old data is returned) (not supported by MySQL). Preload is not supported. 

Finishing methods
^^^^^^^^^^^^^^^^^^^^^^^

Finishing methods are those that cause the query to be executed:

- Exec: executes the delete

Example
------------------------

.. code-block:: go

    type MyModel struct {
        model.UUIDModel

        Name string
    }

    deletedCount, err := cql.Delete[MyModel](
        context.Background(),
        db,
        conditions.MyModel.Name.Is().Eq(cql.String("a_string")),
    ).Exec()

Joins
------------------------

It is also possible to perform joins in the first part of the delete (except for MySQL):

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

    deletedCount, err := cql.Delete[MyModel](
        context.Background(),
        db,
        conditions.MyModel.Related(
            conditions.MyOtherModel.Name.Is().Eq(cql.String("a_string")),
        ),
    ).Exec()

Here the only limitation is that only the the initial models will be deleted (not of the joined models).

Soft delete
------------------------

Soft delete is also supported by CQL. 
For this, your model must contain a base model that has timestamps: model.UUIDModelWithTimestamps or model.UIntModelWithTimestamps.

For example:

.. code-block:: go

    type MyModel struct {
        model.UUIDModelWithTimestamps

        Name string
    }

Once this is done, cql will automatically take care of:

- Replace DELETE statements with UPDATEs to the deleted_at of the entity.
- Add the condition ``deleted_at is not null`` to your queries, to avoid receiving entities that have been deleted (unless a condition on deleted_at is part of the query you are performing).


Type safety
------------------------

Delete uses the same system of compilable conditions as cql.Query, 
so it shares its features and limitations in terms of type safety at compile time. 

.. TODO actualizar si se mueve
For more details, see :doc:`/cql/type_safety`.

As an added bonus, in cql.Delete, the Returning method is also safe at compile time, 
allowing you to only obtain results in a list of the correct type:

.. code-block:: go
    :caption: Correct
    :linenos:

    myModelsDeleted := []MyModel{}

    deletedCount, err := cql.Delete[MyModel](
        context.Background(),
        db,
        conditions.MyModel.Name.Is().Eq(cql.String("a_string")),
    ).Returning(&deletedModels).Exec()

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 1,7
    :linenos:

    myModelsDeleted := []MyOtherModel{}

    deletedCount, err := cql.Delete[MyModel](
        context.Background(),
        db,
        conditions.MyModel.Name.Is().Eq(cql.String("a_string")),
    ).Returning(&myModelsDeleted).Exec()

In this case, the compilation error will be:

.. code-block:: none

    cannot use &myModelsDeleted (value of type *[]MyOtherModel) as *[]MyModel value in argument to 
    cql.Delete[MyModel](context.Background(), db, conditions.MyModel.Name.Is().Eq(cql.String("a_string"))).Returning