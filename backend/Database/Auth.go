package Database

import (
	"assaultrifle/Form"
	"assaultrifle/Utils"
	"errors"
)

func Login(email, password string) (string, error) {
	var user Form.User
	result := DB.Where("email = ? AND password = ?", email, password).First(&user)
	if result.Error != nil {
		return "", errors.New("Kullanıcı bulunamadı")
	}
	return user.Token, nil
}

func Register(email, password, username string) (bool, string) {
	var existing Form.User
	if err := DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return false, ""
	}

	token, _:= Utils.Token(10)
	if (token == "err") {
		return false, ""
	}
	
	user := Form.User{
		Username: username,
		Password: password,
		Email:    email,
		Token:    token,
	}

	if err := DB.Create(&user).Error; err != nil {
		return false, ""
	}

	return true, token
}

func Users(token string) (Form.User, error) {
	var user Form.User
	if err := DB.Where("token = ?", token).First(&user).Error; err != nil {
		return Form.User{}, errors.New("kullanıcı bulunamadı")
	}
	return user, nil
}

