package middleware

import (
	"fmt"
	"net/http"

	"github.com/bkstephens/go_graphql_todo/server/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > len(BEARER_SCHEMA) {
			tokenString := authHeader[len(BEARER_SCHEMA):]
			token, _ := service.JWTAuthService().ValidateToken(tokenString)
			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				fmt.Println(claims)
				return
			}
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
