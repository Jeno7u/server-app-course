package model

type UserCreateRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Age          *int   `json:"age" binding:"omitempty,min=0"`
	IsSubscribed *bool  `json:"is_subscribed"`
}

type UsersResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserResponse struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	IsAdult bool   `json:"is_adult"`
}
