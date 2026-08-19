package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Heiji57/ETB-BE/config"
	"github.com/Heiji57/ETB-BE/internal/domain"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			domain.Fail(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, domain.ErrUnauthorized
			}
			return []byte(cfg.JWT.AccessSecret), nil
		})
		if err != nil || !token.Valid {
			domain.Fail(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			domain.Fail(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			domain.Fail(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("userId", userID)
		c.Next()
	}
}
