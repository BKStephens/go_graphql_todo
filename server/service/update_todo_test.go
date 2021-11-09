package service

import (
	"context"
	"testing"

	"github.com/bkstephens/go_graphql_todo/server/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuiteUpdateTodo struct {
	suite.Suite
	pool *pgxpool.Pool
}

func (suite *testSuiteUpdateTodo) SetupSuite() {
	var err error
	suite.pool, err = pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/go_graphql_todo_test")
	if err != nil {
		panic(err)
	}
}

func (suite *testSuiteUpdateTodo) TearDownTest() {
	if _, err := suite.pool.Exec(context.Background(), "TRUNCATE users CASCADE;"); err != nil {
		suite.T().Fatal(err)
		return
	}
}

func (suite *testSuiteUpdateTodo) TearDownSuite() {
	suite.pool.Close()
}

func (suite *testSuiteUpdateTodo) TestUpdateTodo() {
	userId, _ := InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, suite.pool)

	todo := InsertTodo(types.InsertTodoParams{UserId: userId}, suite.pool)
	done := true
	updatedTodo := UpdateTodo(types.UpdateTodoParams{UserId: &userId, ID: &todo.ID, Done: &done}, suite.pool)
	assert.Equal(suite.T(), updatedTodo.Done, done)
}

func (suite *testSuiteUpdateTodo) TestCantUpdateOthersTodo() {
	userId1, _ := InsertUser(types.InsertUserParams{Username: "testuser1", Email: "testuser1@example.com", Password: "password"}, suite.pool)
	userId2, _ := InsertUser(types.InsertUserParams{Username: "testuser2", Email: "testuser2@example.com", Password: "password"}, suite.pool)

	todo := InsertTodo(types.InsertTodoParams{UserId: userId1}, suite.pool)

	done := true
	updatedTodo := UpdateTodo(types.UpdateTodoParams{UserId: &userId2, ID: &todo.ID, Done: &done}, suite.pool)
	assert.Equal(suite.T(), updatedTodo, types.Todo{})
}

func TestSuiteUpdateTodo(t *testing.T) {
	suite.Run(t, new(testSuiteUpdateTodo))
}
