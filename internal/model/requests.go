package model

type UsersLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
