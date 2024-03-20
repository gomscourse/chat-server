#!/bin/bash
source /root/.env

# wait for the postgres to be ready
while ! nc -z postgres 5432; do
  >&2 echo "postgres is unavailable, wait..."
  sleep 2
done

# run migrations
export MIGRATION_DSN="host=postgres port=5432 dbname=$PG_DATABASE_NAME user=$PG_USER password=$PG_PASSWORD sslmode=disable"

# Run goose and check if the command was successful
if ! goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v; then
    echo "Migration failed"
    exit 1  # Exit with a non-zero status to indicate error
fi

# run server
./chat_server

