package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Account string `json:"account"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

// Removes "bearer " from auth header
func stripBearer(tok string) (string, error) {
	if len(tok) > 6 && strings.ToLower(tok[0:7]) == "bearer " {
		return tok[7:], nil
	}
	return tok, nil
}

// Check if user has a valid token
func RequiredLoggedIn(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := stripBearer(c.Request.Header.Get("Authorization"))
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		tokenClaims, parseErr := jwt.ParseWithClaims(
			token,
			&Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			},
		)
		if parseErr != nil {
			c.IndentedJSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if tokenClaims != nil {
			claims, ok := tokenClaims.Claims.(*Claims)

			if ok && tokenClaims.Valid {
				c.Set("user_id", claims.Account)
				c.Set("user_role", claims.Role)
				c.Next()
				return
			}
		}

		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
}
