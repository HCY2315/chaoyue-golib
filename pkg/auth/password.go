package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Password string

const (
	defaultCost = bcrypt.DefaultCost
)

func NewPasswordFromString(s string) Password {
	return Password(s)
}

func (p *Password) MustEncryptFrom(plainText string) *Password {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainText), defaultCost)
	if err != nil {
		fmt.Printf("加密Password失败:%s\n", err.Error())
		return p
	}
	*p = Password(hashed)
	return p
}

func (p Password) MatchWithPlain(plainText string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p), []byte(plainText)) == nil
}

func (p *Password) String() string {
	return string(*p)
}
