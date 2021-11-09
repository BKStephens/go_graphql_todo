package controller

import (
	"fmt"
	"net/http"

	"github.com/bkstephens/go_graphql_todo/server/service"
	"github.com/bkstephens/go_graphql_todo/server/types"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SignUpRequest struct {
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}

func SignupHandler(c *gin.Context) {
	signUpRequest := &SignUpRequest{}
	err := c.ShouldBindJSON(signUpRequest)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "username, email, and password must be provided")
		return
	}

	pool := c.MustGet("pool").(*pgxpool.Pool)
	err = service.InsertUser(types.InsertUserParams(*signUpRequest), pool)
	if err != nil {
		c.String(http.StatusInternalServerError, "Signup failed")
		return
	}
	c.String(http.StatusOK, "Signup successful")
}
