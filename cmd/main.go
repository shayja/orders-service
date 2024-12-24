package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shayja/orders-service/docs"
	"github.com/shayja/orders-service/internal/adapters/controllers"
	"github.com/shayja/orders-service/internal/adapters/middleware"
	repositories "github.com/shayja/orders-service/internal/adapters/repositories/orders"
	"github.com/shayja/orders-service/internal/usecases"
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

	// load
	LoadENV() 

	host := os.Getenv("DB_HOST")
	if host == "" {
		log.Fatal("DB_HOST environment variable not set")
	}

	// Load environment variables and Format the connection string to the database
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("SSL_MODE"))
	
	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize repository, usecase, and controller
	orderRepo := &repositories.OrderRepository{Db: db}
	orderUsecase := &usecases.OrderUsecase{OrderRepo: orderRepo}
	orderController := &controllers.OrderController{OrderUsecase: orderUsecase}

	// Initialize Gin
	r := gin.Default()

	// Define your secret key for token validation
	secretKey := os.Getenv("ACCESS_TOKEN_SECRET")

	// Register routes
	RegisterOrderRoutes(r, orderController, secretKey)

	RegisterSwagger(r)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func LoadENV() {
	// load .env file
	if err := godotenv.Load(); err != nil {
        fmt.Printf("Error getting env, not comming through %v", err)
    }
}

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