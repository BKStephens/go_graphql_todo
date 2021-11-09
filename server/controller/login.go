package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bkstephens/go_graphql_todo/server/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"username"`
}

func LoginHandler(c *gin.Context) {
	credentials := &Credentials{}
	err := c.ShouldBindJSON(credentials)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "username and password must be provided")
		return
	}

	pool := c.MustGet("pool").(*pgxpool.Pool)
	storedCreds := &Credentials{}
	err = pool.QueryRow(context.Background(), "SELECT username, password FROM users where username=$1;", credentials.Username).Scan(&storedCreds.Username, &storedCreds.Password)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Could not find user")
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(credentials.Password)); err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	jwtService := service.JWTAuthService()
	token := jwtService.GenerateToken(credentials.Username)

	c.JSON(http.StatusOK, gin.H{"token": token})
}
