package handlers

import (
	"GoFinal/helpers"
	"GoFinal/middlewares"
	"GoFinal/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		var foundUser model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(user);
		if err := db.Where("email = ?", user.Email).First(&foundUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
			return
		}
		user.Balance = 0
		user.Role = "customer"

		hassedPassword, err := helpers.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password", "error": err})
			return
		}
		user.Password = hassedPassword
		err = user.Validate();
		
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message" : "Validation Error" ,"error": err.Error()})
			return
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":         user.ID,
			"full_name":  user.FullName,
			"email":      user.Email,
			"password":   user.Password,
			"balance":    0,
			"created_at": user.CreatedAt,
		})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		var foundUser model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email not registered"})
			return
		}
		if err := helpers.VerifyPassword(user.Password, foundUser.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong password"})
			return
		}
		token, err := middlewares.CreateToken(foundUser.Email, int(foundUser.ID), foundUser.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func TopUp(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		var foundUser model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get email"})
			return
		}


		if err := db.Where("email = ?", email).First(&foundUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email not registered"})
			return
		}
		foundUser.Balance = foundUser.Balance + user.Balance

		if err := db.Save(&foundUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user balance", "error": err})
			return
		}

		response:=  fmt.Sprintf("Your balance has been succesfully updated to Rp. %d", foundUser.Balance);
		c.JSON(http.StatusOK, gin.H{
			"message": response,
		})
	}
}