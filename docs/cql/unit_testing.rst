==============================
Unit testing
==============================

Inherited from gorm, cql is compatible with `go-sqlmock <github.com/DATA-DOG/go-sqlmock>` 
to configure a mock database in order to run unit tests, for example:

.. code-block:: go

    conn, mock, err := sqlmock.New()
    require.NoError(t, err)

    defer conn.Close()

    db, err := cql.Open(postgres.New(postgres.Config{
        Conn: conn,
    }))
    require.NoError(t, err)

    // set mock expectations
    mock.ExpectBegin() // transaction start
    mock.ExpectExec(`INSERT INTO "mymodels"`).
        WithArgs(sqlmock.AnyArg(), "John Doe").
        WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit() // transaction commit

    _, err = cql.Insert(context.Background(), db, &models.MyModel{Name: "John Doe"}).Exec()
    require.NoError(t, err)

    assert.NoError(t, mock.ExpectationsWereMet())