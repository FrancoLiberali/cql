==============================
Insert
==============================

While insert operations can still be performed using gorm's Create or Save methods 
(see `gorm documentation <https://gorm.io/docs/create.htm>`_), 
cql offers cql.Insert, which provides an interface similar to the other methods 
and also adds some advantages in error handling.

Insert methods
------------------------

Insert operations are divided into two parts: the Insert function and the Exec method. 
In the first one, we must express the models to be inserted into the database.

The object obtained using `cql.Insert` has different methods that 
will allow you to modify the query:

Modifier methods
^^^^^^^^^^^^^^^^^^^^^^^^^^

Modifier methods allow the query to be modified for error handling.

They are divided into two groups: those that allow you to select which case you want to control 
and those that allow you to determine the action to be taken in that case.

The first ones are:

- OnConflict: allows to set the action to be taken when any conflict happens.
- OnConflictOn(fields ...FieldOfModel[T]): allows to set the action to be taken when a conflict with the fields specified by parameter happens.
- OnConstraint(constraintName string): allows to set the action to be taken when a conflict with the constraint specified by parameter happens.

Then the possible actions to be taken are:

- DoNothing: not take any action, simply preventing an error from being responded.
- Update(fields ...FieldOfModel[T]): allows you to choose which fields to update with the values of the models that already exist.
- UpdateAll: will update all model attributes with the values of the models that already exist.
- Set(sets ...*Set[T]): allows you to configure specific updates to be performed. Its syntax is the same as the :ref:`sets in the cql.Update method <cql/update:Finishing methods>`.

.. warning::

    In postgres OnConflict can be used only with DoNothing. For UpdateAll, Update and Set, OnConflictOn must be used.

In addition, the action Set allows a third method:

- Where: Where allows to choose to execute the sets only in some of the conflict cases. Only available for postgres and sqlite.

Finishing methods
^^^^^^^^^^^^^^^^^^^^^^^

Finishing methods are those that cause the query to be executed:

- Exec: executes the insert, returning the amount of rows inserted.
- ExecInBatches(batchSize uint): execute the insert statement in batches of batchSize, returning the amount of rows inserted.

Examples
------------------------

.. code-block:: go
    :caption: Model definition

    type MyModel struct {
        model.UUIDModel

        Name   string
        Status int
    }

.. code-block:: go
    :caption: Single insert

    myModel := &MyModel{
        Name: "myModelName",
    }

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel,
    ).Exec()

.. code-block:: go
    :caption: Multiple insert

    myModel1 := &MyModel{
        Name: "myModelName1",
    }

    myModel2 := &MyModel{
        Name: "myModelName2",
    }

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel1,
        myModel2,
    ).Exec()

.. code-block:: go
    :caption: Multiple insert in batches

    myModel1 := &MyModel{
        Name: "myModelName1",
    }

    myModel2 := &MyModel{
        Name: "myModelName2",
    }

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel1,
        myModel2,
    ).ExecInBatches(1)

.. code-block:: go
    :caption: On conflict do nothing

    myModel := &MyModel{
        Name: "myModelName",
    }

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel,
    ).OnConflict().DoNothing().Exec()

.. code-block:: go
    :caption: On conflict on name update all

    myModel := &MyModel{
        Name: "myModelName",
        Status: 1,
    }

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel,
    ).OnConflictOn(conditions.MyModel.Name).UpdateAll().Exec()

.. code-block:: go
    :caption: On constraint on name update the status

    myModel := &MyModel{
        Name: "myModelName",
        Status: 1,
    }

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel,
    ).OnConstraint("mymodel_unique_name").Update(conditions.MyModel.Status).Exec()

.. code-block:: go
    :caption: On conflict set a value

    myModel := &MyModel{
        Name: "myModelName",
        Status: 1,
    }

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel,
    ).OnConflict().Set(
        conditions.MyModel.Status.Set().Eq(cql.Int(3)),
    ).Exec()

.. code-block:: go
    :caption: On conflict set a value in some cases

    myModel1 := &MyModel{
        Name: "myModelName1",
        Status: 1,
    }

    myModel2 := &MyModel{
        Name: "myModelName2",
        Status: 1,
    }

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel,
    ).OnConflict().Set(
        conditions.MyModel.Status.Set().Eq(cql.Int(3)),
    ).Where(
        conditions.MyModel.Name.Is().Eq(cql.String("myModelName1")),
    ).Exec()

Type safety
------------------------

OnConflictOn and Update
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

In terms of type safety, methods that receive fields (OnConflictOn and Update) only allow fields from the initial model.

OnConflictOn:

.. code-block:: go
    :caption: Correct
    :linenos:

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel,
    ).OnConflictOn(conditions.MyModel.Name).UpdateAll().Exec()

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 5
    :linenos:

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel,
    ).OnConflictOn(conditions.MyOtherModel.Name).UpdateAll().Exec()

In this case, the compilation error will be:

.. code-block:: none

    cannot use conditions.MyOtherModel.Name (variable of struct type condition.StringField[models.MyOtherModel])
    as condition.FieldOfModel[models.MyModel] value in argument to cql.Insert(context.Background(), db, myModel).OnConflictOn:
    condition.StringField[models.MyOtherModel] does not implement condition.FieldOfModel[models.MyModel] (wrong type for method getModel)

Update:

.. code-block:: go
    :caption: Correct
    :linenos:

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel,
    ).OnConflict().Update(conditions.MyModel.Name).Exec()

.. code-block:: go
    :class: with-errors
    :caption: Incorrect
    :emphasize-lines: 5
    :linenos:

    insertedCount, err := cql.Insert(
        context.Background(),
        db,
        myModel,
    ).OnConflict().Update(conditions.MyOtherModel.Name).Exec()

In this case, the compilation error will be:

.. code-block:: none

    cannot use conditions.MyOtherModel.Name (variable of struct type condition.StringField[models.MyOtherModel])
    as condition.FieldOfModel[models.MyModel] value in argument to cql.Insert(context.Background(), db, myModel).OnConflict().Update:
    condition.StringField[models.MyOtherModel] does not implement condition.FieldOfModel[models.MyModel] (wrong type for method getModel)

Set
^^^^^^^^^^^^^^^^^^^^^^^

In the case of Set, since it is the same system as cql.Update, it shares its features and limitations in terms of type safety at compile time.
For details, see :ref:`cql/update:Set`.

Where
^^^^^^^^^^^^^^^^^^^^^^^

In the case of Where, since it is the same system as cql.Query, it shares its features and limitations in terms of type safety at compile time.

For more details, see :doc:`/cql/query_type_safety`.

Type safety limitations and cqllint
------------------------------------------------

The OnConstraint method is not safe at compile time, since CQL has no way of knowing which constraints are defined in the database. 
If you try to use one that is not defined, the error returned will be the error returned by the database:

.. code-block:: none

    ERROR: constraint "do_not_exists" for table "cities" does not exist (SQLSTATE 42704)
