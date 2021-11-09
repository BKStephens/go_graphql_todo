# go_graphql_todo

## Installation

Install tern:
```
cd ~ && go get -u github.com/jackc/tern
``

Create databases:
```
createdb go_graphql_todo_dev
createdb go_graphql_todo_test
```

Migrate databases:
```
APP_ENV="dev" ~/go/bin/tern migrate -c server/db/tern.conf -m server/db/migrations
APP_ENV="test" ~/go/bin/tern migrate -c server/db/tern.conf -m server/db/migrations
```

## Commands

Rollback migration:
```
~/go/bin/tern migrate -d=-1 -c server/db/tern.conf -m server/db/migrations
```

## Run tests

```
go test -p 1 ./... -count=1
```

## Run server

```
JWT_SECRET_KEY=secret go run ./server
```
