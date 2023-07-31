# BaDaaS ORM: Backend and Distribution ORM (Object Relational Mapping) <!-- omit in toc -->

BaDaaS ORM is the BaDaaS component that allows for easy persistence and querying of objects. It is built on top of gorm and adds for each entity a service and a repository that allows complex queries without any extra effort.

BaDaaS ORM can be used both within a BaDaaS application and as a stand-alone application.

- [Installation](#installation)
- [Provided functionalities](#provided-functionalities)
  - [Base models](#base-models)
  - [CRUDServiceModule](#crudservicemodule)

## Installation

Once you have started your project with `go init`, you must add the dependency to BaDaaS:

```bash
go get -u github.com/ditrit/badaas
```

## Provided functionalities

### Base models

badaas-orm gives you two types of base models for your classes: `orm.UUIDModel` and `orm.UIntModel`.

To use them, simply embed the desired model in any of your classes:

```go
type MyClass struct {
  orm.UUIDModel

  // your code here
}
```

Once done your class will be considered a **BaDaaS Model**.

The difference between them is the type they will use as primary key: a random uuid and an auto incremental uint respectively. Both provide date created, edited and deleted (<https://gorm.io/docs/delete.html#Soft-Delete>).

### CRUDServiceModule

`CRUDServiceModule` provides you a CRUDService and a CRUDRepository for your badaas Model. After calling it as, for example, `orm.GetCRUDServiceModule[models.Company](),` the following can be used by dependency injection:

- `crudCompanyService orm.CRUDService[models.Company, orm.UUID]`
- `crudCompanyRepository orm.CRUDRepository[models.Company, orm.UUID]`

These classes will allow you to perform queries using the compilable query system generated with badaas-cli. For details on how to do this visit [badaas-cli docs](github.com/ditrit/badaas-cli/README.md).
