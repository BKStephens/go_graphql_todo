package service

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InsertTodoList(pool *pgxpool.Pool, userId int) int {
	var todoListId int
	pool.QueryRow(context.Background(), "INSERT INTO todo_lists(user_id) VALUES ($1) RETURNING id;", userId).Scan(&todoListId)
	return todoListId
}
