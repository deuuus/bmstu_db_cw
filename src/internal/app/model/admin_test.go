package model_test

import (
	"testing"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestAdmin_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		a       func() *model.Admin
		isValid bool
	}{
		{
			name: "valid",
			a: func() *model.Admin {
				return model.TestAdmin(t)
			},
			isValid: true,
		},
		{
			name: "with encrypted password",
			a: func() *model.Admin {
				a := model.TestAdmin(t)
				a.Password = ""
				a.EncryptedPassword = "encryptedpassword"

				return a
			},
			isValid: true,
		},
		{
			name: "empty email",
			a: func() *model.Admin {
				a := model.TestAdmin(t)
				a.Email = ""

				return a
			},
			isValid: false,
		},
		{
			name: "invalid email",
			a: func() *model.Admin {
				a := model.TestAdmin(t)
				a.Email = "invalid"

				return a
			},
			isValid: false,
		},
		{
			name: "empty password",
			a: func() *model.Admin {
				a := model.TestAdmin(t)
				a.Password = ""

				return a
			},
			isValid: false,
		},
		{
			name: "short password",
			a: func() *model.Admin {
				a := model.TestAdmin(t)
				a.Password = "short"

				return a
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.a().Validate())
			} else {
				assert.Error(t, tc.a().Validate())
			}
		})
	}
}

func TestAdmin_BeforeCreate(t *testing.T) {
	a := model.TestAdmin(t)
	assert.NoError(t, a.BeforeCreate())
	assert.NotEmpty(t, a.EncryptedPassword)
}