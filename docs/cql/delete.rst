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
Here, the whole system of compilable queries is valid (for details visit :ref:`cql/query`). 

The object obtained using `cql.Delete` has different methods that 
will allow you to modify the query:

Modifier methods
^^^^^^^^^^^^^^^^^^^^^^^^^^

Modifier methods are those that modify the query in a certain way, affecting the models delete:
- Limit: specifies the number of models to be deleted. (only supported by MySQL)
- Ascending: specifies an ascending order when deleted models. (only supported by MySQL)
- Descending: specifies a descending order when deleted models. (only supported by MySQL)
- Returning: specifies that the models models must be fetched from the database after being deleted 
(the old data is returned) (not supported by MySQL). Preload of related data is also possible (not supported by SQLite). 

Finishing methods
^^^^^^^^^^^^^^^^^^^^^^^

Finishing methods are those that cause the query to be executed:

- Exec: executes the delete

Example
^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: go

    type MyModel struct {
        model.UUIDModel

        Name string
    }

    deletedCount, err := cql.Delete[MyModel](
        gormDB,
        conditions.MyModel.Name.Is().Eq("a_string"),
    ).Exec()

Joins
^^^^^^^^^^^^^^^^^^^^^^^

It is also possible to perform joins in the first part of the delete (Delete method):

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
        gormDB,
        conditions.MyModel.Related(
            conditions.MyOtherModel.Name.Is().Eq("a_string"),
        ),
    ).Exec()

Here the only limitation is that only the the initial models will be deleted (not of the joined models). 