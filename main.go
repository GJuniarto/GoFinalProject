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
	

	r.POST("/users/register", handlers.Register(db))
	r.POST("/users/login", handlers.Login(db))
	r.Use(middlewares.TokenAuthMiddleware())
	r.PATCH("/users/topup", handlers.TopUp(db))

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed!")
	}
}
