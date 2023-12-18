==============================
Advanced query
==============================

Dynamic operators
--------------------------------

In :doc:`/badaas-orm/query` we have seen how to use the operators 
to make comparisons between the attributes of a model and static values such as a string, 
a number, etc. But if we want to make comparisons between two or more attributes of 
the same type we need to use the dynamic operators. 
These, instead of a dynamic value, receive a FieldIdentifier, that is, 
an object that identifies the attribute with which the operation is to be performed.

These identifiers are also generated during the generation of conditions 
as attributes of the condition model 
(if models.MyModel.Name exists, then conditions.MyModel.Name is generated).

For example we query all MyModels that has the same value in its Name attribute that 
its related MyOtherModel's Name attribute.

.. code-block:: go

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

    myModels, err := orm.NewQuery[MyModel](
        gormDB,
        conditions.MyModel.Related(
            conditions.MyOtherModel.NameIs().Dynamic().Eq(conditions.MyModel.Name),
        ),
    ).Find()

**Attention**, when using dynamic operators the verification that the FieldIdentifier 
is concerned by the query is performed at run time, returning an error otherwise. 
For example:

.. code-block:: go

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

    myModels, err := orm.NewQuery[MyModel](
        gormDB,
        conditions.MyModel.NameIs().Dynamic().Eq(conditions.MyOtherModel.Name),
    ).Find()

will respond orm.ErrFieldModelNotConcerned in err.

All operators supported by badaas-orm that receive any value are available in their dynamic version 
after using the Dynamic() method of the FieldIs object.

Select join
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

In case the attribute to be used by the dynamic operator is present more 
than once in the query, it will be necessary to select the join to be used, 
to avoid getting the error orm.ErrJoinMustBeSelected. 
To do this, you must use the SelectJoin method, as in the following example:

.. code-block:: go

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

    models, err := orm.NewQuery[Child](
        gormDB,
        conditions.Child.Parent1(
            conditions.Parent1.ParentParent(),
        ),
        conditions.Child.Parent2(
            conditions.Parent2.ParentParent(),
        ),
        conditions.Child.NameIs().Dynamic().Eq(conditions.ParentParent.Name).SelectJoin(
            0, // for the value 0 (conditions.ParentParent.Name),
            0, // choose the first (0) join (made by conditions.Child.Parent1())
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

All operators supported by badaas-orm that receive any value are available 
in their unsafe version after using the Unsafe() method of the FieldIs object.


Unsafe conditions (raw SQL)
--------------------------------

In case you need to use operators that are not supported by badaas-orm 
(please create an issue in our repository if you think we have forgotten any), 
you can always run raw SQL with unsafe.NewCondition, as in the following example:

.. code-block:: go

    myModels, err := orm.NewQuery[MyModel](
        gormDB,
        unsafe.NewCondition[MyModel]("%s.name = NULL"),
    ).Find()

As you can see in the example, "%s" can be used in the raw SQL to be replaced 
by the table name of the model to which the condition belongs.

Of course, its use is not recommended because it can generate errors in the execution 
of the query that will depend on the database used.