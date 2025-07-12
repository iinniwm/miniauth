package middleware
import (
    "net/http"
    "strings"

    "miniauth/utils"

    "github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
    token, err := c.Cookie("Authorization")
    if err != nil || strings.TrimSpace(token) == "" {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    claims, err := utils.ValidateJWT(token)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    // แนบ user ID ลง context เผื่อใช้ต่อ
    c.Set("user_id", claims.UserID)
    c.Next()
}