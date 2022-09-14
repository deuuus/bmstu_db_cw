package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type Patient struct {
	ID                int    `json:"id"`
	Name		      string `json:"name"`
	Surname           string `json:"surname"`
	Gender            string `json:"gender"`
	Birth_year        int    `json:"birth_year"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
}

func (p *Patient) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.Email,      validation.Required, is.Email),
		validation.Field(&p.Gender,     validation.In("Женский", "Мужской")),
		validation.Field(&p.Birth_year, validation.Min(1922), validation.Max(2022)),
		validation.Field(&p.Password,   validation.By(requiredIf(p.EncryptedPassword == "")), validation.Length(6, 15)),
	)
}

func (p *Patient) BeforeCreate() error {
	if len(p.Password) > 0 {
		enc, err := encryptString(p.Password)
		if err != nil {
			return err
		}
		p.EncryptedPassword = enc
	}
	return nil
}

func (p *Patient) Sanitize() {
	p.Password = ""
}

func (p *Patient) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
