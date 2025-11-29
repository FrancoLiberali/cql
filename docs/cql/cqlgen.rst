==============================
cql-gen
==============================

`cql-gen` is the command line tool that generates the conditions to query your objects. 
For each cql Model found in the input packages, a file containing all possible Conditions 
on that object will be generated, allowing you to use cql.

Installation
----------------------------

For simply installing it, use:

.. code-block:: bash

    go install github.com/FrancoLiberali/cql/cql-gen@latest

.. warning::

    The version of cql-gen used must be the same as the version of cql. 
    You can install a specific version using `go install github.com/FrancoLiberali/cql/cql-gen@vX.Y.Z`, 
    where X.Y.Z is the version number.

Conditions generation
--------------------------------------

While conditions can be generated executing cql-gen, it's recommended to use `go generate`:

Once cql-gen is installed, inside our project we will have to create a package called conditions 
(or another name if you wish) and inside it a file with the following content:

.. code-block:: go

    package conditions

    //go:generate cql-gen ../models_path_1 ../models_path_2

where ../models_path_1 ../models_path_2 are the relative paths between the package conditions 
and the packages containing the definition of your models (can be only one).

`Example file <https://github.com/FrancoLiberali/cql-quickstart/blob/main/conditions/cql.go>`_.

Now, from the root of your project you can execute:

.. code-block:: bash

  go generate ./...

and the conditions for each of your models will be created in the conditions package.

Use of the conditions
--------------------------------------

After performing the conditions generation, 
your conditions package will have a replica of your models package, 
i.e. if, for example, the type models.MyModel is part of your models, 
the variable conditions.MyModel will be in the conditions package. 
This variable is called the condition model and it has:

- An attribute for each attribute of your original model with the same name 
  (if models.MyModel.Name exists, then conditions.MyModel.Name is generated), 
  that allows to use that attribute in your queries.
- A method for each relation of your original model with the same name 
  (if models.MyModel.MyOtherModel exists, then conditions.MyModel.MyOtherModel() is generated), 
  which will allow you to perform joins in your queries.
- Methods for :doc:`/cql/preloading`.

Then, combining these conditions you will be able to make all the queries you need in a safe way.

For details about querying, see :doc:`/cql/query`.