package service

import (
	"context"
	"reflect"
	"strconv"

	"github.com/bkstephens/go_graphql_todo/server/types"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetTodos(params types.GetTodosParams, pool *pgxpool.Pool) []types.Todo {
	where := ""
	args := []interface{}{}
	args = append(args, strconv.Itoa(*params.UserId))
	values := reflect.ValueOf(params)
	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		if !field.IsValid() || field.IsNil() {
			continue
		}
		dbField := values.Type().Field(i).Tag.Get("db")
		if dbField == "" {
			continue
		}

		v := field.Elem().Interface()
		args = append(args, v)
		if i+1 < values.NumField() {
			where += " AND "
		}
		where += dbField + " = $" + strconv.Itoa(len(args))
	}

	var todos []types.Todo
	rows, err := pool.Query(
		context.Background(),
		`SELECT id, text, done, todo_list_id
		 FROM todos
		 WHERE todo_list_id IN (SELECT id FROM todo_lists WHERE user_id = $1) `+where,
		args...,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var text string
		var listId int
		var done bool
		rows.Scan(&id, &text, &done, &listId)
		todos = append(todos, types.Todo{ID: id, Text: text, Done: done, ListId: listId})
	}

	return todos
}
