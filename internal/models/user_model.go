package models

type UserModel struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewUserModel(id string, login string, password string) *UserModel {
	return &UserModel{
		ID:       id,
		Login:    login,
		Password: password,
	}
}
