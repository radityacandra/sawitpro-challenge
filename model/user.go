package model

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

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
		PhoneNumber: phoneNumber[1:],
	}
}

func (u *User) Authenticate(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}

	return true
}

func (u *User) UpdateFullName(fullName *string) {
	if fullName != nil && *fullName != u.FullName {
		u.FullName = *fullName
	}
}

func (u *User) UpdatePhoneNumber(phoneNumber *string) {
	if phoneNumber != nil && *phoneNumber != fmt.Sprintf("+%s", u.PhoneNumber) {
		u.PhoneNumber = *phoneNumber
		u.PhoneNumber = u.PhoneNumber[1:]
	}
}

func (u *User) PhoneNumberChanged(phoneNumber *string) bool {
	if phoneNumber != nil && *phoneNumber != fmt.Sprintf("+%s", u.PhoneNumber) {
		return true
	}

	return false
}

func (u *User) GetFormalizedPhoneNumber() string {
	return fmt.Sprintf("+%s", u.PhoneNumber)
}
