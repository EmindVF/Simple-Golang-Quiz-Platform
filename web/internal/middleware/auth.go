package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"quiz_platform/internal/handler/repository"
	"quiz_platform/internal/misc/config"
	"quiz_platform/internal/utility"
)

type SessionData struct {
	UserId      int32
	Permissions int64
	UserName    string
}

// Middleware to check JWT
func GetTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		tokenString, err := c.Cookie("authorization")
		if tokenString == "" || err != nil {
			c.Set("BaseH", &gin.H{
				"username":    "",
				"authorized":  false,
				"permissions": int64(0),
			})
			c.Next()
			return
		}

		token, err := utility.ValidateToken(tokenString, config.GlobalConfig.TokenInfo.PublicKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		permissions, err := repository.UserRepositoryInstance.
			GetUserPermissions(ctx, token.UserId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		user, err := repository.UserRepositoryInstance.
			GetUserById(ctx, token.UserId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("sessionData", &SessionData{
			UserId:      token.UserId,
			Permissions: permissions,
			UserName:    user.UserName,
		})
		c.Set("BaseH", &gin.H{
			"username":    user.UserName,
			"authorized":  true,
			"permissions": permissions,
		})
		c.Next()
	}
}

func RequirePermissionMiddleware(perm int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, ok := c.Get("sessionData")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if sessionData, ok := data.(*SessionData); ok {
			if sessionData.Permissions&perm != perm {
				c.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
