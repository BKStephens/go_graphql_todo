# go_graphql_todo

## Installation

Install tern:
```
cd ~ && go get -u github.com/jackc/tern
``

Create database:
```
createdb go_graphql_todo_dev
```

Migrate database:
```
~/go/bin/tern migrate -c server/db/tern.conf -m server/db/migrations
```

## Commands

Rollback migration:
```
~/go/bin/tern migrate -d=-1 -c server/db/tern.conf -m server/db/migrations
```
