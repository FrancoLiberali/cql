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
        context.Background(),
        db,
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
        context.Background(),
        db,
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
        context.Background(),
        db,
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

Model collections can also be preloaded (relations has many or many to many): 

.. code-block:: go
    :caption: Example model

    type Seller struct {
        model.UUIDModel

        Company   *Company
        CompanyID *model.UUID // Company HasMany Seller (Company 0..1 -> 0..* Seller)
    }

    type Company struct {
        model.UUIDModel

        Sellers *[]Seller // Company HasMany Seller (Company 0..1 -> 0..* Seller)
    }

.. code-block:: go
    :caption: Query

    company, err := cql.Query[Company](
        context.Background(),
        db,
        conditions.Company.Sellers.Preload(),
    ).FindOne()

    if err == nil {
        sellers, err := company.GetSellers()
        if err == nil {
            // you can safely apply your business logic
        } else {
            // err is cql.ErrRelationNotLoaded
        }
    }

Nested preloads can also be applied to preload model relationships within the collection:

.. code-block:: go
    :caption: Example model

    type Office struct {
        model.UUIDModel

        Seller   *Seller
        SellerID *model.UUID `gorm:"not null"` // Seller HasOne Office (Seller 1 -> 1 Office)
    }

    type Seller struct {
        model.UUIDModel

        Office   *Office // Seller HasOne Office (Seller 1 -> 1 Office)

        Company   *Company
        CompanyID *model.UUID // Company HasMany Seller (Company 0..1 -> 0..* Seller)
    }

    type Company struct {
        model.UUIDModel

        Sellers *[]Seller // Company HasMany Seller (Company 0..1 -> 0..* Seller)
    }

.. code-block:: go
    :caption: Query

    company, err := cql.Query[Company](
        context.Background(),
        db,
        conditions.Company.Sellers.Preload(
            conditions.Seller.Office().Preload()
        ),
    ).FindOne()

    if err == nil {
        sellers, err := company.GetSellers()
        if err == nil {
            for _, seller := range sellers {
                office, err := seller.GetOffice()
                if err == nil {
                    // you can safely apply your business logic
                } else {
                    // err is cql.ErrRelationNotLoaded
                }
            }
        } else {
            // err is cql.ErrRelationNotLoaded
        }
    }