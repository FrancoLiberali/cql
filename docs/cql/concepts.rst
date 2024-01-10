==============================
Concepts
==============================

Model
------------------------------

A model is any object (struct) of go that you want to persist 
in the database and on which you can perform queries. 
For this, the struct must have an embedded cql base model.

For details visit :ref:`cql/declaring_models:model declaration`.

Base model
-----------------------------

It is a struct that when embedded allows your structures to become cql models, 
adding ID, CreatedAt, UpdatedAt and DeletedAt attributes and the possibility to persist, 
create conditions and perform queries on these structures.

For details visit :ref:`cql/declaring_models:base models`.

Model ID
-----------------------------

The id is a unique identifier needed to persist a model in the database. 
It can be a model.UIntID or a model.UUID, depending on the base model used.

For details visit :ref:`cql/declaring_models:base models`.

Auto Migration
----------------------------------------------------------

To persist the models it is necessary to migrate the database, 
so that the structure of the tables corresponds to the definition of the model. 
This migration is performed by gorm through the gormDB.

For details visit :ref:`cql/connecting_to_a_database:migration`.

GormDB
-----------------------------

GormDB is a gorm.DB object that allows communication with the database. 
This object will be needed as a parameter for the main cql functions (Query, Update and Delete).

For details visit :ref:`cql/connecting_to_a_database:connection`.

Condition
-----------------------------

Conditions are the basis of the cql query system, every query is composed of a set of conditions. 
Conditions belong to a particular model and there are 4 different types: 
WhereConditions, ConnectionConditions and JoinConditions.

For details visit :doc:`/cql/query`.

WhereCondition
-----------------------------

Type of condition that allows filters to be made on the model to which they belong 
and an attribute of this model. These filters are performed through operators.

For details visit :doc:`/cql/query`.

ConnectionCondition
-----------------------------

Type of condition that allows the use of logical operators 
(and, or, or, not) between WhereConditions.

For details visit :doc:`/cql/query`.

JoinCondition
-----------------------------

Condition type that allows to navigate relationships between models, 
which will result in a join in the executed query 
(don't worry, if you don't know what a join is, 
you don't need to understand the queries that cql executes).

For details visit :doc:`/cql/query`.

Operator
-----------------------------

Concept similar to database operators, 
which allow different operations to be performed on an attribute of a model, 
such as comparisons, predicates, pattern matching, etc.

Operators can be classified as static, dynamic and unsafe.

For details visit :doc:`/cql/query`.

Static operator
-----------------------------

Static operators are those that perform operations on an attribute and static values, 
such as a boolean value, an integer, etc.

For details visit :doc:`/cql/query`.

Dynamic operator
-----------------------------

Dynamic operators are those that perform operations between an attribute and other attributes, 
either from the same model or from a different model, as long as the type of these attributes is the same.

For details visit :doc:`/cql/advanced_query`.

Unsafe operator
-----------------------------

Unsafe operators are those that can perform operations between an attribute and 
any type of value or attribute.

For details visit :doc:`/cql/advanced_query`.

Nullable types
-----------------------------

Nullable types are the types provided by the sql library 
that are a nullable version of the basic types: 
sql.NullString, sql.NullTime, sql.NullInt64, sql.NullInt32, 
sql.NullBool, sql.NullFloat64, etc..

For details visit <https://pkg.go.dev/database/sql>.

Compiled query system
-----------------------------

The set of conditions that are received by the 
`cql.Query`, `cql.Update` and `cql.Delete` methods form the cql compiled query system. 
It is so named because the conditions will verify at compile time that the query to be executed is correct.

For details visit :ref:`cql/query:conditions` and :doc:`/cql/compile_time_safety`.

Conditions generation
----------------------------

Conditions are the basis of the compiled query system. 
They are generated for each model and attribute and can then be used. 
Their generation is done with cql-gen.

For details visit :ref:`cql/cqlgen:Conditions generation`.

Relation getter
-----------------------------------

Relationships between objects can be loaded from the database using the Preload method. 
In order to safely navigate the relations in the loaded model cql provides methods 
called "relation getters".

For details visit :doc:`/cql/preloading`.