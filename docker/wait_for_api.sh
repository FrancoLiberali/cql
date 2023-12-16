#!/bin/sh

until $(curl --output /dev/null --silent --fail http://localhost:$1); do
    printf '.'
    sleep 5
done