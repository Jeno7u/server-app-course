package model

type Product struct {
	ProductID int     `json:"product_id" binding:"required"`
	Name      string  `name:"email" binding:"required,email"`
	Category  string  `category:"age" binding:"omitempty,min=0"`
	Price     float64 `json:"price"`
}
