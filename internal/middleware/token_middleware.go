package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rustingoff/admin_panel_rep/pkg/jwt"
	"net/http"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := jwt.TokenValid(c.Request)
		if err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", nil)
			c.Abort()
		}

		c.Next()
	}
}
