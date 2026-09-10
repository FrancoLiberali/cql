===========
Performance
===========

CQL generates, at compile time, a per-model *scanner* that copies each row from
the database into your Go structs **without reflection**. gorm — and tools built
on top of it such as `gorm-gen <https://gorm.io/gen/>`_ — instead scan every row
reflectively at run time. The larger the result set, the more that difference
adds up.

Benchmarks
==========

The CQL benchmarks live in ``cql/benchmarks`` and mirror gorm's own
``gorm/tests/benchmark_test.go`` shapes (same ``users`` table layout); the gorm
figures come from running that gorm suite, and the gorm-gen figures from a
generated DAO over the same table. All were produced on the same machine,
against an in-memory SQLite database, with ``go test -bench``. The numbers are
relative — what matters is CQL versus the alternatives on the *same* run, not
the absolute nanoseconds.

Reading 10,000 rows (``.Find()`` into a slice)
----------------------------------------------

This is where the generated scanner pays off:

.. list-table::
   :header-rows: 1

   * - Library
     - Time
     - Allocations
   * - **CQL**
     - **14.5 ms**
     - **119,882**
   * - gorm
     - 18.0 ms
     - 190,091
   * - gorm-gen
     - 18.3 ms
     - 190,138

CQL is about **20% faster** and does about **37% fewer allocations** than gorm,
because it never falls back to reflection while materialising the rows.

Reading a single row by primary key (``.First()``)
--------------------------------------------------

.. list-table::
   :header-rows: 1

   * - Library
     - Time
     - Allocations
   * - **CQL**
     - **7.0 µs**
     - **95**
   * - gorm
     - 7.1 µs
     - 94
   * - gorm-gen
     - 8.0 µs
     - 106

For a single row the per-query cost dominates, so CQL matches gorm and is a bit
faster than gorm-gen. The scanner's advantage grows with the number of rows
returned — the common case for list queries.

Why gorm-gen tracks gorm
========================

gorm-gen adds a type-safe query-building API, but the generated code executes
queries through gorm's runtime and therefore uses gorm's reflective scanner. Its
per-row scan cost is, by construction, the same as gorm's — the ``gorm`` and
``gorm-gen`` rows above are within noise of each other. CQL replaces exactly that
reflective step with generated code, which is what produces the difference.
