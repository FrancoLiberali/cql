==============================
Preloading
==============================

When doing a join, conditions can be applied on joined models but, 
by default, only the information of the main model is returned as a result. 
To also get the joined models, it is necessary to use the Preload() method.

Examples
----------------------------------

**Preload a related model**

In this example we query all MyModels and preload whose related MyOtherModel.

.. code-block:: go

    type MyOtherModel struct {
        model.UUIDModel
    }

    type MyModel struct {
        model.UUIDModel

        Related   MyOtherModel
        RelatedID model.UUID
    }

    myModels, err := cql.Query[MyModel](
        gormDB,
        conditions.MyModel.Related().Preload(),
    ).Find()

**Nested preloads**

.. code-block:: go

    type Parent struct {
        model.UUIDModel
    }

    type MyOtherModel struct {
        model.UUIDModel

        Parent   Parent
        ParentID model.UUID
    }

    type MyModel struct {
        model.UUIDModel

        Related   MyOtherModel
        RelatedID model.UUID
    }

    myModels, err := cql.Query[MyModel](
        gormDB,
        conditions.MyModel.Related(
            conditions.MyOtherModel.Parent().Preload(),
        ),
    ).Find()

As we can see, it is not necessary to add the preload to all joins, 
it is enough to do it in the deepest one, 
to recover, in this example, both Related and Parent.

Relation getters
--------------------------------------

At the moment, with the PreloadConditions, we can choose whether or not to preload a relation. 
The problem is that once we get the result of the query, we cannot determine if a null value 
corresponds to the fact that the relation is really null or that the preload was not performed, 
which means a big risk of making decisions in our business logic on incomplete information.

For this reason, cql provides the Relation getters. 
These are methods that will be added to your models to safely navigate a relation, 
responding `cql.ErrRelationNotLoaded` in case you try to navigate a relation 
that was not loaded from the database. 
They are created in a file called cql.go in your model package when 
:ref:`generating conditions <cql/concepts:conditions generation>`.

Here is an example of its use:

.. code-block:: go

    type MyOtherModel struct {
        model.UUIDModel
    }

    type MyModel struct {
        model.UUIDModel

        Related   MyOtherModel
        RelatedID model.UUID
    }

    myModel, err := cql.Query[MyModel](
        conditions.MyModel.Related().Preload(),
    ).FindOne()

    if err == nil {
        firstRelated, err := myModel.GetRelated()
        if err == nil {
            // you can safely apply your business logic
        } else {
            // err is cql.ErrRelationNotLoaded
        }
    }

Unfortunately, these relation getters cannot be created in all cases but only in those in which:

- The relation is made with an object directly instead of a pointer 
  (which is not recommended as described :ref:`here <cql/declaring_models:references>`).
- The relation is made with pointers and the foreign key (typically the ID) is in the same model.
- The relation is made with a pointer to a list.

Preload collections
---------------------------

During the :ref:`conditions generation <cql/query:conditions generation>` the following 
methods will also be created for the condition models:

- Preload<Collection>() for each of the collection of models of your model, 
  where <Collection> is the name of the collection,  to preload it when doing a query. 