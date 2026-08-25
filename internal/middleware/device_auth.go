package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

func DeviceKeyAuth(q *repository.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Device-Key")
		if key == "" {
			domain.Fail(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		device, err := q.GetDeviceByAPIKey(c.Request.Context(), key)
		if err != nil {
			domain.Fail(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("deviceId", device.ID)
		c.Set("userId", device.UserID)
		c.Next()
	}
}
