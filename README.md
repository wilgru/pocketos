# PocketSystem

PocketSystem is an adaptive productivity ecosystem built around PocketBase.

This repository is currently set up as a backend-only PocketBase boilerplate with extension points for both Go and JavaScript.

## Requirements

- Go 1.27+

## Getting started

1. Install dependencies:

   ```bash
   go mod tidy
   ```

2. Start the PocketBase server:

   ```bash
   go run . serve
   ```

PocketBase stores its SQLite database and runtime data in `pb_data/`.

## Project structure

- `main.go` boots PocketBase and enables JavaScript hooks and JS migrations.
- `internal/extensions/` is where Go routes, hooks, and other backend extensions live.
- `pb_hooks/` contains JavaScript hook files loaded by the PocketBase JS VM plugin.
- `pb_migrations/` is reserved for JavaScript migrations created with the PocketBase migrate command.

## Included extension examples

- Go extension route: `GET /api/healthz/go`
- JavaScript hook route: `GET /api/healthz/js`
- Go collection registry: the `_tools` system collection is created by a JavaScript migration; each new base collection is added with its name in the required, unique `targetCollectionName` field.

These routes return simple JSON responses so you can verify both extension systems are wired correctly.

## Extending the project

### Add Go extensions

Add new PocketBase hooks, routes, or services in `internal/extensions/` and register them from `Register`.

### Add JavaScript extensions

Create new `*.pb.js` files in `pb_hooks/`. PocketBase will auto-load them on startup, and changes are watched automatically while developing.

### Add JavaScript migrations

Use PocketBase's migrate command to create migration files in `pb_migrations/`.
