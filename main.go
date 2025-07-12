// main.go
package main

import (
	"github.com/gin-gonic/gin"
	"miniauth/routes"
	"miniauth/model"
	"miniauth/config"
	"github.com/gin-contrib/cors"
	"log"
)
func main() {
    /*config.Port = os.Getenv("PORT")
	config.DBUrl = os.Getenv("DB_URL")
	config.JwtSecret = os.Getenv("JWT_SECRET")
	if config.Port == "" || config.DBUrl == "" || config.JwtSecret == "" {
    log.Fatal("❌ Missing environment variables")
	}*/

	config.LoadEnv()
	config.ConnectDB()
	if err := config.DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("❌ Failed to auto-migrate: %v", err)
	}

	log.Println("✅ Auto-migration completed")


	// สร้าง router และกำหนดเส้นทาง
	gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.DebugMode) // ใช้สำหรับการพัฒนา	
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

    routes.AuthRoutes(router)
    router.Run(":" + config.Port)

	

}