package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName  string    `json:"full_name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"`
	Role      string    `json:"role" validate:"required,oneof=admin customer"`
	Balance   int       `json:"balance" validate:"required,min=0,max=100000000"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	gorm.Model
	Title      string    `json:"title" validate:"required"`
	Price      int       `json:"price" validate:"required,min=0,max=50000000"`
	Stock      int       `json:"stock" validate:"required,min=5"`
	CategoryID uint      `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Category   Category  `json:"category" gorm:"foreignKey:CategoryID"`
}

type Category struct {
	gorm.Model
	Type              string    `json:"type" validate:"required"`
	SoldProductAmount int       `json:"sold_product_amount"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type TransactionHistory struct {
	gorm.Model
	ProductID  uint      `json:"product_id"`
	UserID     uint      `json:"user_id"`
	Quantity   int       `json:"quantity" validate:"required"`
	TotalPrice int       `json:"total_price" validate:"required"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Product    Product   `json:"product" gorm:"foreignKey:ProductID"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (p *Product) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (c *Category) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

func (t *TransactionHistory) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
