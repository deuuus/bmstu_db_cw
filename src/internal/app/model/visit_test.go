package model_test

import (
	"testing"

	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestVisit_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		v       func() *model.Visit
		isValid bool
	}{
		{
			name: "valid (active)",
			v: func() *model.Visit {
				return model.TestVisit(t)
			},
			isValid: true,
		},
		{
			name: "valid (done)",
			v: func() *model.Visit {
				v := model.TestVisit(t)
				v.Status = "Done"

				return v
			},
			isValid: true,
		},
		{
			name: "invalid status",
			v: func() *model.Visit {
				v := model.TestVisit(t)
				v.Status = "invalid"

				return v
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.v().Validate())
			} else {
				assert.Error(t, tc.v().Validate())
			}
		})
	}
}