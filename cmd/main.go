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
	"strings"

	"github.com/gin-contrib/cors"
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
	// Load environment variables from .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Could not load .env file, using system environment variables")
	}

	log.Println("JWT_SECRET:", os.Getenv("JWT_SECRET")) // DEBUG: Check loaded value

	// Set Gin mode from env (default: debug)
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	middleware.InitJWTSecret() // Initialize JWT secret after loading .env

	db.InitCassandra()
	defer func() {
		sess := db.GetSession()
		if sess != nil {
			sess.Close()
		}
	}()

	r := gin.Default()

	// Configurar CORS para permitir conexiones desde React y Vite
	config := cors.DefaultConfig()

	// Leer orígenes permitidos desde .env
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins != "" {
		config.AllowOrigins = strings.Split(allowedOrigins, ",")
	} else {
		// Fallback por defecto
		config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	}

	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	r.Use(cors.New(config))

	// Middleware adicional para manejar preflight OPTIONS manualmente
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept")
		c.Header("Access-Control-Max-Age", "86400")
		c.Status(200)
	})

	// Configure trusted proxies from .env
	proxies := os.Getenv("TRUSTED_PROXIES")
	proxyList := []string{}
	if proxies != "" {
		proxyList = append(proxyList, proxies)
	}
	if err := r.SetTrustedProxies(proxyList); err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}
	r.POST("/images/:image_id/like", middleware.AuthMiddleware(), handlers.LikeImage)
	r.DELETE("/images/:image_id/like", middleware.AuthMiddleware(), handlers.UnlikeImage)
	r.GET("/images/:image_id/like/status", middleware.AuthMiddleware(), handlers.GetImageLikeStatus)
	r.GET("/images/:image_id/likes/count", handlers.GetImageLikesCount)
	r.DELETE("/images/:image_id", middleware.AuthMiddleware(), handlers.DeleteImage)
	r.GET("/users/me", middleware.AuthMiddleware(), handlers.GetCurrentUser)
	r.PATCH("/users/me", middleware.AuthMiddleware(), handlers.UpdateOwnUser)
	r.GET("/users/me/share-link", middleware.AuthMiddleware(), handlers.GetMyShareLink)

	// Ruta raíz con información de la API
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OSOHUB API",
			"version": "1.0",
			"swagger": "/swagger/index.html",
		})
	})

	// Ruta pública para perfiles (sin autenticación)
	r.GET("/profile/:username", handlers.GetPublicProfile)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/users/:user_id", handlers.GetUserByID)
	r.POST("/users", handlers.CreateUser)
	r.PATCH("/users/:user_id/ban", handlers.BanUser)
	r.POST("/auth/login", handlers.Login)
	r.GET("/images/byid/:image_id", handlers.GetImageByIDByOnlyID)
	r.POST("/images", middleware.AuthMiddleware(), handlers.UploadImage)
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
