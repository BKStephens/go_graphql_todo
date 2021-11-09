package service

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/bkstephens/go_graphql_todo/server/types"
	"github.com/jackc/pgx/v4/pgxpool"
)

func UpdateTodo(params types.UpdateTodoParams, pool *pgxpool.Pool) types.Todo {
	var updatedTodoId int
	var updatedText string
	var updatedDone bool
	var updatedListId int
	args := []interface{}{}
	allowedParams := parseParams(params)
	if len(allowedParams) == 0 {
		return types.Todo{}
	}
	var updateFragment string
	for i := 0; i < len(allowedParams); i++ {
		value := allowedParams[i]["value"]
		dbField := allowedParams[i]["dbField"]
		args = append(args, value)
		updateFragment += dbField + " = $" + strconv.Itoa(len(args))
		if i+1 < len(allowedParams) {
			updateFragment += ", "
		}
	}

	args = append(args, *params.UserId, strconv.Itoa(*params.ID))
	err := pool.QueryRow(
		context.Background(),
		`UPDATE todos SET `+updateFragment+
			` FROM (SELECT id AS listId FROM todo_lists WHERE user_id = $`+strconv.Itoa(len(args)-1)+`) subquery`+
			` WHERE subquery.listId = todo_list_id AND id = $`+strconv.Itoa(len(args))+
			` RETURNING id, text, done, todo_list_id;`,
		args...,
	).Scan(&updatedTodoId, &updatedText, &updatedDone, &updatedListId)
	if err != nil {
		fmt.Println(err)
		return types.Todo{}
	}

	return types.Todo{ID: updatedTodoId, Text: updatedText, Done: updatedDone, ListId: updatedListId}
}

func parseParams(params types.UpdateTodoParams) []map[string]string {
	values := reflect.ValueOf(params)
	allowedParams := []map[string]string{}
	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		if !field.IsValid() || field.IsNil() {
			continue
		}
		var value string
		switch reflect.TypeOf(field.Elem().Interface()).String() {
		case "int":
			value = strconv.Itoa(field.Elem().Interface().(int))
		case "bool":
			value = strconv.FormatBool(field.Elem().Interface().(bool))
		default:
			value = field.Elem().Interface().(string)
		}

		dbField := values.Type().Field(i).Tag.Get("db")
		if dbField == "done" || dbField == "text" {
			allowedParams = append(allowedParams, map[string]string{"dbField": dbField, "value": value})
		}
	}
	return allowedParams
}
