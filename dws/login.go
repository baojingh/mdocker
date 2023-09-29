package dws

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
}

var userList = []User{
	{Username: "hadoop", Password: "hadoop"},
}

func authenticate(username string, password string) bool {
	for _, user := range userList {
		if user.Username == username {
			currPd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			// first param should be hash password, second param should be original password
			err := bcrypt.CompareHashAndPassword(currPd, []byte(user.Password))
			return err == nil
		}
	}
	return false
}
