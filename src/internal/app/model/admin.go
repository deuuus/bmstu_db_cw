package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID                int    `json:"id"`
	Name		      string `json:"name"`
	Surname           string `json:"surname"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
}

func (a *Admin) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.Email,      validation.Required, is.Email),
		validation.Field(&a.Password,   validation.By(requiredIf(a.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

func (a *Admin) BeforeCreate() error {
	if len(a.Password) > 0 {
		enc, err := encryptString(a.Password)
		if err != nil {
			return err
		}
		a.EncryptedPassword = enc
	}
	return nil
}

func (d *Admin) Sanitize() {
	d.Password = ""
}

func (a *Admin) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password)) == nil
}