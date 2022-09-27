# Golang-ninja-project
## Work in progress 

### Todo:
- Upload to external storage + save info to Postgres(/upload endpoint)
- Get All Uploaded info from Postgres(/files endpoint)
- Docker
- Ext file Storage
- Tests

### ADV:
- Linter
- Logging
- Google Doc with Design
- CI/CD

### Tools:
- go 1.19
- postgres

### How to use
Run container with postgres:
```cmd
docker run -d --name=file-db -e POSTGRES_PASSWORD=postgres -v ${HOME}/pgdata/:/var/lib/postgresql/data -p 5432:5432 --rm postgres
```
Migrate tables using migrate tool(https://github.com/golang-migrate/migrate):
```cmd
migrate -path ./schema -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
```
