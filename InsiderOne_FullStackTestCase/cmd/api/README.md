# cmd/api

## Purpose
Entry point of the application.

## What goes here
- `main.go` — starts the HTTP server, loads config, wires all dependencies together (repositories, services, handlers)

## What needs to be done
- Initialize the database connection
- Register all route handlers
- Start listening on the configured port (e.g. `:8080`)
