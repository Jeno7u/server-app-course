package models

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title" binding:"required,min=3,max=80"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=todo in_progress done"`
	Priority    int    `json:"priority" binding:"required,min=1,max=5"`
	OwnerID     int    `json:"owner_id"`
}

type TaskCreate struct {
	Title       string `json:"title" binding:"required,min=3,max=80"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=todo in_progress done"`
	Priority    int    `json:"priority" binding:"required,min=1,max=5"`
}

type TaskStatusUpdate struct {
	Status string `json:"status" binding:"required,oneof=todo in_progress done"`
}

type CurrentUser struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
}
