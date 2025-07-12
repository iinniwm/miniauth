package routes
import (
	"github.com/gin-gonic/gin"
	"miniauth/controller"
	"miniauth/middleware"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/register", controller.Register)
    router.POST("/login", controller.Login)
    router.POST("/refresh", controller.RefreshToken)

    auth := router.Group("/")
    auth.Use(middleware.RequireAuth)
    auth.GET("/profile", controller.Profile)
    auth.POST("/logout", controller.Logout)
	
}