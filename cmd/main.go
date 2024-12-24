package main

import (
	"database/sql"
	"fmt"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/shayja/orders-service/config"
	"github.com/shayja/orders-service/docs"
	"github.com/shayja/orders-service/internal/adapters/controllers"
	"github.com/shayja/orders-service/internal/adapters/middleware"
	repositories "github.com/shayja/orders-service/internal/adapters/repositories/orders"
	"github.com/shayja/orders-service/internal/usecases"
	//"github.com/shayja/orders-service/pkg/jwt"
)

// Swagger
//
//  @title                      Orders Microservice
//  @version                    1.0
//  @description                API documentation for the Orders microservice.
//  @contact.name               Shay Jacoby
//  @contact.url                https://github.com/shayja/
//  @contact.email              shayja@gmail.com
//  @license.name               Apache 2.0
//  @license.url                http://www.apache.org/licenses/LICENSE-2.0.html
//  @host                       localhost:8080
//  @BasePath                   /api/v1
//  @schemes                    http https
//  @securityDefinitions.apiKey apiKey
//  @in                         header
//  @name                       Authorization
//  @description                Type "Bearer" followed by a space and JWT token. Please add it in the format "Bearer {AccessToken}" to authorize your requests.
func main() {

	// Load environment variables and configuration
	cfg, err := config.LoadENV()
	if err != nil {
		panic(err)
	}

	if cfg.DBHost == "" {
		panic("DB_HOST environment variable not set")
	}

	// Load environment variables and Format the connection string to the database
	// Connect to database
	db := RegisterDb(cfg)

	// Initialize repository, usecase, and controller
	orderRepo := &repositories.OrderRepository{Db: db}
	orderUsecase := &usecases.OrderUsecase{OrderRepo: orderRepo}
	orderController := &controllers.OrderController{OrderUsecase: orderUsecase}

	// Initialize Gin
	r := gin.Default()

	// Define your secret key for token validation
	secretKey := cfg.AccessTokenSecret

	//GenerateToken(secretKey)

	// Register routes
	RegisterOrderRoutes(r, orderController, secretKey)

	RegisterSwagger(r)

	// Start server
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		panic(err)
	}
}

func RegisterDb(cfg *config.Config) *sql.DB {
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return db
}


/*
// No user auth in this microservice so we call this func to generate a JWT token using the provided secret key, the secret stored in .env file.
func GenerateToken(secretKey string) {
	// User ID and secret key (should be stored securely, not hard-coded)
	userId := "451fa817-41f4-40cf-8dc2-c9f22aa98a4f"

	// Generate JWT
	token, err := jwt.GenerateJWT(userId, secretKey)
	if err != nil {
		panic(err)
	}

	// Print the token
	fmt.Println("Generated Token:", token)
}
*/
func RegisterOrderRoutes(r *gin.Engine, orderController *controllers.OrderController, secretKey string) {
	// Version 1 routes
	orderRoutes := r.Group("/api/v1/order")
	{
		// Apply AuthMiddleware globally or for specific routes
		orderRoutes.Use(middleware.AuthMiddleware(secretKey))

		// Register version 1 order routes
		orderRoutes.GET("", orderController.GetOrders)
		orderRoutes.POST("", orderController.Create)
		orderRoutes.GET(":id", orderController.GetById)
		orderRoutes.PUT(":id/status", orderController.UpdateStatus)
	}
}

func RegisterSwagger(r *gin.Engine) {
	// Swagger setup
	docs.SwaggerInfo.Title = "Go simple Microservice"
	docs.SwaggerInfo.Description = "API documentation for the Go simple Microservice"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}