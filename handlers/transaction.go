package handlers

import (
	"GoFinal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
type TransactionData struct {
	ProductID uint `json:"product_id"`
	Quantity int `json:"quantity"`
}

func CreateTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactionData TransactionData
		var product model.Product
		var transaction model.TransactionHistory
		var user model.User

		userId, exists := c.Get("id")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		if err := c.ShouldBindJSON(&transactionData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Check data is exists
		if err := db.First(&product, transactionData.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
			return
		}
		// Check stock is available or not
		if product.Stock < transactionData.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock not available"})
			return
		}
		// Total price calculation
		transaction.TotalPrice = product.Price * transactionData.Quantity

		if err := db.First(&user, userId).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		// Check user balance
		if user.Balance < transaction.TotalPrice {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Balance not enough"})
			return
		}

		// Update product stock
		product.Stock = product.Stock - transactionData.Quantity
		if err := db.Save(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product", "error": err})
			return
		}

		// Update user balance
		user.Balance = user.Balance - transaction.TotalPrice
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user balance", "error": err})
			return
		}

		transaction.ProductID = transactionData.ProductID
		transaction.UserID = uint(userId.(float64))
		transaction.Quantity = transactionData.Quantity
		transaction.CreatedAt = time.Now()
		transaction.UpdatedAt = time.Now()

		err := transaction.Validate();
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message" : "Validation Error" ,"error": err.Error()})
			return
		}

		if err := db.Create(&transaction).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create transaction", "error": err})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "You have successfully purchased the product",
			"transaction_bill": gin.H{
				"total_price":   transaction.TotalPrice,
				"quantity":      transaction.Quantity,
				"product_title": product.Title,
			}})
	}
}

func GetMyTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactions []model.TransactionHistory
		userId, exists := c.Get("id")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}
		if err := db.Where("user_id = ?", userId).Find(&transactions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch transactions", "error": err})
			return
		}
		responseTransactions := make([]map[string]interface{}, len(transactions))

		for index, transaction := range transactions {
			var product model.Product;
			if err := db.First(&product, transaction.ProductID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
				return
			}
			responseTransactions[index] = map[string]interface{}{
				"id":          transaction.ID,
				"product_id":  transaction.ProductID,
				"user_id":     transaction.UserID,
				"quantity":    transaction.Quantity,
				"total_price": transaction.TotalPrice,
				"Product": map[string]interface{}{
					"id":          product.ID,
					"title":       product.Title,
					"price":       product.Price,
					"stock":       product.Stock,
					"category_id": product.CategoryID,
					"created_at":  product.CreatedAt,
					"updated_at":  product.UpdatedAt,
				},
			}
		}

		c.JSON(http.StatusOK, responseTransactions);
	}
}

func GetUserTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactions []model.TransactionHistory
		if err := db.Find(&transactions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch transactions", "error": err})
			return
		}
		responseTransactions := make([]map[string]interface{}, len(transactions))

		for index, transaction := range transactions {
			var product model.Product;
			var user model.User;
			if err := db.First(&product, transaction.ProductID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
				return
			}
			if err := db.First(&user, transaction.UserID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
				return
			}

			responseTransactions[index] = map[string]interface{}{
				"id":          transaction.ID,
				"product_id":  transaction.ProductID,
				"user_id":     transaction.UserID,
				"quantity":    transaction.Quantity,
				"total_price": transaction.TotalPrice,
				"Product": map[string]interface{}{
					"id":          product.ID,
					"title":       product.Title,
					"price":       product.Price,
					"stock":       product.Stock,
					"category_id": product.CategoryID,
					"created_at":  product.CreatedAt,
					"updated_at":  product.UpdatedAt,
				},
				"User": map[string]interface{}{
					"id":         user.ID,
					"email":      user.Email,
					"full_name":  user.FullName,
					"balance":    user.Balance,
					"created_at": user.CreatedAt,
					"updated_at": user.UpdatedAt,
				},
			}
			
		}

		c.JSON(http.StatusOK, responseTransactions);
	}
}