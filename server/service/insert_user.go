package service

import (
	"context"

	"github.com/bkstephens/go_graphql_todo/server/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(params types.InsertUserParams, pool *pgxpool.Pool) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(params.Password), 8)
	if _, err := pool.Exec(context.Background(), "INSERT INTO users(username, email, password) VALUES ($1, $2, $3)", params.Username, params.Email, string(hashedPassword)); err != nil {
		return err
	}
	return nil
}
