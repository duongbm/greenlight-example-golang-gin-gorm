# Project structure

- `bin`: contains compiled application binaries that ready for deployment to production
- `cmd/api`: contain primary logic code
- `internal`: contain various packages used by api such as interacting with db, validation, job, task,...
- `migrations`: contain sql migration files
- `remote`: contain configuration files and setup scripts
- `Makefile`: contain recipes for automating common administrative task like: auditing, building binary, execute
  migrations

* Note: any packages that lives under `internal` can only import by code inside the parent of the this directory.

# Features

- [x] Encapsulate Routing
- [x] Custom encode/decode JSON
- [x] Manage HTTP response
- [x] Validating JSON input
- [x] Setup GORM, optimize connection pool
- [x] SQL Migrate
- [x] Perform simple CRUD (GET, POST, PUT, PATCH, DELETE)
- [x] Optimistic Concurrency Control
- [x] Filter(Simple full text search), Sort, Pagination
- [x] Structured JSON log entries
- [x] Panic Recover Middleware
- [ ] Rate Limit
- [ ] Manage SQL query timeout
- [ ] User Model & Registration
- [ ] Sending emails
- [ ] User Activation
- [ ] Authentication
- [ ] Permissions
- [ ] Metrics
- [ ] Building, Versioning and Quality control