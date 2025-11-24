==============================
cqllint
==============================

`cqllint` is a Go linter that checks that cql queries will not generate run-time errors. 

While, in most cases, queries created using cql are checked at compile time, 
there are still some cases that can generate run-time errors (see :ref:`cql/type_safety:Runtime errors`).

cqllint analyses the Go code written to detect these cases and fix them without the need to execute the query. 
It also adds other detections that would not generate runtime errors but are possible misuses of cql.

.. note::

    At the moment, only the errors cql.ErrFieldModelNotConcerned, cql.ErrFieldIsRepeated, 
    cql.ErrAppearanceMustBeSelected and cql.ErrAppearanceOutOfRange are detected.

We recommend integrating cqllint into your CI so that the use of cql ensures 100% that your queries will be executed correctly.

Installation
----------------------------

For simply installing it, use:

.. code-block:: bash

    go install github.com/FrancoLiberali/cql/cqllint@latest

.. warning::

    The version of cqllint used must be the same as the version of cql. 
    You can install a specific version using `go install github.com/FrancoLiberali/cql/cqllint@vX.Y.Z`, 
    where X.Y.Z is the version number.

Execution
----------------------------

cqllint can be used independently by running:

.. code-block:: bash

    cqllint ./...

or using `go vet`:

.. code-block:: bash

    go vet -vettool=$(which cqllint) ./...

Errors
-------------------------------



Misuses
-------------------------

Although some cases would not generate runtime errors, cqllint will detect them as they are possible misuses of cql.

.. TODO poner un ejemplo aca de error y misuse y luego poner links a cada seccion
.. TODO poner las limitaciones de dentro de la misma funcion y eso