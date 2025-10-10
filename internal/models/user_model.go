package models

type UserModel struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	Current   int64  `json:"current"`
	Withdrawn int64  `json:"withdrawn"`
}

func NewUserModel(id string, login string, password string, current int64, withdrawn int64) *UserModel {
	return &UserModel{
		ID:        id,
		Login:     login,
		Password:  password,
		Current:   current,
		Withdrawn: withdrawn,
	}
}
