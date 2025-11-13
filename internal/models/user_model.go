package models

type UserModel struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

var EmptyUserModel = UserModel{}

func NewUserModel(id int, login string, password string, current int64, withdrawn int64) *UserModel {
	return &UserModel{
		ID:    id,
		Login: login,
	}
}
