# Golang Boilerplate

This is a source boilerplate for Golang using Gin framework help you save your time to start with new project.

## Installation

```bash
go mod init
```

## Environment variables

```bash
touch .env
```

Paste these to your .env file

```
PORT=3000
DB_URL=postgres://username:password@host:port/database_name?sslmode=disable
```

## Run

Install make before run your app

```bash
make run
```

Then your app will be run

```bash
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/v1/                  --> github.com/duyanhitbe/go-boilerplate/internal/routes.(*Server).registerIndexRoutes.func1 (3 handlers)
[GIN-debug] GET    /api/v1/todo/             --> github.com/duyanhitbe/go-boilerplate/internal/handlers.(*TodoHandler).GetAllTodo-fm (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :3000
```

## Generation

### Migration

Create a new migration by this command

```bash
make migrate name=migration_name
```

Then two files was created in `internal/database/migrations`<br/>  <br />
Write some sql, and run:

```bash
make migrate-up
```

If you want to revert what you did, run:

```bash
make migrate-down
```

When your migration got dirty, run:

```bash
make migrate-force step=step_number
```

### Database

Create a `sql` file in `internal/database/sql`, write some operations. Then run:

```bash
make sqlc
```
Some `go` files will be generated in `internal/database/generated` <br />  <br />
You can use those functions to get data from database
