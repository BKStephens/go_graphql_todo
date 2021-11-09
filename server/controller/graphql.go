package controller

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/bkstephens/go_graphql_todo/server/schema"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func GraphQLHandler(c *gin.Context) {
	b := new(bytes.Buffer)
	query := c.Request.URL.Query()
	for _, value := range query {
		for _, value2 := range value {
			fmt.Fprintf(b, value2)
		}
	}
	result := executeQuery(b.String(), schema.GetTodoSchema(c))
	c.JSON(http.StatusOK, result)
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
