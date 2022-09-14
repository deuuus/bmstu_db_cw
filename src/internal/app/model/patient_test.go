package model_test

import (
	"testing"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestPatient_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		p       func() *model.Patient
		isValid bool
	}{
		{
			name: "valid",
			p: func() *model.Patient {
				return model.TestPatient(t)
			},
			isValid: true,
		},
		{
			name: "with encrypted password",
			p: func() *model.Patient {
				p := model.TestPatient(t)
				p.Password = ""
				p.EncryptedPassword = "encryptedpassword"

				return p
			},
			isValid: true,
		},
		{
			name: "empty email",
			p: func() *model.Patient {
				p := model.TestPatient(t)
				p.Email = ""

				return p
			},
			isValid: false,
		},
		{
			name: "invalid email",
			p: func() *model.Patient {
				p := model.TestPatient(t)
				p.Email = "invalid"

				return p
			},
			isValid: false,
		},
		{
			name: "empty password",
			p: func() *model.Patient {
				p := model.TestPatient(t)
				p.Password = ""

				return p
			},
			isValid: false,
		},
		{
			name: "short password",
			p: func() *model.Patient {
				p := model.TestPatient(t)
				p.Password = "short"

				return p
			},
			isValid: false,
		},
		{
			name: "invalid birth_year (negative)",
			p: func() *model.Patient {
				p := model.TestPatient(t)
				p.Birth_year = -1000

				return p
			},
			isValid: false,
		},
		{
			name: "invalid birth_year (less than 1922)",
			p: func() *model.Patient {
				p := model.TestPatient(t)
				p.Birth_year = 1900

				return p
			},
			isValid: false,
		},
		{
			name: "invalid birth_year (more than 2022)",
			p: func() *model.Patient {
				p := model.TestPatient(t)
				p.Birth_year = 2222

				return p
			},
			isValid: false,
		},
		{
			name: "invalid gender",
			p: func() *model.Patient {
				p := model.TestPatient(t)
				p.Gender = "Female"

				return p
			},
			isValid: false,
		},
		{
			name: "invalid gender (case)",
			p: func() *model.Patient {
				p := model.TestPatient(t)
				p.Gender = "жЕнский"

				return p
			},
			isValid: false,
		},
		{
			name: "invalid phone",
			p: func() *model.Patient {
				p := model.TestPatient(t)
				p.Gender = "123456789"

				return p
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.p().Validate())
			} else {
				assert.Error(t, tc.p().Validate())
			}
		})
	}
}

func TestPatient_BeforeCreate(t *testing.T) {
	p := model.TestPatient(t)
	assert.NoError(t, p.BeforeCreate())
	assert.NotEmpty(t, p.EncryptedPassword)
}