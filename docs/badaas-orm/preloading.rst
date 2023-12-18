==============================
Preloading
==============================

PreloadConditions
---------------------------

During the :ref:`conditions generation <badaas-orm/query:conditions generation>` the following 
PreloadConditions are also generated which are useful for preloading:

- One PreloadCondition for each of your models, that will allow to preload this model when doing a query.
  The name of these conditions will be <Model>PreloadAttributes where 
  <Model> is the model type.
- One PreloadCondition for each of the relations of your model, 
  to preload that the related object when doing a query. 
  This is really just a facility that translates to using the JoinCondition of 
  that relation and then the PreloadAttributes of the related model.
  The name of these conditions will be <Model>Preload<Relation> where 
  <Model> is the model type and <Relation> is the name of the attribute that creates the relation.
- One PreloadCondition to preload all the related models of your model.
  The name of these conditions will be <Model>PreloadRelations where 
  <Model> is the model type.

Examples
----------------------------------

**Preload a related model**

In this example we query all YourModels and preload whose related Related.

.. code-block:: go

    type Related struct {
        model.UUIDModel
    }

    type YourModel struct {
        model.UUIDModel

        Related   Related
        RelatedID model.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelRelated(
            conditions.RelatedPreloadAttributes,
        ),
    )

Or using the PreloadCondition to avoid the JoinCondition 
(only useful when you don't want to add other conditions to that Join):

.. code-block:: go

    type Related struct {
        model.UUIDModel
    }

    type YourModel struct {
        model.UUIDModel

        Related   Related
        RelatedID model.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelPreloadRelated,
    )

**Preload a list of models**

.. code-block:: go

    type Related struct {
        model.UUIDModel

        YourModel *YourModel
        YourModelID *model.UUID
    }

    type YourModel struct {
        model.UUIDModel

        Related *[]Related
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelPreloadRelated,
    )

**Nested preloads**

.. code-block:: go

    type Parent struct {
        model.UUIDModel
    }

    type Related struct {
        model.UUIDModel

        Parent   Parent
        ParentID model.UUID
    }

    type YourModel struct {
        model.UUIDModel

        Related   Related
        RelatedID model.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelRelated(
            conditions.RelatedPreloadParent,
        ),
    )

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

    type Related struct {
        model.UUIDModel
    }

    type YourModel struct {
        model.UUIDModel

        Related   Related
        RelatedID model.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelPreloadRelated,
    )

    if err == nil && len(yourModels) > 1 {
        firstRelated, err := yourModels[0].GetRelated()
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