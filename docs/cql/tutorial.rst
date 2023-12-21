==============================
Tutorial
==============================

In this short tutorial you will learn the main functionalities of cql. 
The code to be executed in each step can be found in this `repository <https://github.com/FrancoLiberali/cql-tutorial>`_.

Model and data
-----------------------

In the file `models/models.go` you find the definition of the following model:

.. image:: /img/cql-tutorial-model.png
  :width: 700
  :alt: cql tutorial model

For details about the definition of models you can read :doc:`/cql/declaring_models`.

In `sqlite:db` you will find a sqlite database with the following data:

.. list-table:: Countries
   :header-rows: 1

   * - ID
     - Name
     - CapitalID
   * - 3739a825-bc5c-4350-a2bc-6e77e22fe3f4
     - France
     - eaa480a3-694e-4be3-9af5-ad935cdd57e2
   * - 0c4404f6-83c2-4bdf-93d5-a5ff2fe4f921
     - United States of America
     - df44272e-c3db-4e18-876c-f9f579488716

.. list-table:: Cities
   :header-rows: 1

   * - ID
     - Name
     - Population
     - CountryID
   * - eaa480a3-694e-4be3-9af5-ad935cdd57e2
     - Paris
     - 2161000
     - 3739a825-bc5c-4350-a2bc-6e77e22fe3f4
   * - df44272e-c3db-4e18-876c-f9f579488716
     - Washington D. C.
     - 689545
     - 0c4404f6-83c2-4bdf-93d5-a5ff2fe4f921
   * - 8c3dfc38-1fc6-4ec9-a89b-e41018a54b4a
     - Paris
     - 25171
     - 0c4404f6-83c2-4bdf-93d5-a5ff2fe4f921

As you can see, there are two cities called Paris in this database: 
the well known Paris, capital of France and site of the iconic Eiffel tower, 
and Paris in the United States of America, site of the Eiffel tower with the cowboy hat 
(no joke, just search for paris texas eiffel tower in your favorite search engine).

In this tutorial we will explore the cql functions that will allow us to differentiate these two Paris.

Tutorial 1: simple query
-------------------------------

In this first tutorial we are going to perform a simple query to obtain all the cities called Paris. 

In the tutorial_1.go file you will find that we can perform this query as follows:

.. code-block:: go

    cities, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
    ).Find()

We can run this tutorial with `make tutorial_1` and we will obtain the following result:

