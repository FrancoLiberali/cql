==============================
Advanced query
==============================

Collections
-------------------------------

cql also allows you to set conditions on a collection of models (one to many or many to many relationships):

.. code-block:: go
    :caption: Example model

    type Seller struct {
        model.UUIDModel

        Name string

        Company   *Company
        CompanyID *model.UUID // Company HasMany Seller (Company 0..1 -> 0..* Seller)
    }

    type Company struct {
        model.UUIDModel

        Sellers *[]Seller // Company HasMany Seller (Company 0..1 -> 0..* Seller)
    }

.. code-block:: go
    :caption: Query

    companies, err := cql.Query[Company](
        context.Background(),
        db,
        conditions.Company.Sellers.Any(
            conditions.Seller.Name.Is().Eq(cql.String("franco")),
        ),
    ).Find()

The methods for collections are:

- None: generates a condition that is true if no model in the collection fulfills the conditions.
- Any: generates a condition that is true if at least one model in the collection fulfills the conditions.
- All: generates a condition that is true if all models in the collection fulfill the conditions (or is empty).

Dynamic operators
--------------------------------

In :doc:`/cql/query` we have seen how to use the operators 
to make comparisons between the attributes of a model and static values such as a string, 
a number, etc. But if we want to make comparisons between two or more attributes of 
the same type we need to use the dynamic operators. 
These receive a Field, an object that identifies the attribute with which the operation is to be performed.

These identifiers are also generated during the generation of conditions 
as attributes of the condition model 
(if models.MyModel.Name exists, then conditions.MyModel.Name is generated).

For example we query all MyModels that has the same value in its Name attribute that 
its related MyOtherModel's Name attribute.

.. code-block:: go
    :caption: Example model

    type MyOtherModel struct {
        model.UUIDModel

        Name string
    }

    type MyModel struct {
        model.UUIDModel

        Name string

        Related   MyOtherModel
        RelatedID model.UUID
    }

.. code-block:: go
    :caption: Query

    myModels, err := cql.Query[MyModel](
        context.Background(),
        db,
        conditions.MyModel.Related(
            conditions.MyOtherModel.Name.Is().Eq(conditions.MyModel.Name),
        ),
    ).Find()

**Attention**, when using dynamic operators the verification that the Field 
is concerned by the query is performed at run time, returning an error otherwise. 
For example:

.. code-block:: go
    :caption: Example model

     type MyOtherModel struct {
        model.UUIDModel

        Name string
    }

    type MyModel struct {
        model.UUIDModel

        Name string

        Related   MyOtherModel
        RelatedID model.UUID
    }

.. code-block:: go
    :caption: Query

    myModels, err := cql.Query[MyModel](
        context.Background(),
        db,
        conditions.MyModel.Name.Is().Eq(conditions.MyOtherModel.Name),
    ).Find()

will respond cql.ErrFieldModelNotConcerned in err.

Dynamic functions
--------------------------------

Functions can also be applied between different attributes:

.. code-block:: go
    conditions.MyModel.Attribute1.Divided(conditions.MyModel.Attribute2).Is().Lt(cql.Int64(10))

within dynamic operators:

.. code-block:: go
    conditions.MyModel.Attribute1.Is().Lt(conditions.MyModel.Attribute2.Minus(cql.Int64(10)))

or both of them:

.. code-block:: go
    conditions.MyModel.Attribute1.Divided(conditions.MyModel.Attribute2).Is().Lt(
        conditions.MyModel.Attribute2.Minus(conditions.MyModel.Attribute1),
    )

In all cases, the attributes to be used in the functions may belong to the same or different entities.

For example, if we seek to obtain the cities whose population represents at least half of the population of their country:

.. code-block:: go
    :caption: Example model

    type Country struct {
        model.UUIDModel

        Population int
    }

    type City struct {
        model.UUIDModel

        Population int

        Country   Country
        CountryID model.UUID
    }

.. code-block:: go
    :caption: Query
    :linenos:
    :emphasize-lines: 5

    cities, err := cql.Query[City](
        context.Background(),
        db,
        conditions.City.Country(
            conditions.Country.Population.Is().Lt(
                conditions.City.Population.Times(2),
            ),
        ),
    ).Find()


Appearance
-------------------------

In case the attribute to be used is present more 
than once in the query, it will be necessary to select select its appearance number, 
to avoid getting the error cql.ErrAppearanceMustBeSelected. 
To do this, you must use the Appearance method of the field, as in the following example:

.. code-block:: go
    :caption: Example model

    type ParentParent struct {
        model.UUIDModel
    }

    type Parent1 struct {
        model.UUIDModel

        ParentParent   ParentParent
        ParentParentID model.UUID
    }

    type Parent2 struct {
        model.UUIDModel

        ParentParent   ParentParent
        ParentParentID model.UUID
    }

    type Child struct {
        model.UUIDModel

        Parent1   Parent1
        Parent1ID model.UUID

        Parent2   Parent2
        Parent2ID model.UUID
    }

.. code-block:: go
    :caption: Query
    :linenos:
    :emphasize-lines: 11

    models, err := cql.Query[Child](
        context.Background(),
        db,
        conditions.Child.Parent1(
            conditions.Parent1.ParentParent(),
        ),
        conditions.Child.Parent2(
            conditions.Parent2.ParentParent(),
        ),
        conditions.Child.Name.Is().Eq(
            conditions.ParentParent.Name.Appearance(0), // choose the first (0) appearance (made by conditions.Child.Parent1())
        ),
    ).Find()

Unsafe operators
--------------------------------

In case you want to avoid the type validations performed by the operators, 
unsafe operators should be used. 
Although their use is not recommended, this can be useful when the database 
used allows operations between different types or when attributes of different 
types map at the same time in the database (see <https://gorm.io/docs/data_types.html>).

If it is neither of these two cases, the use of an unsafe operator will result in 
an error in the execution of the query that depends on the database used.

All operators supported by cql that receive any value are available 
in their unsafe version after using the IsUnsafe() method of the Field object.


Unsafe conditions (raw SQL)
--------------------------------

In case you need to use operators that are not supported by cql
(please create an issue in our repository if you think we have forgotten any), 
you can always run raw SQL with unsafe.NewCondition, as in the following example:

.. code-block:: go

    myModels, err := cql.Query[MyModel](
        context.Background()
        db,
        unsafe.NewCondition[MyModel]("%s.name = NULL"),
    ).Find()

As you can see in the example, "%s" can be used in the raw SQL to be replaced 
by the table name of the model to which the condition belongs.

Of course, its use is not recommended because it can generate errors in the execution 
of the query that will depend on the database used.