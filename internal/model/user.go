package model

type ProfileResponse struct {
	UserId   string `json:"userId" `
	Username string `json:"username" `
}

type User struct {
	UserId   string
	Username string
	Password string
}
