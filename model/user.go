package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id          int
	FullName    string
	Password    string
	PhoneNumber string
}

func NewUser(fullName, password, phoneNumber string) *User {
	return &User{
		FullName:    fullName,
		Password:    password,
		PhoneNumber: phoneNumber,
	}
}

func (u *User) Authenticate(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}

	return true
}
