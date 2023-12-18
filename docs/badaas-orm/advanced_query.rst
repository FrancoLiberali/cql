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

These identifiers are also generated during the generation of conditions and 
their name of these FieldIdentifiers will be <Model><Attribute>Field where 
<Model> is the model type and <Attribute> is the attribute name.

For example we query all YourModels that has the same value in its String attribute that 
its related Related's String attribute.

.. code-block:: go

    type Related struct {
        model.UUIDModel

        String string
    }

    type YourModel struct {
        model.UUIDModel

        String string

        Related   Related
        RelatedID model.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelRelated(
            conditions.RelatedString(
                dynamic.Eq(conditions.YourModelStringField),
            ),
        ),
    )

**Attention**, when using dynamic operators the verification that the FieldIdentifier 
is concerned by the query is performed at run time, returning an error otherwise. 
For example:

.. code-block:: go

    type Related struct {
        model.UUIDModel

        String string
    }

    type YourModel struct {
        model.UUIDModel

        String string

        Related   Related
        RelatedID model.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelString(
            dynamic.Eq(conditions.RelatedStringField),
        ),
    )

will respond orm.ErrFieldModelNotConcerned in err.

All operators supported by badaas-orm that receive any value are available in their dynamic version at
<https://pkg.go.dev/github.com/ditrit/badaas/orm/dynamic>. 

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

    models, err := ts.crudChildService.Query(
        conditions.ChildParent1(
            conditions.Parent1ParentParent(),
        ),
        conditions.ChildParent2(
            conditions.Parent2ParentParent(),
        ),
        conditions.ChildName(
            // for the value 0 (conditions.ParentParentNameField),
            // choose the first (0) join (made by conditions.ChildParent1())
            dynamic.Eq(conditions.ParentParentNameField).SelectJoin(0, 0),
        ),
    )

Unsafe operators
--------------------------------

In case you want to avoid the type validations performed by the operators, unsafe operators should be used. 
Although their use is not recommended, this can be useful when the database 
used allows operations between different types or when attributes of different 
types map at the same time in the database (see <https://gorm.io/docs/data_types.html>).

If it is neither of these two cases, the use of an unsafe operator will result in 
an error in the execution of the query that depends on the database used.

All operators supported by badaas-orm that receive any value are available in their unsafe version at
<https://pkg.go.dev/github.com/ditrit/badaas/orm/unsafe>. 

Unsafe conditions (raw SQL)
--------------------------------

In case you need to use operators that are not supported by badaas-orm 
(please create an issue in our repository if you think we have forgotten any), 
you can always run raw SQL with unsafe.NewCondition, as in the following example:

.. code-block:: go

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelString(
            unsafe.NewCondition[models.YourModel]("%s.name = NULL"),
        ),
    )

As you can see in the example, "%s" can be used in the raw SQL to be replaced 
by the table name of the model to which the condition belongs.

Of course, its use is not recommended because it can generate errors in the execution 
of the query that will depend on the database used.