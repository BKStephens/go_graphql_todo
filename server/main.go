package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	r.POST("/login", LoginHandler)
	r.POST("/signup", SignupHandler)

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
	err = InsertUser(signUpRequest, pool)
	if err != nil {
		c.String(http.StatusInternalServerError, "Signup failed")
		return
	}
	c.String(http.StatusOK, "Signup successful")
}

func InsertUser(signUpRequest *SignUpRequest, pool *pgxpool.Pool) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), 8)
	if _, err := pool.Exec(context.Background(), "INSERT INTO users(username, email, password) VALUES ($1, $2, $3)", signUpRequest.Username, signUpRequest.Email, string(hashedPassword)); err != nil {
		return err
	}
	return nil
}

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
	c.String(http.StatusOK, "User Logged In")
}
