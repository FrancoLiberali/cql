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
   * - 1
     - United States of America
     - 2
   * - 2
     - France
     - 3

.. list-table:: Cities
   :header-rows: 1

   * - ID
     - Name
     - Population
     - CountryID
   * - 1
     - Paris
     - 25171
     - 1
   * - 2
     - Washington D. C.
     - 689545
     - 1
   * - 3
     - Paris
     - 2161000
     - 2

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
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
    ).Find()

We can run this tutorial with `make tutorial_1` and we will obtain the following result:

.. code-block:: none

    Cities named 'Paris' are:
        1: City{ID: 1, Name: Paris, Population: 25171, CountryID:1, Country:<nil> }
        2: City{ID: 3, Name: Paris, Population: 2161000, CountryID:2, Country:<nil> }

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
    :emphasize-lines: 5

    cities, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
        conditions.City.Population.Is().Gt(cql.Int64(1000000)),
    ).Find()

We can run this tutorial with `make tutorial_2` and we will obtain the following result:

.. code-block:: none

    Cities named 'Paris' with a population bigger than 1.000.000 are:
        1: City{ID: 3, Name: Paris, Population: 2161000, CountryID:2, Country:<nil> }

As you can see, in this case we only get one city, Paris in France.

In this tutorial we have used the operator Gt to obtain this city, 
for more details you can read :ref:`cql/query:Operators`.

Tutorial 3: modifiers
-------------------------------

Although in the previous tutorial we achieved our goal of differentiating the two Paris, 
the way to do it is debatable since the population of Paris, Texas may increase to over 1.000.000 someday
and then, the result of this query can change. 
Therefore, we will search only for the city with the largest population.

In the tutorial_3.go file you will find that we can perform this query as follows:

.. code-block:: go
    :emphasize-lines: 5,6,7

    parisFrance, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
    ).Descending(
        conditions.City.Population,
    ).Limit(1).FindOne()

We can run this tutorial with `make tutorial_3` and we will obtain the following result:

.. code-block:: none

    City named 'Paris' with the largest population is: City{ID: 3, Name: Paris, Population: 2161000, CountryID:2, Country:<nil> }

As you can see, again we get only the Paris in France. 
As you may have noticed, in this case we have used the `FindOne` method instead of `Find`. 
This is because in this case we are sure that the result is a single model, 
so instead of getting a list we get a single city.

In this tutorial we have used query modifier methods, 
for more details you can read :ref:`cql/query:Query methods`.

Tutorial 4: functions
-------------------------------

Another alternative could be to try applying functions to the values of the cities to determine which one
is Paris, France.
As an example, let's look for cities where twice the population is greater than 1.000.000.

In the tutorial_4.go file you will find that we can perform this query as follows:

.. code-block:: go
    :emphasize-lines: 5

    cities, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
        conditions.City.Population.Times(cql.Int64(2)).Is().Gt(cql.Int64(1000000)),
    ).Find()

We can run this tutorial with `make tutorial_4` and we will obtain the following result:

.. code-block:: none

    Cities named 'Paris' with twice its population bigger than 1.000.000 are:
        1: City{ID: 3, Name: Paris, Population: 2161000, CountryID:2, Country:<nil> }

As you can see, in this case we only get one city, Paris in France.

In this tutorial we have used the function Times to multiply the population of cities, 
for more details you can read :ref:`cql/query:Functions`.

Tutorial 5: joins
-------------------------------

Again, the solution of the previous tutorial is debatable because the evolution 
of populations could make Paris, Texas have more inhabitants than Paris, France one day. 
Therefore, we are now going to improve this query by obtaining the city called 
Paris whose country is called France. 

In the tutorial_5.go file you will find that we can perform this query as follows:

.. code-block:: go
    :emphasize-lines: 5,6,7

    parisFrance, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
        conditions.City.Country(
            conditions.Country.Name.Is().Eq(cql.String("France")),
        ),
    ).FindOne()

We can run this tutorial with `make tutorial_5` and we will obtain the following result:

.. code-block:: none

    City named 'Paris' in 'France' is: City{ID: 3, Name: Paris, Population: 2161000, CountryID:2, Country:<nil> }

As you can see, again we get only the Paris in France. 

In this tutorial we have used a condition that performs a join.

Tutorial 6: preloading
-------------------------------

You may have noticed that in the results of the previous tutorials the Country field of the cities was null (Country:<nil>). 
This is because, to ensure performance, cql will retrieve only the attributes of the model 
you are querying (City in this case because the method used is cql.Query[models.City]) 
but not of its relationships. If we also want to obtain this data, we must perform preloading.

In the tutorial_6.go file you will find that we can perform this query as follows:

.. code-block:: go
    :emphasize-lines: 5

    cities, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
        conditions.City.Country().Preload(),
    ).Find()

We can run this tutorial with `make tutorial_6` and we will obtain the following result:

.. code-block:: none

    Cities named 'Paris' are:
        1: City{ID: 1, Name: Paris, Population: 25171, CountryID:1, Country:Country{ID: 1, Name: United States of America, CapitalID:2, Capital:<nil> } } with country: Country{ID: 1, Name: United States of America, CapitalID:2, Capital:<nil> }
        2: City{ID: 3, Name: Paris, Population: 2161000, CountryID:2, Country:Country{ID: 2, Name: France, CapitalID:3, Capital:<nil> } } with country: Country{ID: 2, Name: France, CapitalID:3, Capital:<nil> }

