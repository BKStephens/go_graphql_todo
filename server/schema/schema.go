package schema

import (
	"github.com/bkstephens/go_graphql_todo/server/service"
	"github.com/bkstephens/go_graphql_todo/server/types"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/jackc/pgx/v4/pgxpool"
)

var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"text": &graphql.Field{
			Type: graphql.String,
		},
		"done": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

func getRootMutation(c *gin.Context) *graphql.Object {
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createTodo": &graphql.Field{
				Type:        todoType,
				Description: "Create new todo",
				Args: graphql.FieldConfigArgument{
					"text": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"listId": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					text, _ := params.Args["text"].(string)
					listId, _ := params.Args["listId"].(int)
					pool := c.MustGet("pool").(*pgxpool.Pool)
					userId := int(c.MustGet("userId").(float64))
					newTodo := service.InsertTodo(types.InsertTodoParams{Text: text, ListId: listId, UserId: userId}, pool)

					return newTodo, nil
				},
			},
			"updateTodo": &graphql.Field{
				Type:        todoType,
				Description: "Update existing todo, mark it done or not done",
				Args: graphql.FieldConfigArgument{
					"done": &graphql.ArgumentConfig{
						Type: graphql.Boolean,
					},
					"text": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					pool := c.MustGet("pool").(*pgxpool.Pool)
					id, _ := params.Args["id"].(int)
					userId := int(c.MustGet("userId").(float64))
					updateParams := types.UpdateTodoParams{ID: &id, UserId: &userId}
					if params.Args["done"] != nil {
						done := params.Args["done"].(bool)
						updateParams.Done = &done
					}
					if params.Args["text"] != nil {
						text := params.Args["text"].(string)
						updateParams.Text = &text
					}

					updatedTodo := service.UpdateTodo(updateParams, pool)

					return updatedTodo, nil
				},
			},
		},
	})
	return rootMutation
}

func getRootQuery(c *gin.Context) *graphql.Object {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{

			"todo": &graphql.Field{
				Type:        todoType,
				Description: "Get single todo",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					idQuery, _ := params.Args["id"].(int)
					userId := int(c.MustGet("userId").(float64))
					pool := c.MustGet("pool").(*pgxpool.Pool)
					todos := service.GetTodos(types.GetTodosParams{UserId: &userId, TodoId: &idQuery}, pool)
					return todos[0], nil
				},
			},

			"todoList": &graphql.Field{
				Type:        graphql.NewList(todoType),
				Description: "List of todos",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userId := int(c.MustGet("userId").(float64))
					pool := c.MustGet("pool").(*pgxpool.Pool)
					todos := service.GetTodos(types.GetTodosParams{UserId: &userId}, pool)
					return todos, nil
				},
			},
		},
	})
	return rootQuery
}

func GetTodoSchema(c *gin.Context) graphql.Schema {
	var todoSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    getRootQuery(c),
		Mutation: getRootMutation(c),
	})

	return todoSchema
}
