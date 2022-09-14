package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type Doctor struct {
	ID                int    `json:"id"`
	Name		      string `json:"name"`
	Surname           string `json:"surname"`
	Work_since        int    `json:"work_since"`
	Spec_id           int    `json:"spec_id"` 
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
}

type DoctorView struct {
	Doctor         *Doctor
	Specialization *Specialization
}

func (d *Doctor) Validate() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.Email,      validation.Required, is.Email),
		validation.Field(&d.Work_since, validation.Min(1942), validation.Max(2022)),
		validation.Field(&d.Password,   validation.By(requiredIf(d.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

func (d *Doctor) BeforeCreate() error {
	if len(d.Password) > 0 {
		enc, err := encryptString(d.Password)
		if err != nil {
			return err
		}
		d.EncryptedPassword = enc
	}
	return nil
}

func (d *Doctor) Sanitize() {
	d.Password = ""
}

func (d *Doctor) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(d.EncryptedPassword), []byte(password)) == nil
}