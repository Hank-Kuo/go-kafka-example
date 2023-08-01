package utils

import "golang.org/x/crypto/bcrypt"

func HashText(text string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckTextHash(text, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	return err
}
