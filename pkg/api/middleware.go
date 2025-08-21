package api

import (
	"BackendTemplate/pkg/database"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin) // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// Basic认证中间件
func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			// 返回WWW-Authenticate头，触发浏览器的弹框
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatus(401)
			return
		}

		encodedCreds := authHeader[len("Basic "):]
		creds, err := base64.StdEncoding.DecodeString(encodedCreds)
		if err != nil {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatus(401)
			return
		}

		credParts := strings.SplitN(string(creds), ":", 2)
		if len(credParts) != 2 {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatus(401)
			return
		}
		user, pass := credParts[0], credParts[1]

		var user_pass database.Users
		database.Engine.Where("username = ?", user).Get(&user_pass)
		if user_pass.Password != pass || user_pass.Password == "" {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatus(401)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}