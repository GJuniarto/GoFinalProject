package main

import (
	"GoFinal/config"
	"GoFinal/handlers"
	"GoFinal/middlewares"
	"GoFinal/model"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func main() {
	config := config.NewConfig()

	r := gin.Default()
	dsn := config.GetPostgresConfig()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed connect db!" , err);
	}

	err = db.AutoMigrate(&model.User{}, &model.Product{}, &model.Category{}, &model.TransactionHistory{})

	if err != nil {
		log.Fatal("Failed migrate db!" , err);
	}
	
	// User
	r.POST("/users/register", handlers.Register(db))
	r.POST("/users/login", handlers.Login(db))
	r.Use(middlewares.TokenAuthMiddleware())
	r.PATCH("/users/topup", handlers.TopUp(db))

	// Categories
	r.POST("/categories", middlewares.AdminAuth(), handlers.CreateCategory(db))
	r.GET("/categories", middlewares.AdminAuth(), handlers.GetCategories(db))
	r.PATCH("/categories/:categoryId", middlewares.AdminAuth(), handlers.EditCategory(db))
	r.DELETE("/categories", middlewares.AdminAuth(), handlers.DeleteCategory(db))

	// Products
	r.POST("/products", middlewares.AdminAuth(), handlers.CreateProduct(db))
	r.GET("/products", handlers.GetProducts(db))
	r.PUT("/products/:productId", middlewares.AdminAuth(), handlers.EditProduct(db))
	r.DELETE("/products/:productId", middlewares.AdminAuth(), handlers.DeleteProduct(db))

	r.POST("/transactions", handlers.CreateTransaction(db))
	r.GET("/transactions/my-transactions", handlers.GetMyTransactions(db))
	r.GET("/transactions/user-transactions", middlewares.AdminAuth(), handlers.GetUserTransactions(db))

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed!")
	}
}
