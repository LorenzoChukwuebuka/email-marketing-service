#!/bin/sh
set -e

echo "Waiting for database to be ready..."
until pg_isready -h postgres-db -p 5432 -U myuser; do
  echo "Database not ready, waiting..."
  sleep 2
done

echo "Running migrations..."
migrate -database="$DB_URL" -path=/migrations up

echo "Migrations completed!"