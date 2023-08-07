==============================
Query
==============================

Query methods
------------------------

In CRUDRepository you will find different methods that will 
allow you to perform queries on the model to which that repository belongs:

- GetByID: will allow you to obtain a model by its id.
- QueryOne: will allow you to obtain the model that meets the conditions received by parameter.
- Query: will allow you to obtain the models that meet the conditions received by parameter.

Compilable query system
------------------------

The set of conditions that are received by the read operations of the CRUDService 
and CRUDRepository form the badaas-orm compilable query system. 
It is so named because the conditions will verify at compile time that the query to be executed is correct.

These conditions are objects of type Condition that contain the 
necessary information to perform the queries in a safe way. 
They are generated from the definition of your models using badaas-cli.

Conditions generation
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The generation of conditions is done with badaas-cli. For this, we need to install badaas-cli:

.. code-block:: bash

    go install github.com/ditrit/badaas-cli

Then, inside our project we will have to create a package called conditions 
(or another name if you wish) and inside it a file with the following content:

.. code-block:: go

    package conditions

    //go:generate badaas-cli gen conditions ../models_path_1 ../models_path_2

where ../models_path_1 ../models_path_2 are the relative paths between the package conditions 
and the packages containing the definition of your models (can be only one).

Now, from the root of your project you can execute:

.. code-block:: bash

  go generate ./...

and the conditions for each of your models will be created in the conditions package.

Use of the conditions
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

After generating the conditions you will have the following conditions:

- One condition for each attribute of each of your models. 
  The name of these conditions will be <Model><Attribute> where 
  <Model> is the model type and <Attribute> is the attribute name. 
  These conditions are of type WhereCondition.
- One condition for each relationship with another model that each of your models has. 
  The name of these conditions will be <Model><Relation> where 
  <Model> is the model type and <Relation> is the name of the attribute that creates the relation. 
  These conditions are of type JoinCondition because using them will 
  mean performing a join within the executed query.

Then, combining these conditions, the Connection Conditions (orm.And, orm.Or, orm.Not) 
and the Operators (orm.Eq, orm.Lt, etc.) you will be able to make all 
the queries you need in a safe way.

Examples
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**Filter by an attribute**

In this example we query all YourModel that has "a_string" in the Attribute attribute.

.. code-block:: go

    type YourModel struct {
        model.UUIDModel

        Attribute string
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelAttribute(orm.Eq("a_string")),
    )

**Filter by an attribute of a related model**

In this example we query all YourModels whose related Related has "a_string" in its Attribute attribute.

.. code-block:: go

    type Related struct {
        model.UUIDModel

        Attribute string
    }

    type YourModel struct {
        model.UUIDModel

        Related   Related
        RelatedID model.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelRelated(
            conditions.RelatedAttribute(orm.Eq("a_string")),
        ),
    )

**Multiple conditions**

In this example we query all YourModels that has a 4 in the IntAttribute attribute and 
whose related Related has "a_string" in its Attribute attribute.

.. code-block:: go

    type Related struct {
        model.UUIDModel

        Attribute string
    }

    type YourModel struct {
        model.UUIDModel

        IntAttribute int

        Related   Related
        RelatedID model.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelIntAttribute(orm.Eq(4)),
        conditions.YourModelRelated(
            conditions.RelatedAttribute(orm.Eq("a_string")),
        ),
    )

Operators
------------------------

Below you will find the complete list of available operators:

- orm.Eq(value): EqualTo
- orm.EqOrIsNull(value): if value is not NULL returns a Eq operator but if value is NULL returns a IsNull operator
- orm.NotEq(value): NotEqualTo
- orm.NotEqOrIsNotNull(value): if value is not NULL returns a NotEq operator but if value is NULL returns a IsNotNull operator
- orm.Lt(value): LessThan
- orm.LtOrEq(value): LessThanOrEqualTo
- orm.Gt(value): GreaterThan
- orm.GtOrEq(value): GreaterThanOrEqualTo
- orm.IsNull()
- orm.IsNotNull()
- orm.Between(v1, v2): Equivalent to v1 < attribute < v2
- orm.NotBetween(v1, v2): Equivalent to NOT (v1 < attribute < v2)
- orm.IsTrue()
- orm.IsNotTrue()
- orm.IsFalse()
- orm.IsNotFalse()
- orm.IsUnknown()
- orm.IsNotUnknown()
- orm.IsDistinct(value)
- orm.IsNotDistinct(value)
- orm.Like(pattern)
- orm.Like(pattern).Escape(escape)
- orm.ArrayIn(values)
- orm.ArrayNotIn(values)