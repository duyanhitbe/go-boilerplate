package hash

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct {
}

func NewBcrypt() Hash {
	return &Bcrypt{}
}

func (b *Bcrypt) Create(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (b *Bcrypt) Compare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