.. code-block:: bash

    Cities named 'Paris' are:
        1: &{UUIDModel:{ID:eaa480a3-694e-4be3-9af5-ad935cdd57e2 CreatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 UpdatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Name:Paris Population:2161000 Country:<nil> CountryID:3739a825-bc5c-4350-a2bc-6e77e22fe3f4}
        2: &{UUIDModel:{ID:8c3dfc38-1fc6-4ec9-a89b-e41018a54b4a CreatedAt:2023-08-11 16:43:27.468149185 +0200 +0200 UpdatedAt:2023-08-11 16:43:27.468149185 +0200 +0200 DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Name:Paris Population:25171 Country:<nil> CountryID:0c4404f6-83c2-4bdf-93d5-a5ff2fe4f921}

As you can see, in this case we will get both cities which we can differentiate by their population and the id of the country.

In this tutorial we have used the cql compiled queries system to get these cities, 
for more details you can read :ref:`cql/query:conditions`.

Tutorial 2: operators
-------------------------------

Now we are going to try to obtain only the Paris of France and in a first 
approximation we could do it using its population: we will only look for the Paris 
whose population is greater than one million inhabitants. 

In the tutorial_2.go file you will find that we can perform this query as follows:

.. code-block:: go

    cities, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
        conditions.City.Population.Is().Gt(1000000),
    ).Find()

We can run this tutorial with `make tutorial_2` and we will obtain the following result:

.. code-block:: bash

    Cities named 'Paris' with a population bigger than 1.000.000 are:
        1: &{UUIDModel:{ID:eaa480a3-694e-4be3-9af5-ad935cdd57e2 CreatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 UpdatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Name:Paris Population:2161000 Country:<nil> CountryID:3739a825-bc5c-4350-a2bc-6e77e22fe3f4}

As you can see, in this case we only get one city, Paris in France.

In this tutorial we have used the operator Gt to obtain this city, 
for more details you can read :ref:`cql/query:Operators`.

Tutorial 3: modifiers
-------------------------------

Although in the previous tutorial we achieved our goal of differentiating the two Paris, 
the way to do it is debatable since the population of Paris, Texas may increase to over 1000000 someday 
and then, the result of this query can change. 
Therefore, we will search only for the city with the largest population.

In the tutorial_3.go file you will find that we can perform this query as follows:

.. code-block:: go

    parisFrance, err := cql.Query[models.City](
		db,
		conditions.City.Name.Is().Eq("Paris"),
	).Descending(
		conditions.City.Population,
	).Limit(1).FindOne()

We can run this tutorial with `make tutorial_3` and we will obtain the following result:

.. code-block:: bash

    City named 'Paris' with the largest population is: &{UUIDModel:{ID:eaa480a3-694e-4be3-9af5-ad935cdd57e2 CreatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 UpdatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Name:Paris Population:2161000 Country:<nil> CountryID:3739a825-bc5c-4350-a2bc-6e77e22fe3f4}

As you can see, again we get only the Paris in France. 
As you may have noticed, in this case we have used the `FindOne` method instead of `Find`. 
This is because in this case we are sure that the result is a single model, 
so instead of getting a list we get a single city.

In this tutorial we have used query modifier methods, 
for more details you can read :ref:`cql/query:Query methods`.

Tutorial 4: joins
-------------------------------

Again, the solution of the previous tutorial is debatable because the evolution 
of populations could make Paris, Texas have more inhabitants than Paris, France one day. 
Therefore, we are now going to improve this query by obtaining the city called 
Paris whose country is called France. 

In the tutorial_4.go file you will find that we can perform this query as follows:

.. code-block:: go

    parisFrance, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
        conditions.City.Country(
            conditions.Country.Name.Is().Eq("France"),
        ),
    ).FindOne()

We can run this tutorial with `make tutorial_4` and we will obtain the following result:

.. code-block:: bash

    Cities named 'Paris' in 'France' are:
        1: &{UUIDModel:{ID:eaa480a3-694e-4be3-9af5-ad935cdd57e2 CreatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 UpdatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Name:Paris Population:2161000 Country:<nil> CountryID:3739a825-bc5c-4350-a2bc-6e77e22fe3f4}

As you can see, again we get only the Paris in France. 

In this tutorial we have used a condition that performs a join, 
for more details you can read :ref:`cql/query:Use of the conditions`.

Tutorial 5: preloading
-------------------------------

You may have noticed that in the results of the previous tutorials the Country field of the cities was null (Country:<nil>). 
This is because, to ensure performance, cql will retrieve only the attributes of the model 
you are querying (City in this case because the method used is cql.Query[models.City]) 
but not of its relationships. If we also want to obtain this data, we must perform preloading.

In the tutorial_5.go file you will find that we can perform this query as follows:

.. code-block:: go

    cities, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
        conditions.City.PreloadCountry(),
    ).Find()

We can run this tutorial with `make tutorial_5` and we will obtain the following result:

.. code-block:: bash

    Cities named 'Paris' are:
        1: &{UUIDModel:{ID:eaa480a3-694e-4be3-9af5-ad935cdd57e2 CreatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 UpdatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Name:Paris Population:2161000 Country:0xc0001d1600 CountryID:3739a825-bc5c-4350-a2bc-6e77e22fe3f4}
            with country: &{UUIDModel:{ID:3739a825-bc5c-4350-a2bc-6e77e22fe3f4 CreatedAt:2023-08-11 16:43:27.445202858 +0200 +0200 UpdatedAt:2023-08-11 16:43:27.457191337 +0200 +0200 DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Name:France Capital:<nil> CapitalID:eaa480a3-694e-4be3-9af5-ad935cdd57e2}
        2: &{UUIDModel:{ID:8c3dfc38-1fc6-4ec9-a89b-e41018a54b4a CreatedAt:2023-08-11 16:43:27.468149185 +0200 +0200 UpdatedAt:2023-08-11 16:43:27.468149185 +0200 +0200 DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Name:Paris Population:25171 Country:0xc0001d1780 CountryID:0c4404f6-83c2-4bdf-93d5-a5ff2fe4f921}
            with country: &{UUIDModel:{ID:0c4404f6-83c2-4bdf-93d5-a5ff2fe4f921 CreatedAt:2023-08-11 16:43:27.462357133 +0200 +0200 UpdatedAt:2023-08-11 16:43:27.479800337 +0200 +0200 DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Name:United States of America Capital:<nil> CapitalID:df44272e-c3db-4e18-876c-f9f579488716}

As you can see, now the country attribute is a valid pointer to a Country object (Country:0xc0001d1600).
Then the Country object information is accessed with the `GetCountry` method. 
This method is not defined in the `models/models.go` file but is a :ref:`relation getter <cql/concepts:relation getter>` 
that is generated by cql-cli together with the conditions. 
These methods allow us to differentiate null objects from objects not loaded from the database, 
since when trying to browse a relation that was not loaded we will get `errors.ErrRelationNotLoaded`. 

In this tutorial we have used preloading and relation getters, 
for more details you can read :doc:`/cql/preloading`.

Tutorial 6: dynamic operators
-------------------------------

So far we have performed operations that take as input a static value (equal to "Paris" or greater than 1000000) 
but what if now we would like to differentiate these two Paris from each other based on whether they 
are the capital of their country.

In the tutorial_6.go file you will find that we can perform this query as follows:

.. code-block:: go

    cities, err := cql.Query[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
        conditions.City.Country(
            conditions.Country.CapitalID.Is().Dynamic().Eq(conditions.City.ID),
        ),
    ).Find()

We can run this tutorial with `make tutorial_6` and we will obtain the following result:

.. code-block:: bash

    Cities named 'Paris' that are the capital of their country are:
        1: &{UUIDModel:{ID:eaa480a3-694e-4be3-9af5-ad935cdd57e2 CreatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 UpdatedAt:2023-08-11 16:43:27.451393348 +0200 +0200 DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Name:Paris Population:2161000 Country:<nil> CountryID:3739a825-bc5c-4350-a2bc-6e77e22fe3f4}

As you can see, again we only get the Paris in France.

In this tutorial we have used dynamic conditions, 
for more details you can read :ref:`cql/advanced_query:Dynamic operators`.

Tutorial 7: update
-------------------------------

So far we have only made select queries, but in this tutorial we want to edit the population of Paris.

In the tutorial_7.go file you will find that we can perform this query as follows:

.. code-block:: go

    updated, err := cql.Update[models.City](
        db,
        conditions.City.Name.Is().Eq("Paris"),
        conditions.City.Country(
            conditions.Country.Name.Is().Eq("France"),
        ),
    ).Returning(&cities).Set(
        conditions.City.Population.Set().Eq(2102650),
    )

We can run this tutorial with `make tutorial_7` and we will obtain the following result:

.. code-block:: bash

    Updated 1 city: {{eaa480a3-694e-4be3-9af5-ad935cdd57e2 2023-08-11 16:43:27.451393348 +0200 +0200 2023-12-21 10:02:36.420763701 -0300 -0300 {0001-01-01 00:00:00 +0000 UTC false}} Paris 2102650 <nil> 3739a825-bc5c-4350-a2bc-6e77e22fe3f4}
    Initial population was 2161000

As you can see, first we can know the number of updated models with the value "updated" returned by the Set method 
(according to the number of models that meet the conditions entered in the Update method). 
On the other hand, it is also possible to obtain the information of the updated models using the Returning method.

In this tutorial we have used updates, 
for more details you can read :ref:`cql/update`.

Tutorial 8: create and delete
-------------------------------

In this tutorial we want to create a new city called Rennes and then delete it.

In the tutorial_8.go file you will find that we can perform this query as follows:

Create:

.. code-block:: go

    rennes := models.City{
        Country:    france,
        Name:       "Rennes",
        Population: 215366,
    }
    if err := db.Create(&rennes).Error; err != nil {
        log.Panicln(err)
    }

Delete:

.. code-block:: go

    deleted, err := cql.Delete[models.City](
		db,
		conditions.City.Name.Is().Eq("Rennes"),
	).Exec()

We can run this tutorial with `make tutorial_8` and we will obtain the following result:

.. code-block:: bash

    Deleted 1 city

Here, we simply get the number of deleted models through the "deleted" variable returned by the Exec method 
(according to the number of models that meet the conditions entered in the Delete method).

In this tutorial we have used create and delete, 
for more details you can read :ref:`cql/create` and :ref:`cql/delete`.
