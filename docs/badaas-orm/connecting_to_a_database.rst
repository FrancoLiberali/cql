==============================
Connecting to a database
==============================

Connection
-----------------------------

badaas-orm supports the databases MySQL, PostgreSQL, SQLite, SQL Server using gorm's driver. 
Some databases may be compatible with the mysql or postgres dialect, 
in which case you could just use the dialect for those databases (from which CockroachDB is tested).

To communicate with the database badaas-orm need a :ref:`GormDB <badaas-orm/concepts:GormDB>` object. 
To create it, you can use the function `orm.Open` that will allow you to connect to a database 
using the specified dialector. This function is equivalent to `gorm.Open` 
but with the difference that in case of not adding any configuration, 
the badaas-orm default logger will be configured instead of the gorm one. 
For details about this logger visit :doc:`/badaas-orm/logger`. 
For details about gorm configuration visit `gorm documentation <https://gorm.io/docs/connecting_to_the_database.html>`_.

Migration
----------------------------

Migration is done by gorm using the `gormDB.AutoMigrate` method. 
For details visit `gorm docs <https://gorm.io/docs/migration.html>`_.