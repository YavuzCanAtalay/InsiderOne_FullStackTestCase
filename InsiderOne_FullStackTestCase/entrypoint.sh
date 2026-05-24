#!/bin/sh
set -e

if [ -n "$DATABASE_URL" ]; then
    psql "$DATABASE_URL" -f ./migrations/001_create_tables.sql
    psql "$DATABASE_URL" -f ./migrations/002_seed_data.sql
fi

exec ./football-api
