package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/bkstephens/go_graphql_todo/server/service"
	"github.com/bkstephens/go_graphql_todo/server/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	serverUrl string
}

var envVariables = map[string]string{
	"JWT_SECRET_KEY": "secret",
	"PORT":           "8081",
	"DATABASE_URL":   "postgres://postgres:postgres@localhost:5432/go_graphql_todo_test",
}

func (suite *testSuite) SetupSuite() {
	for k, v := range envVariables {
		os.Setenv(k, v)
	}
	suite.serverUrl = "http://localhost:" + envVariables["PORT"]
	go Initialize()
	for {
		req, _ := http.NewRequest("GET", suite.serverUrl+"/health", nil)
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == 200 {
			break
		}
	}
}

func (suite *testSuite) TearDownSuite() {
	for k := range envVariables {
		os.Unsetenv(k)
	}
	DBPool.Close()
}

func (suite *testSuite) TearDownTest() {
	if _, err := DBPool.Exec(context.Background(), "TRUNCATE users CASCADE;"); err != nil {
		suite.T().Fatal(err)
		return
	}
}

func (suite *testSuite) TestPostSignupWithoutRequiredFields() {
	reader := strings.NewReader("")
	req, err := http.NewRequest("POST", suite.serverUrl+"/signup", reader)
	if err != nil {
		suite.T().Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), 400, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	assert.Equal(suite.T(), "username, email, and password must be provided", bodyString)
}

func (suite *testSuite) TestPostSignupWithRequiredFields() {
	reqBody := []byte(`{"username": "testuser", "email": "testuser@example.com", "password": "password"}`)
	req, err := http.NewRequest("POST", suite.serverUrl+"/signup", bytes.NewBuffer(reqBody))
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	assert.Equal(suite.T(), bodyString, "Signup successful")
}

func (suite *testSuite) TestPostSignupForUserThatExits() {
	service.InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, DBPool)

	reqBody := []byte(`{"username": "testuser", "email": "testuser@example.com", "password": "password"}`)
	req, err := http.NewRequest("POST", suite.serverUrl+"/signup", bytes.NewBuffer(reqBody))
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		suite.T().Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(suite.T(), 500, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)
	assert.Equal(suite.T(), bodyString, "Signup failed")
}

func (suite *testSuite) TestPostLoginForUserThatExists() {
	service.InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, DBPool)

	reqBody := []byte(`{"username": "testuser", "password": "password"}`)
	req, err := http.NewRequest("POST", suite.serverUrl+"/login", bytes.NewBuffer(reqBody))
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	type responseT struct {
		Token string `json:"token"`
	}
	var response responseT
	json.Unmarshal(bodyBytes, &response)
	assert.NotNil(suite.T(), response.Token)
}

func (suite *testSuite) TestPostLoginForUserThatExistsButWrongPassword() {
	service.InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, DBPool)

	reqBody := []byte(`{"username": "testuser", "password": "letmein"}`)
	req, err := http.NewRequest("POST", suite.serverUrl+"/login", bytes.NewBuffer(reqBody))
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), 401, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	assert.Equal(suite.T(), "Unauthorized", bodyString)
}

func (suite *testSuite) TestPostLoginForUserThatDoesntExist() {
	reqBody := []byte(`{"username": "testuser", "password": "password"}`)
	req, err := http.NewRequest("POST", suite.serverUrl+"/login", bytes.NewBuffer(reqBody))
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), 500, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	assert.Equal(suite.T(), "Could not find user", bodyString)
}

func (suite *testSuite) TestAuthorizationLoggedIn() {
	service.InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, DBPool)
	reqBody := []byte(`{"username": "testuser", "password": "password"}`)
	req, err := http.NewRequest("POST", suite.serverUrl+"/login", bytes.NewBuffer(reqBody))
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	type responseT struct {
		Token string `json:"token"`
	}
	var response responseT
	json.Unmarshal(bodyBytes, &response)

	req, err = http.NewRequest("GET", suite.serverUrl+"/api/v1/authorized", nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+response.Token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(suite.T(), 200, resp.StatusCode)
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	assert.Equal(suite.T(), "Reached authenticated endpoint", bodyString)
}

func (suite *testSuite) TestAuthorizationNoAuthorizatonToken() {
	req, err := http.NewRequest("GET", suite.serverUrl+"/api/v1/authorized", nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(suite.T(), 401, resp.StatusCode)
}

func (suite *testSuite) TestTodoLifecycle() {
	// createTodo
	service.InsertUser(types.InsertUserParams{Username: "testuser", Email: "testuser@example.com", Password: "password"}, DBPool)
	reqBody := []byte(`{"username": "testuser", "password": "password"}`)
	req, err := http.NewRequest("POST", suite.serverUrl+"/login", bytes.NewBuffer(reqBody))
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	type responseT struct {
		Token string `json:"token"`
	}
	var response responseT
	json.Unmarshal(bodyBytes, &response)
	resp.Body.Close()

	req, err = http.NewRequest("GET", suite.serverUrl+`/graphql?query=mutation+_{createTodo(text:"My+new+todo"){id,text,done}}`, nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+response.Token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}

	assert.Equal(suite.T(), 200, resp.StatusCode)
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Regexp(suite.T(), `{"data":{"createTodo":{"done":false,"id":[0-9]+,"text":"My new todo"}}}`, bodyString)
	resp.Body.Close()

	// todoList
	req, err = http.NewRequest("GET", suite.serverUrl+`/graphql?query={todoList{id,text,done}}`, nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+response.Token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}

	assert.Equal(suite.T(), 200, resp.StatusCode)
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString = string(bodyBytes)

	assert.Regexp(suite.T(), `{"data":{"todoList":[{"done":false,"id":[0-9]+,"text":"My new todo"}]}}`, bodyString)
	resp.Body.Close()

	var todoId = regexp.MustCompile(`"id":(?P<id>[0-9]+),`).FindStringSubmatch(bodyString)[1]
	// updateTodo
	//req, err = http.NewRequest("GET", suite.serverUrl+`/graphql?query=mutation+_{createTodo(text:"My+new+todo"){id,text,done}}`, nil)
	req, err = http.NewRequest("GET", suite.serverUrl+`/graphql?query=mutation+_{updateTodo(id:`+todoId+`,text:"My+new+todo+updated"){id,text,done}}`, nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+response.Token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		suite.T().Fatal(err)
	}

	assert.Equal(suite.T(), 200, resp.StatusCode)
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString = string(bodyBytes)

	assert.Regexp(suite.T(), `{"data":{"updateTodo":{"done":false,"id":[0-9]+,"text":"My new todo updated"}}}`, bodyString)
	resp.Body.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
