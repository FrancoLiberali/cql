==============================
CRUD Operations
==============================

CRUDServices and CRUDRepositories
--------------------------------------

CRUD operations are made to your models via the CRUDServices and CRUDRepositories. 
The difference between them is that a CRUDService will execute this operations within a transaction 
while the CRUDRepository will be executed within a transaction received by parameter, 
thus allowing defining services that perform multiple CRUD operations within the same transaction.

Create, Save and Delete methods are just hooks to the gorm's corresponding methods. 
For details visit 
<https://gorm.io/docs/create.html>, <https://gorm.io/docs/update.html> and <https://gorm.io/docs/delete.html>. 
On the other hand, read (query) operations are provided by badaas-orm via its 
:ref:`compilable query system <badaas-orm/concepts:compilable query system>` 
(see how in :doc:`/badaas-orm/query`).

Each pair of CRUDService and CRUDRepository corresponds to a model. To create them you must use 
the `orm.GetCRUD[<model>, <modelID>](gormDB)` where 
`<model>` is the type of your :ref:`model <badaas-orm/concepts:model>`, 
`<modelID>` is the type of your :ref:`model's id <badaas-orm/concepts:model id>` 
and `gormDB` is the :ref:`GormDB <badaas-orm/concepts:GormDB>` object.

When using badaas-orm with `fx` as :ref:`dependency injector <badaas-orm/concepts:Dependency injection>` you 
will need to provide to fx `orm.GetCRUDServiceModule[<model>]()` 
where `<model>` is the type of your :ref:`model <badaas-orm/concepts:model>`. 
After that the following can be used by dependency injection:

- `crudYourModelService orm.CRUDService[<model>, <modelID>]`
- `crudYourModelRepository orm.CRUDRepository[<model>, <modelID>]`

For example:

.. code-block:: go


    type YourModel struct {
        model.UUIDModel
    }

    func main() {
        fx.New(
            // connect to db
            fx.Provide(NewGormDBConnection),
            // activate badaas-orm
            fx.Provide(GetModels),
            orm.AutoMigrate,

            orm.GetCRUDServiceModule[YourModel](),
            fx.Invoke(QueryCRUDObjects),
        ).Run()
    }

    func QueryCRUDObjects(crudYourModelService orm.CRUDService[YourModel, model.UUID]) {
        // use crudYourModelService
    }