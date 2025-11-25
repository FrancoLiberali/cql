==============================
Type safety
==============================

One of the most important features of the CQL is

.. code-block:: none

    Is compile-time safe:
        queries are validated at compile time to avoid errors 
        such as comparing attributes that are of different types, 
        trying to use attributes or navigate relationships that do not exist, 
        using information from tables that are not included in the query, etc.; 
        ensuring that a runtime error will not be raised.

While there are other libraries that provide an API type safety 
(`gorm-gen <https://gorm.io/gen/>`_, `jooq <https://www.jooq.org/>`_ (Java), 
`diesel <https://diesel.rs/>`_ (Rust)), CQL is the only one that allows us to be sure 
that the generated query is correct, avoiding runtime errors.

Each of the CQL features is designed to be safe at compile time.
For details of each feature, see:

- Query: :doc:`/cql/query_type_safety`
- Select: :ref:`cql/select:Type safety`.
- Insert: :ref:`cql/insert:Type safety`.
- Update: :ref:`cql/update:Type safety`.
- Delete: :ref:`cql/delete:Type safety`.

In each of these sections, you will see that there are limitations in terms of runtime security, 
with borderline cases where queries can generate runtime errors even though they compile. 
For these cases, there is :doc:`/cql/cqllint`, a utility for analyzing our code and finding these cases, 
to ensure security before runtime.

A runtime error that is common to almost all CQL methods is cql.ErrUnsupportedByDatabase.
This error is generated when an attempt is made to use a method or function that is not supported by the database engine used. 
When using cql, you will see many comments about features that certain database engines do not support:

.. code-block:: go

    // SetMultiple allows updating multiple tables in the same query.
    //
    // available for: mysql
    func (update *Update[T]) SetMultiple(sets ...ISet) (int64, error) {

generating the following error when used in another database:

.. code-block:: none

    method not supported by database; method: SetMultiple

This error is not yet supported by cqllint, so these comments should be reviewed carefully.
