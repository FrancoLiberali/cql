==============================
Preloading
==============================

PreloadConditions
---------------------------

During the :ref:`conditions generation <badaas-orm/query:conditions generation>` the following 
methods will also be created for the condition models:

- Preload() will allow to preload this model when doing a query.
- Preload<Relation>() for each of the relations of your model, 
  where <Relation> is the name of the attribute that creates the relation, 
  to preload that the related object when doing a query. 
  This is really just a facility that translates to using the JoinCondition of 
  that relation and then the Preload method of the related model.
- PreloadRelation() to preload all the related models of your model 
  (only generated if the model has at least one relation).

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

    myModels, err := orm.NewQuery[MyModel](
        gormDB,
        conditions.MyModel.Related(
            conditions.Related.Preload(),
        ),
    ).Find()

Or using the PreloadRelation method to avoid the JoinCondition 
(only useful when you don't want to add other conditions to that Join):

.. code-block:: go

    myModels, err := orm.NewQuery[MyModel](
        gormDB,
        conditions.MyModel.PreloadRelated(),
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

    myModels, err := orm.NewQuery[MyModel](
        gormDB,
        conditions.MyModel.Related(
            conditions.MyOtherModel.PreloadParent(),
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

For this reason, badaas-orm provides the Relation getters. 
These are methods that will be added to your models to safely navigate a relation, 
responding `errors.ErrRelationNotLoaded` in case you try to navigate a relation 
that was not loaded from the database. 
They are created in a file called badaas-orm.go in your model package when 
:ref:`generating conditions <badaas-orm/concepts:conditions generation>`.

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

    myModel, err := orm.NewQuery[MyModel](
        conditions.MyModel.PreloadRelated(),
    ).FindOne()

    if err == nil {
        firstRelated, err := myModel.GetRelated()
        if err == nil {
            // you can safely apply your business logic
        } else {
            // err is errors.ErrRelationNotLoaded
        }
    }

Unfortunately, these relation getters cannot be created in all cases but only in those in which:

- The relation is made with an object directly instead of a pointer 
  (which is not recommended as described :ref:`here <badaas-orm/declaring_models:references>`).
- The relation is made with pointers and the foreign key (typically the ID) is in the same model.
- The relation is made with a pointer to a list.