==============================
Connecting to a database
==============================

Connection
-----------------------------

badaas-orm supports the PostgreSQL databases using gorm's driver. 
Some databases may be compatible with the postgres dialect, 
in which case you could just use the dialect for those databases (from which CockroachDB is tested).

To communicate with the database badaas-orm need a :ref:`GormDB <badaas-orm/concepts:GormDB>` object. 
To create it, you can use the function `orm.Open` that will allow you to connect to a database 
using the specified dialector. This function is equivalent to `gorm.Open` 
but with the difference that in case of not adding any configuration, 
the badaas-orm default logger will be configured instead of the gorm one. 
For details about this logger visit :doc:`/badaas-orm/logger`. 
For details about gorm configuration visit `gorm documentation <https://gorm.io/docs/connecting_to_the_database.html>`_.

When using badaas-orm with `fx` as :ref:`dependency injector <badaas-orm/concepts:Dependency injection>` you 
will need to provide (`fx.Provide`) a function that returns a `*gorm.DB`.

Migration
----------------------------

Migration is done by gorm using the `gormDB.AutoMigrate` method. 
For details visit `gorm docs <https://gorm.io/docs/migration.html>`_.

When using badaas-orm with `fx` as :ref:`dependency injector <badaas-orm/concepts:Dependency injection>` 
this method can't be called directly. In that case, badaas-orm will execute the migration by providing 
`orm.AutoMigrate` to fx. For this to work, you will need to provide also a method that returns 
`orm.GetModelsResult` with the models you want to include in the migration. 
Remember that the order in this list is important for gorm to be able to execute the migration.