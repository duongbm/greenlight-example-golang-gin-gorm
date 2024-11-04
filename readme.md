# Project structure

- `bin`: contains compiled application binaries that ready for deployment to production
- `cmd/api`: contain primary logic code
- `internal`: contain various packages used by api such as interacting with db, validation, job, task,...
- `migrations`: contain sql migration files
- `remote`: contain configuration files and setup scripts
- `Makefile`: contain recipes for automating common administrative task like: auditing, building binary, execute
  migrations

* Note: any packages that lives under `internal` can only import by code inside the parent of the this directory.

# Advantage

- Encapsulate Routing
- Custom encode/decode JSON
- Manage HTTP response
- Validating JSON input
- Setup GORM, optimize connection pool
- SQL Migrate
- Perform CRUD (GET, POST, PUT, PATCH, DELETE)
- Optimistic Concurrency Control