# CQL: Compiled Query Language <!-- omit in toc -->

[![Build Status](https://github.com/FrancoLiberali/cql/actions/workflows/cql.yml/badge.svg?branch=main)](https://github.com/FrancoLiberali/cql/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/FrancoLiberali/cql)](https://goreportcard.com/report/github.com/FrancoLiberali/cql)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=FrancoLiberali_cql&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=FrancoLiberali_cql)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=FrancoLiberali_cql&metric=coverage)](https://sonarcloud.io/summary/new_code?id=FrancoLiberali_cql)

[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/FrancoLiberali/cql)
[![Documentation Status](https://readthedocs.org/projects/compiledquerylenguage/badge/?version=latest)](https://compiledquerylenguage.readthedocs.io/en/latest/?badge=latest)

Originally part of [BaDaaS](https://github.com/ditrit/badaas), CQL allows easy and safe persistence and querying of objects.

It's built on top of [gorm](https://gorm.io/), a library that actually provides the functionality of an ORM: mapping objects to tables in the SQL database. While gorm does this job well with its automatic migration then performing queries on these objects is somewhat limited, forcing us to write SQL queries directly when they are complex. CQL seeks to address these limitations with a query system that:

- Is compile-time safe: its query system is validated at compile time to avoid errors such as comparing attributes that are of different types, trying to use attributes or navigate relationships that do not exist, using information from tables that are not included in the query, etc.
- Is easy to use: the use of this system does not require knowledge of databases, SQL languages or complex concepts. Writing queries only requires programming in go and the result is easy to read.
- Is designed for real applications: the query system is designed to work well in real-world cases where queries are complex, require navigating multiple relationships, performing multiple comparisons, etc.
- Is designed so that developers can focus on the business model: its queries allow easy retrieval of model relationships to apply business logic to the model and it provides mechanisms to avoid errors in the business logic due to mistakes in loading information from the database.
- It is designed for high performance: the query system avoids as much as possible the use of reflection and aims that all the necessary model data can be retrieved in a single query to the database.

<style>
table th:first-of-type {
    width: 33%;
}
table th:nth-of-type(2) {
    width: 33%;
}
table th:nth-of-type(3) {
    width: 33%;
}
</style>

| Language | Query |
|---|---|
| SQL | SELECT cities.* FROM cities <br> INNER JOIN countries ON <br>&emsp;&emsp; countries.id = cities.country_id AND <br>&emsp;&emsp; countries.name = "France" <br> WHERE cities.name = "Paris" |
| GORM | db.Where(<br>&emsp;"cities.name = ?",<br>&emsp;"Paris",<br>).Joins(<br>&emsp;"Country",<br>&emsp;db.Where( <br>&emsp;&emsp; "Country.name = ?", <br>&emsp;&emsp; "France", <br>&emsp; ), <br> ).Find(&cities) |
| CQL | cql.Query[models.City]( <br>&emsp; db, <br>&emsp; conditions.City.Name.Is().Eq("Paris"), <br>&emsp; conditions.City.Country( <br>&emsp;&emsp; conditions.Country.Name.Is().Eq("France"), <br>&emsp; ), <br> ).FindOne() |

## Documentation

<https://compiledquerylenguage.readthedocs.io/en/latest/>

## Contributing

See [this section](../docs/contributing/contributing.md) to view the cql contribution guidelines.

## License

CQL is Licensed under the [Mozilla Public License Version 2.0](../LICENSE).
