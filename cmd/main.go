// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"os"
	"osohub/db"
	_ "osohub/docs" // swaggo docs
	"osohub/handlers"
	"osohub/middleware"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title OSOHUB API
// @version 1.0
// @description API for OSOHUB using Cassandra
// @host localhost:8080
// @BasePath /

// @tag.name Auth & Users
// @tag.description Endpoints for authentication and user management
// @tag.name Images
// @tag.description Endpoints for image management and feed

func main() {
	// Cargar variables de entorno desde .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables de entorno del sistema")
	}
	log.Println("JWT_SECRET:", os.Getenv("JWT_SECRET")) // DEBUG: Verifica el valor cargado

	middleware.InitJWTSecret() // Inicializa el secreto JWT después de cargar el .env

	db.InitCassandra()
	defer db.Session.Close()

	r := gin.Default()
	r.POST("/images/:image_id/like", middleware.AuthMiddleware(), handlers.LikeImage)
	r.DELETE("/images/:image_id/like", middleware.AuthMiddleware(), handlers.UnlikeImage)
	r.GET("/images/:image_id/likes/count", handlers.GetImageLikesCount)
	r.DELETE("/images/:image_id", middleware.AuthMiddleware(), handlers.DeleteImage)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/users/:user_id", handlers.GetUserByID)
	r.POST("/users", handlers.CreateUser)
	r.PATCH("/users/:user_id/ban", handlers.BanUser)
	r.POST("/auth/login", handlers.Login)
	r.GET("/images/byid/:image_id", handlers.GetImageByIDByOnlyID)
	r.POST("/images", handlers.UploadImage)
	r.GET("/users/:user_id/images", handlers.GetImagesByUser)
	r.GET("/feed", handlers.GetFeed)
	r.POST("/images/:image_id/report", middleware.AuthMiddleware(), handlers.ReportImage)
	r.GET("/images/:image_id/reports/count", handlers.GetImageReportsCount)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