As you can see, now the country attribute is a valid pointer to a Country object.
Then the Country object information is accessed with the `GetCountry` method. 
This method is not defined in the `models/models.go` file but is a :ref:`relation getter <cql/concepts:relation getter>` 
that is generated by cql-gen together with the conditions. 
These methods allow us to differentiate null objects from objects not loaded from the database, 
since when trying to browse a relation that was not loaded we will get `cql.ErrRelationNotLoaded`. 

In this tutorial we have used preloading and relation getters, 
for more details you can read :doc:`/cql/preloading`.

Tutorial 7: dynamic operators
-------------------------------

So far we have performed operations that take as input a static value (equal to "Paris" or greater than 1000000) 
but what if now we would like to differentiate these two Paris from each other based on whether they 
are the capital of their country.

In the tutorial_7.go file you will find that we can perform this query as follows:

.. code-block:: go
    :emphasize-lines: 6

    cities, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
        conditions.City.Country(
            conditions.Country.CapitalID.IsDynamic().Eq(conditions.City.ID),
        ),
    ).Find()

We can run this tutorial with `make tutorial_7` and we will obtain the following result:

.. code-block:: none

    Cities named 'Paris' that are the capital of their country are:
        1: City{ID: 3, Name: Paris, Population: 2161000, CountryID:2, Country:<nil> }

As you can see, again we only get the Paris in France.

In this tutorial we have used dynamic conditions, 
for more details you can read :ref:`cql/advanced_query:Dynamic operators`.

Tutorial 8: update
-------------------------------

So far we have only made select queries, but in this tutorial we want to edit the population of Paris.

In the tutorial_8.go file you will find that we can perform this query as follows:

.. code-block:: go

    updated, err := cql.Update[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Paris")),
        conditions.City.Country(
            conditions.Country.Name.Is().Eq(cql.String("France")),
        ),
    ).Returning(&cities).Set(
        conditions.City.Population.Set().Eq(cql.Int64(2102650)),
    )

We can run this tutorial with `make tutorial_8` and we will obtain the following result:

.. code-block:: none

    Updated 1 city: City{ID: 3, Name: Paris, Population: 2102650, CountryID:2, Country:<nil> }
    Initial population was 2161000

As you can see, first we can know the number of updated models with the value "updated" returned by the Set method 
(according to the number of models that meet the conditions entered in the Update method). 
On the other hand, it is also possible to obtain the information of the updated models using the Returning method.

In this tutorial we have used updates, 
for more details you can read :doc:`/cql/update`.

Tutorial 9: create and delete
-------------------------------

In this tutorial we want to create a new city called Rennes and then delete it.

In the tutorial_9.go file you will find that we can perform this query as follows:

.. code-block:: go
    :caption: Create

    rennes := models.City{
        CountryID:  france.ID,
        Name:       "Rennes",
        Population: 215366,
    }

    inserted, err := cql.Insert(context.Background(), db, &rennes).Exec()

.. code-block:: go
    :caption: Delete

    deleted, err := cql.Delete[models.City](
        context.Background(),
        db,
        conditions.City.Name.Is().Eq(cql.String("Rennes")),
    ).Exec()

We can run this tutorial with `make tutorial_9` and we will obtain the following result:

.. code-block:: none
    Inserted 1 city
    Deleted 1 city

Here, we simply get the number of inserted and deleted models through the variable returned by the Exec method
(according to the number of models that meet the conditions entered in the Insert/Delete method).

In this tutorial we have used create and delete, 
for more details you can read :doc:`/cql/create` and :doc:`/cql/delete`.

Tutorial 10: Collections
-------------------------------

In this tutorial we want to obtain all the countries that have a city called 'Paris'

In the tutorial_10.go file you will find that we can perform a query as follows:

.. code-block:: go

    countries, err := cql.Query[models.Country](
        context.Background(),
        db,
        conditions.Country.Cities.Any(
            conditions.City.Name.Is().Eq(cql.String("Paris")),
        ),
    ).Find()

We can run this tutorial with `make tutorial_10` and we will obtain the following result:

.. code-block:: none

    Countries that have a city called 'Paris' are:
        1: Country{ID: 1, Name: United States of America, CapitalID:2, Capital:<nil> }
        2: Country{ID: 2, Name: France, CapitalID:3, Capital:<nil> }

As you can see, again we only get the Paris in France.

In this tutorial we have used conditions over collections, 
for more details you can read :ref:`cql/advanced_query:Collections`.

Tutorial 11: Compile type safety
-----------------------------------

In this tutorial we want to verify that cql is compile-time safe.

In the tutorial_11.go file you will find that we try to perform a query as follows:

.. code-block:: go

    _, err := cql.Query[models.City](
        context.Background(),
        db,
        conditions.Country.Name.Is().Eq(cql.String("Paris")),
    ).Find()

We can run this tutorial with `make tutorial_11` and we will obtain the following error during compilation:

.. code-block:: none

    ./tutorial_11.go:20:3:
        cannot use conditions.Country.Name.Is().Eq(cql.String("Paris"))
        (value of interface type condition.WhereCondition[models.Country]) as condition.Condition[models.City]...

As you can see, in this tutorial we are trying to put a condition on Country 
(conditions.Country) to a Query whose main model is City (Query[models.City]). 
This would be equivalent to trying to execute the following SQL query:

.. code-block:: SQL

    SELECT * FROM cities
    WHERE countries.name = "Paris"

Therefore, we will get a compilation error and this incorrect code will never be executed.

For more details you can read :doc:`/cql/type_safety`.
