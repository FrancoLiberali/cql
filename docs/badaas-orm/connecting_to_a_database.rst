==============================
Connecting to a database
==============================

Connection
-----------------------------

badaas-orm supports the PostgreSQL databases using gorm's driver. 
Some databases may be compatible with the postgres dialect, 
in which case you could just use the dialect for those databases (from which CockroachDB is tested).

To communicate with the database badaas-orm need a :ref:`GormDB <badaas-orm/concepts:GormDB>` object, 
that can be created by following `gorm documentation <https://gorm.io/docs/connecting_to_the_database.html>`_. 

badaas-orm also offers the `orm.ConnectToDialector` method that will allow you to connect to a database 
using the specified dialector with retrying. 
It also configures the `gorm's logger <https://gorm.io/docs/logger.html>`_ to work with 
`zap logger <https://github.com/uber-go/zap>`_.

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



