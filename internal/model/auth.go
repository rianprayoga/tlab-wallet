package model

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=100"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=100"`
}
