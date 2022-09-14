package model_test

import (
	"testing"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestDoctor_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		d       func() *model.Doctor
		isValid bool
	}{
		{
			name: "valid",
			d: func() *model.Doctor {
				return model.TestDoctor(t)
			},
			isValid: true,
		},
		{
			name: "invalid email",
			d: func() *model.Doctor {
				d := model.TestDoctor(t)
				d.Email = "invalid"

				return d
			},
			isValid: false,
		},
		{
			name: "empty password",
			d: func() *model.Doctor {
				d := model.TestDoctor(t)
				d.Password = ""

				return d
			},
			isValid: false,
		},
		{
			name: "empty email",
			d: func() *model.Doctor {
				d := model.TestDoctor(t)
				d.Email = ""

				return d
			},
			isValid: false,
		},
		{
			name: "with encrypted password",
			d: func() *model.Doctor {
				d := model.TestDoctor(t)
				d.Password = ""
				d.EncryptedPassword = "encpassword"

				return d
			},
			isValid: true,
		},
		{
			name: "short password",
			d: func() *model.Doctor {
				d := model.TestDoctor(t)
				d.Password = "short"

				return d
			},
			isValid: false,
		},
		{
			name: "invalid_work_since (negative)",
			d: func() *model.Doctor {
				d := model.TestDoctor(t)
				d.Work_since = -1

				return d
			},
			isValid: false,
		},
		{
			name: "invalid_work_since (more than 2022)",
			d: func() *model.Doctor {
				d := model.TestDoctor(t)
				d.Work_since = 2023

				return d
			},
			isValid: false,
		},
		{
			name: "invalid_work_since (less than 1942)",
			d: func() *model.Doctor {
				d := model.TestDoctor(t)
				d.Work_since = 1941

				return d
			},
			isValid: false,
		},
		{
			name: "invalid spec_id (negative)",
			d: func() *model.Doctor {
				d := model.TestDoctor(t)
				d.Spec_id = -1

				return d
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.d().Validate())
			} else {
				assert.Error(t, tc.d().Validate())
			}
		})
	}
}

func TestDoctor_BeforeCreate(t *testing.T) {
	d := model.TestDoctor(t)
	assert.NoError(t, d.BeforeCreate())
	assert.NotEmpty(t, d.EncryptedPassword)
}