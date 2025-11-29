==============================
Connecting to a database
==============================

Connection
-----------------------------

cql supports the databases MySQL, PostgreSQL, SQLite, SQL Server using gorm's driver. 
Some databases may be compatible with the mysql or postgres dialect, 
in which case you could just use the dialect for those databases.

To communicate with the database, cql needs a :ref:`cql.DB <cql/concepts:cqlDB>` object. 
To create it, you can use the function `cql.Open` that will allow you to connect to a database 
using the specified dialector. 

internally, this function uses `gorm.Open` 
and can receive both cql configuration, such as the logger, and gorm configuration. 
For details about this logger visit :doc:`/cql/logger`. 
For details about gorm configuration visit `gorm documentation <https://gorm.io/docs/connecting_to_the_database.html>`_.

Auto Migration
----------------------------

Auto Migration can be done by gorm using the `gormDB.AutoMigrate` method. 
For details visit `gorm docs <https://gorm.io/docs/migration.html>`_.
