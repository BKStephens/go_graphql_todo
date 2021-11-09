package service

import (
	"context"
	"testing"

	"github.com/bkstephens/go_graphql_todo/server/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuiteGetTodos struct {
	suite.Suite
	pool *pgxpool.Pool
}

func (suite *testSuiteGetTodos) SetupSuite() {
	var err error
	suite.pool, err = pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/go_graphql_todo_test")
	if err != nil {
		panic(err)
	}
}

func (suite *testSuiteGetTodos) TearDownTest() {
	if _, err := suite.pool.Exec(context.Background(), "TRUNCATE users CASCADE;"); err != nil {
		suite.T().Fatal(err)
		return
	}
}

func (suite *testSuiteGetTodos) TearDownSuite() {
	suite.pool.Close()
}

func (suite *testSuiteGetTodos) TestGetTodosReturnsAllOfAUsersTodos() {
	userId1, _ := InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, suite.pool)
	userId2, _ := InsertUser(types.InsertUserParams{Username: "testuser2", Email: "testuser2@example.com", Password: "password"}, suite.pool)
	todo1 := InsertTodo(types.InsertTodoParams{UserId: userId1}, suite.pool)
	todo2 := InsertTodo(types.InsertTodoParams{UserId: userId1}, suite.pool)
	InsertTodo(types.InsertTodoParams{UserId: userId2}, suite.pool)

	todos := GetTodos(types.GetTodosParams{UserId: &userId1}, suite.pool)

	assert.Equal(suite.T(), len(todos), 2)
	assert.Equal(suite.T(), todo1.ID, todos[0].ID)
	assert.Equal(suite.T(), todo2.ID, todos[1].ID)
}

func (suite *testSuiteGetTodos) TestGetTodosReturnsOneOfAUsersTodos() {
	userId, _ := InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, suite.pool)
	todo1 := InsertTodo(types.InsertTodoParams{UserId: userId}, suite.pool)
	InsertTodo(types.InsertTodoParams{UserId: userId}, suite.pool)

	todos := GetTodos(types.GetTodosParams{UserId: &userId, TodoId: &todo1.ID}, suite.pool)

	assert.Equal(suite.T(), len(todos), 1)
	assert.Equal(suite.T(), todo1.ID, todos[0].ID)
}

func (suite *testSuiteGetTodos) TestGetTodosReturnsTodosForOneList() {
	userId, _ := InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, suite.pool)
	todo1 := InsertTodo(types.InsertTodoParams{UserId: userId}, suite.pool)
	InsertTodo(types.InsertTodoParams{UserId: userId}, suite.pool)

	todos := GetTodos(types.GetTodosParams{UserId: &userId, TodoId: &todo1.ID}, suite.pool)

	assert.Equal(suite.T(), len(todos), 1)
	assert.Equal(suite.T(), todo1.ID, todos[0].ID)
}

func TestSuiteGetTodos(t *testing.T) {
	suite.Run(t, new(testSuiteGetTodos))
}
