package model

type (
	UsersRegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	UsersLoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	TodoCreateRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)
