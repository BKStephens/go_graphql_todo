# go_graphql_todo

## Installation

Install tern:
```
cd ~ && go get -u github.com/jackc/tern
```

Create databases:
```
createdb go_graphql_todo_dev
createdb go_graphql_todo_test
```

Migrate databases:
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/go_graphql_todo_dev go run server/cmd/migrate/migrate.go
DATABASE_URL=postgres://postgres:postgres@localhost:5432/go_graphql_todo_test go run server/cmd/migrate/migrate.go
```

## Commands

Rollback migration:
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/go_graphql_todo_dev go run server/cmd/migrate_rollback/migrate_rollback.go
```

## Run tests

```
go test -p 1 ./... -count=1
```

## Run server

```
cd client && yarn install && yarn build && cd .. && JWT_SECRET_KEY=secret go run ./server
```

Now you can navigate to [http://localhost:8080](http://locahost:8080) and try
out the app. If you make changes to the frontend code then run `yarn build` and
refresh the page.

### Optional React hot reloading

```
cd client && yarn start
```

Now you can navigate to [http://localhost:3000](http://localhost:3000) and start
using the app.
