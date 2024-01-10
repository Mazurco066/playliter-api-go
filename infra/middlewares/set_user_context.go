package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SetUserContext set context if user has valid token otherwise none
func SetUserContext(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := stripBearer(c.Request.Header.Get("Authorization"))

		tokenClaims, _ := jwt.ParseWithClaims(
			token,
			&Claims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			},
		)

		if tokenClaims != nil {
			claims, ok := tokenClaims.Claims.(*Claims)
			if ok && tokenClaims.Valid {
				c.Set("user_email", claims.Account)
				c.Set("user_role", claims.Role)
				c.Request = setToContext(c, "user_email", claims.Account)
				c.Request = setToContext(c, "user_role", claims.Role)
			}
		}

		c.Next()
	}
}

// Set context key-value pair
func setToContext(c *gin.Context, key interface{}, value interface{}) *http.Request {
	return c.Request.WithContext(context.WithValue(c.Request.Context(), key, value))
}
