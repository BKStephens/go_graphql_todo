package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/bkstephens/go_graphql_todo/server/controller"
	"github.com/bkstephens/go_graphql_todo/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "go_graphql_todo_dev"
)

var DBPool *pgxpool.Pool

func main() {
	Initialize()
}

func Initialize() {
	pool := initDb()
	DBPool = pool
	defer pool.Close()
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("pool", pool)
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
	r.POST("/login", controller.LoginHandler)
	r.POST("/signup", controller.SignupHandler)
	r.GET("/graphql", middleware.AuthorizeJWT(), controller.GraphQLHandler)
	authorized := r.Group("/api/v1")
	authorized.Use(middleware.AuthorizeJWT())
	authorized.GET("/authorized", func(c *gin.Context) {
		c.String(http.StatusOK, "Reached authenticated endpoint")
	})

	var port string
	if envVar := os.Getenv("PORT"); envVar != "" {
		port = envVar
	} else {
		port = "8080"
	}
	r.Run(fmt.Sprintf(":%s", port))
}

func initDb() *pgxpool.Pool {
	var databaseUrl string
	if envVar := os.Getenv("DATABASE_URL"); envVar != "" {
		databaseUrl = envVar
	} else {
		databaseUrl = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	}

	pool, err := pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return pool
}
