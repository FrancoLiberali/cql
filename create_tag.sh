#!/bin/bash

git tag v$1
git push origin v$1
git tag cql-gen/v$1
git push origin cql-gen/v$1
git tag cqllint/v$1
git push origin cqllint/v$1