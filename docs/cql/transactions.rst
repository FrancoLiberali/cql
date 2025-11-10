==============================
Transactions
==============================

The cql methods Query, Update, Delete, Insert have as their second parameter:

.. code-block:: go
    tx *cql.DB

This, in addition to logically receiving the database where you want to run the query as 
seen in the previous sections, runs multiple queries within the same transaction.

For this, the cql.DB class provides the method Transaction. 
The function passed by parameter will be executed inside a gorm transaction 
(for more information visit https://gorm.io/docs/transactions.html). 
Using this method will also allow the transaction execution time to be logged.

For example, the following transaction would be committed
(if no errors occur during the execution of the queries):

.. code-block:: go

    var db *cql.DB
    ctx := context.Background()

    db.Transaction(ctx, func(tx *cql.DB) error {
        myModels, err := cql.Query[MyModel](
            ctx,
            tx,
            conditions.MyModel.Code.Is().Eq(cql.Int64(4)),
        ).Find()
        if err != nil {
            return err
        }

        updated, err := cql.Update[MyModel](
            ctx,
            tx,
            conditions.MyModel.Code.Is().Eq(cql.Int64(4)),
        ).Set(
            conditions.MyModel.Code.Set().Eq(cql.Int(2)),
        )
        if err != nil {
            return err
        }
    })

But the next one will be rolled back:

.. code-block:: go

    var db *cql.DB
    ctx := context.Background()

    db.Transaction(ctx, func(tx *cql.DB) error {
        updated, err := cql.Update[MyModel](
            ctx,
            tx,
            conditions.MyModel.Code.Is().Eq(cql.Int64(4)),
        ).Set(
            conditions.MyModel.Code.Set().Eq(cql.Int(2)),
        )
        if err != nil {
            return err
        }

        return errors.New("example error to rollback")
    })

