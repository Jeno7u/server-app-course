package models

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type TodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CustomUserRequest struct {
	Username string  `json:"username" binding:"required"`
	Age      int     `json:"age" binding:"required,gt=18"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=8,max=16"`
	Phone    *string `json:"phone,omitempty"`
}
