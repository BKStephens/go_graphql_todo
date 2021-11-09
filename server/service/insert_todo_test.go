package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/bkstephens/go_graphql_todo/server/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	pool *pgxpool.Pool
}

func (suite *testSuite) SetupSuite() {
	var err error
	suite.pool, err = pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/go_graphql_todo_test")
	if err != nil {
		panic(err)
	}
}

func (suite *testSuite) TearDownTest() {
	if _, err := suite.pool.Exec(context.Background(), "TRUNCATE users CASCADE;"); err != nil {
		suite.T().Fatal(err)
		return
	}
}

func (suite *testSuite) TearDownSuite() {
	suite.pool.Close()
}

func (suite *testSuite) TestInsertTodoCreatesTodoListWhenUserHasNone() {
	var beforeCount int
	suite.pool.QueryRow(context.Background(), `SELECT COUNT(1) FROM todo_lists;`).Scan(&beforeCount)

	userId, err := InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, suite.pool)
	if err != nil {
		fmt.Println(err)
	}

	todo := InsertTodo(types.InsertTodoParams{UserId: userId}, suite.pool)

	var afterCount int
	suite.pool.QueryRow(context.Background(), `SELECT COUNT(1) FROM todo_lists;`).Scan(&afterCount)
	assert.NotNil(suite.T(), todo.ID)
	assert.Equal(suite.T(), beforeCount+1, afterCount)
}

func (suite *testSuite) TestInsertTodoUsesFirstTodoListIfNoneSpecified() {
	userId1, _ := InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, suite.pool)
	userId2, _ := InsertUser(types.InsertUserParams{Username: "testuser2", Email: "testuser2@example.com", Password: "password"}, suite.pool)

	todoList1 := InsertTodoList(suite.pool, userId1)
	InsertTodoList(suite.pool, userId1)
	InsertTodoList(suite.pool, userId2)
	InsertTodoList(suite.pool, userId2)

	todo := InsertTodo(types.InsertTodoParams{UserId: userId1}, suite.pool)

	var selectedTodoList int
	suite.pool.QueryRow(context.Background(), `SELECT todo_list_id FROM todos WHERE id = $1;`, todo.ID).Scan(&selectedTodoList)
	assert.Equal(suite.T(), selectedTodoList, todoList1)
}

func (suite *testSuite) TestInsertTodoUsesSpecifiedTodoList() {
	userId1, _ := InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, suite.pool)
	userId2, _ := InsertUser(types.InsertUserParams{Username: "testuser2", Email: "testuser2@example.com", Password: "password"}, suite.pool)

	InsertTodoList(suite.pool, userId1)
	todoList2 := InsertTodoList(suite.pool, userId1)
	InsertTodoList(suite.pool, userId2)
	InsertTodoList(suite.pool, userId2)

	todo := InsertTodo(types.InsertTodoParams{UserId: userId1, ListId: todoList2}, suite.pool)

	var selectedTodoList int
	suite.pool.QueryRow(context.Background(), `SELECT todo_list_id FROM todos WHERE id = $1;`, todo.ID).Scan(&selectedTodoList)
	assert.Equal(suite.T(), selectedTodoList, todoList2)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
