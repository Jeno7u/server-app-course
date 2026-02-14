package model


type UsersResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type UserPayload struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

type UserResponse struct {
	Name string `json:"name"`
	Age int `json:"age"`
	IsAdult bool `json:"is_adult"`
}