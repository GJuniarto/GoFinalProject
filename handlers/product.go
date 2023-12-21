package handlers

import (
	"GoFinal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product model.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()

		err := product.Validate();
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message" : "Validation Error" ,"error": err.Error()})
			return
		}

		if err := db.Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create product", "error": err})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":         product.ID,
			"title":       product.Title,
			"price":      product.Price,
			"stock":      product.Stock,
			"category_id": product.CategoryID,
			"created_at": product.CreatedAt,
		})
	}
}

func GetProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []model.Product
		if err := db.Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch products", "error": err})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

func EditProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product model.Product
		var foundProduct model.Product;

		productId := c.Param("productId")
		if err := db.First(&foundProduct, productId).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
			return
		}

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product.UpdatedAt = time.Now()

		err := product.Validate();
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message" : "Validation Error" ,"error": err.Error()})
			return
		}

		if err := db.Model(&foundProduct).Updates(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product", "error": err})
			return
		}

		responseProduct := map[string]interface{}{
			"id":         foundProduct.ID,
			"title":       foundProduct.Title,
			"price":      foundProduct.Price,
			"stock":      foundProduct.Stock,
			"category_id": foundProduct.CategoryID,
			"created_at": foundProduct.CreatedAt,
			"updated_at": foundProduct.UpdatedAt,
		}

		c.JSON(http.StatusOK, gin.H{
			"product": responseProduct,
		})
	}
}

func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product model.Product
		var foundProduct model.Product;

		productId := c.Param("productId")
		if err := db.First(&foundProduct, productId).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
			return
		}

		if err := db.Delete(&product, productId).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete product", "error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Product has been successfully deleted",
		})
	}
}