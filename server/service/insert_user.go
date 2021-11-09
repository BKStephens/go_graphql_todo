package service

import (
	"context"

	"github.com/bkstephens/go_graphql_todo/server/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(params types.InsertUserParams, pool *pgxpool.Pool) (int, error) {
	var userId int
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(params.Password), 8)
	err := pool.QueryRow(context.Background(), "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) RETURNING id;", params.Username, params.Email, string(hashedPassword)).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
