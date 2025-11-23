#!/bin/bash
set -e

echo "---- Running migrations ----"

until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q'; do
  >&2 echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

for file in /docker-entrypoint-initdb.d/*.up.sql; do
    echo "Running migration: $(basename $file)"
    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "$file"
done

echo "Migrations completed successfully"