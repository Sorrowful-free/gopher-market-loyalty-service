package models

type UserModel struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

var EMPTY_USER_MODEL = UserModel{}

func NewUserModel(id int64, login string, password string, current int64, withdrawn int64) *UserModel {
	return &UserModel{
		ID:       id,
		Login:    login,
		Password: password,
	}
}
