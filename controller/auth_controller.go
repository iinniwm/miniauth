package controller

import (
	"net/http"


    "miniauth/config"
    "miniauth/model"
    "miniauth/utils"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
    var body struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    hashed, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

    user := model.User{
        Email:    body.Email,
        Password: string(hashed),
    }

    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func Login(c *gin.Context) {
    var body struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
        return
    }

    var user model.User
    if err := config.DB.First(&user, "email = ?", body.Email).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // ✅ Generate both tokens
    accessToken, err := utils.GenerateJWT(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Access token generation failed"})
        return
    }

    refreshToken, err := utils.GenerateRefreshToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Refresh token generation failed"})
        return
    }

    // 🍪 Set tokens in HttpOnly Cookies
    c.SetCookie("Authorization", accessToken, 3600, "", "", true, true)
    c.SetCookie("RefreshToken", refreshToken, 7*24*3600, "", "", true, true)

    c.JSON(http.StatusOK, gin.H{"message": "Login success"})
}

func Profile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "You are logged in!",
        "user_id": userID,
    })
}

func RefreshToken(c *gin.Context) {
    refreshToken, err := c.Cookie("RefreshToken")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
        return
    }

    claims, err := utils.ValidateJWT(refreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
        return
    }

    // สร้าง access token ใหม่
    accessToken, _ := utils.GenerateJWT(claims.UserID)
    c.SetCookie("Authorization", accessToken, 3600, "", "", true, true)

    c.JSON(http.StatusOK, gin.H{"message": "Access token refreshed"})
}

func Logout(c *gin.Context) {
    // ลบทั้ง 2 cookies โดยตั้ง MaxAge = -1
    c.SetCookie("Authorization", "", -1, "", "", true, true)
    c.SetCookie("RefreshToken", "", -1, "", "", true, true)

    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}