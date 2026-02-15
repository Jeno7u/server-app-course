package model


type Feedback struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
	Message string `json:"message" binding:"required,min=10,max=500,notcontains=кринж;рофл;вайб"`
}