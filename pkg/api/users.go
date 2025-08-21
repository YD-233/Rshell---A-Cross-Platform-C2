package api

import (
	"BackendTemplate/pkg/database"
	"BackendTemplate/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 登录处理函数
func LoginHandler(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var users database.Users
	if database.Engine.Where("username = ?", loginData.Username).Get(&users); users.Password == loginData.Password {
		token, err := utils.GenerateJWT(loginData.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{
			"token":    token,
			"refresh":  "mock-refresh-token",
			"username": loginData.Username,
		}})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

// JWT 验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization2")[len("Bearer "):]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}

// 注销处理函数
func LogoutHandler(c *gin.Context) {
	// 这里可以处理注销逻辑，比如删除 refresh token
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Logged out successfully"})
}

// 修改密码处理函数
func ChangePasswordHandler(c *gin.Context) {
	var passwordData struct {
		OldPassword string `form:"old_password"`
		NewPassword string `form:"new_password"`
	}
	if err := c.ShouldBind(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	// 处理密码修改逻辑
	if passwordData.OldPassword != passwordData.NewPassword {
		username := c.MustGet("username").(string)
		var users database.Users
		if database.Engine.Where("username = ?", username).Get(&users); users.Password == passwordData.OldPassword {
			users.Password = passwordData.NewPassword
			database.Engine.Where("username = ?", username).Update(&users)
			c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Password changed successfully"})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 400, "message": "Password changed failed"})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "Password changed failed"})
	}

}
