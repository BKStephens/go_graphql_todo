package service

import (
	"context"

	"github.com/bkstephens/go_graphql_todo/server/types"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InsertTodo(params types.InsertTodoParams, pool *pgxpool.Pool) types.Todo {
	var listId int
	if params.ListId != 0 {
		listId = params.ListId
	} else {
		pool.QueryRow(context.Background(), "SELECT id FROM todo_lists WHERE user_id = $1 ORDER BY created_at;", params.UserId).Scan(&listId)
		if listId == 0 {
			listId = InsertTodoList(pool, params.UserId)
		}
	}
	var todoId int
	var insertedText string
	var insertedDone bool
	var insertedListId int
	err := pool.QueryRow(context.Background(), "INSERT INTO todos(text, done, todo_list_id) VALUES ($1, $2, $3) RETURNING id, text, done, todo_list_id;", params.Text, params.Done, listId).Scan(&todoId, &insertedText, &insertedDone, &insertedListId)
	if err != nil {
		panic(err)
	}

	return types.Todo{ID: todoId, Text: insertedText, Done: insertedDone, ListId: insertedListId}
}
