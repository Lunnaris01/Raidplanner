#!/bin/bash

if [ -f ../.env ]; then

    set -a
    source ../.env
    set +a
fi


cd sql/schema
goose turso "${TURSO_DATABASE_URL}?authToken=${TURSO_AUTH_TOKEN}" up