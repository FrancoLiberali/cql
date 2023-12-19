#!/bin/sh

until [ "`docker inspect -f {{.State.Health.Status}} cql-test-db`"=="healthy" ]; do
    printf '.';
    sleep 1;
done;