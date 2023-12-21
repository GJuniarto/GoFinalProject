package handlers

import (
	"GoFinal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category model.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		category.SoldProductAmount = 0
		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()

		err := category.Validate();
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message" : "Validation Error" ,"error": err.Error()})
		}

		if err := db.Create(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create category", "error": err})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":         category.ID,
			"type":       category.Type,
			"sold_product_amount":  category.SoldProductAmount,
			"created_at": category.CreatedAt,
		})
	}
}

func GetCategories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []model.Category
		if err := db.Preload("Products").Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch categories", "error": err})
			return
		}
		c.JSON(http.StatusOK, categories)
	}
}

func EditCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category model.Category
		var foundCategory model.Category;

		categoryId := c.Param("categoryId")
		if err := db.First(&foundCategory, categoryId).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
		
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := category.Validate();
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message" : "Validation Error" ,"error": err.Error()})
		}

		foundCategory.Type = category.Type
		if err := db.Save(&foundCategory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to edit category", "error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":         foundCategory.ID,
			"type":       foundCategory.Type,
			"sold_product_amount":  foundCategory.SoldProductAmount,
			"created_at": foundCategory.CreatedAt,
			"updated_at": foundCategory.UpdatedAt,
		})
	}
}

func DeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category model.Category
		var foundCategory model.Category;

		categoryId := c.Param("categoryId")
		if err := db.First(&foundCategory, categoryId).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}

		if err := db.Delete(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete category", "error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Category has been successfully deleted ",})
	}
}