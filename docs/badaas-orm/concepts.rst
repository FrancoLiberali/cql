==============================
Concepts
==============================

Model
------------------------------

A model is any object (struct) of go that you want to persist 
in the database and on which you can perform queries. 
For this, the struct must have an embedded badaas-orm base model.

For details visit :ref:`badaas-orm/declaring_models:model declaration`.

Base model
-----------------------------

It is a struct that when embedded allows your structures to become BaDaaS models, 
adding attributes ID, CreatedAt, UpdatedAt and DeletedAt attributes and the possibility to persist, 
create conditions and perform queries on these structures.

For details visit :ref:`badaas-orm/declaring_models:base models`.

Model ID
-----------------------------

The id is a unique identifier needed to persist a model in the database. 
It can be a model.UIntID or a model.UUID, depending on the base model used.

For details visit :ref:`badaas-orm/declaring_models:base models`.

Auto Migration
----------------------------------------------------------

To persist the models it is necessary to migrate the database, 
so that the structure of the tables corresponds to the definition of the model. 
This migration is performed by gorm through the gormDB.

For details visit :ref:`badaas-orm/connecting_to_a_database:migration`.

GormDB
-----------------------------

GormDB is a gorm.DB object that allows communication with the database.

For details visit :ref:`badaas-orm/connecting_to_a_database:connection`.

Condition
-----------------------------

Conditions are the basis of the badaas-orm query system, every query is composed of a set of conditions. 
Conditions belong to a particular model and there are 4 different types: 
WhereConditions, ConnectionConditions, JoinConditions and PreloadConditions.

For details visit :doc:`/badaas-orm/query`.

WhereCondition
-----------------------------

Type of condition that allows filters to be made on the model to which they belong 
and an attribute of this model. These filters are performed through operators.

For details visit :doc:`/badaas-orm/query`.

ConnectionCondition
-----------------------------

Type of condition that allows the use of logical operators 
(and, or, or, not) between WhereConditions.

For details visit :doc:`/badaas-orm/query`.

JoinCondition
-----------------------------

Condition type that allows to navigate relationships between models, 
which will result in a join in the executed query 
(don't worry, if you don't know what a join is, 
you don't need to understand the queries that badaas-orm executes).

For details visit :doc:`/badaas-orm/query`.

PreloadCondition
-----------------------------

Type of condition that allows retrieving information from a model as a result of the database (preload). 
This information can be all its attributes and/or another model that is related to it.

For details visit :doc:`/badaas-orm/preloading`.

Operator
-----------------------------

Concept similar to database operators, 
which allow different operations to be performed on an attribute of a model, 
such as comparisons, predicates, pattern matching, etc.

Operators can be classified as static, dynamic and unsafe.

For details visit :doc:`/badaas-orm/query`.

Static operator
-----------------------------

Static operators are those that perform operations on an attribute and static values, 
such as a boolean value, an integer, etc.

For details visit :doc:`/badaas-orm/query`.

Dynamic operator
-----------------------------

Dynamic operators are those that perform operations between an attribute and other attributes, 
either from the same model or from a different model, as long as the type of these attributes is the same.

For details visit :doc:`/badaas-orm/advanced_query`.

Unsafe operator
-----------------------------

Unsafe operators are those that can perform operations between an attribute and 
any type of value or attribute.

For details visit :doc:`/badaas-orm/advanced_query`.

Nullable types
-----------------------------

Nullable types are the types provided by the sql library 
that are a nullable version of the basic types: 
sql.NullString, sql.NullTime, sql.NullInt64, sql.NullInt32, 
sql.NullBool, sql.NullFloat64, etc..

For details visit <https://pkg.go.dev/database/sql>.

CRUDService
-----------------------------

A CrudService is a service that allows us to perform CRUD (create, read, update and delete) 
operations on a specific model, executing all the necessary operations within a transaction. 
Internally they use the CRUDRepository of that model.

For details visit :ref:`badaas-orm/crud:CRUDServices and CRUDRepositories`.

CRUDRepository
-----------------------------

A CRUDRepository is an object that allows us to perform CRUD operations (create, read, update, delete) 
on a model but, unlike services, its internal operations are performed within a transaction received 
by parameter. 
This is useful to be able to define services that perform multiple CRUD 
operations within the same transaction.

For details visit :ref:`badaas-orm/crud:CRUDServices and CRUDRepositories`.

Compilable query system
-----------------------------

The set of conditions that are received by the read operations of the CRUDService 
and CRUDRepository form the badaas-orm compilable query system. 
It is so named because the conditions will verify at compile time that the query to be executed is correct.

For details visit :ref:`badaas-orm/query:compilable query system`.

Conditions generation
----------------------------

Conditions are the basis of the compilable query system. 
They are generated for each model and attribute and can then be used. 
Their generation is done with badaas-cli.

For details visit :ref:`badaas-orm/query:Conditions generation`.

Dependency injection
-----------------------------------

Dependency injection is a programming technique in which an object or function 
receives other objects or functions that it depends on. badaas-orm is compatible with 
`uber fx <https://uber-go.github.io/fx/>`_ to inject the CRUDServices and 
CRUDRepositories in your objects and functions.

Relation getter
-----------------------------------

Relationships between objects can be loaded from the database using PreloadConditions. 
In order to safely navigate the relations in the loaded model badaas-orm provides methods 
called "relation getters".

For details visit :doc:`/badaas-orm/preloading`.