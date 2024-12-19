#!/bin/sh
# Debugging: log environment variables
echo "Running migrations on $DATABASE_HOST:$DATABASE_PORT"
echo "Using database user: $DATABASE_USER"
echo "Using database: $DATABASE_NAME"

# Set the PostgreSQL password environment variable
export PASSWORD=$DATABASE_PASSWORD

# Wait for PostgreSQL to be available using pg_isready
until pg_isready -h $DATABASE_HOST -p 5432 -U "$DATABASE_USER"; do
  echo "Waiting for database..."
  sleep 2
done

# Run migrations
echo "Running migrations..."
psql -h $DATABASE_HOST -U "$DATABASE_USER" -d "$DATABASE_NAME" -f /app/Migration/001_create_restaurants_table.up.sql
psql -h $DATABASE_HOST -U "$DATABASE_USER" -d "$DATABASE_NAME" -f /app/Migration/002_create_items_table.up.sql

# Start the main service
echo "Starting application..."
exec "$@"