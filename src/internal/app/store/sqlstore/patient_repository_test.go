package sqlstore_test

import (
	"testing"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestPatientRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("patients")

	s := sqlstore.New(db)
	u := model.TestPatient(t)
	assert.NoError(t, s.Patient().Create(u))
	assert.NotNil(t, u.ID)
}

func TestPatientRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("patients")

	s := sqlstore.New(db)
	u1 := model.TestPatient(t)
	s.Patient().Create(u1)
	u2, err := s.Patient().Find(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestPatientRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("patients")

	s := sqlstore.New(db)
	u1 := model.TestPatient(t)
	_, err := s.Patient().FindByEmail(u1.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Patient().Create(u1)
	u2, err := s.Patient().FindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestPatientRepository_GetAll(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("patients")

	s := sqlstore.New(db)
	u1 := model.TestPatient(t)
	s.Patient().Create(u1)

	u2 := model.TestPatient(t)
	u2.Email = "other@mail.ru"
	s.Patient().Create(u2)

	u, err := s.Patient().GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, 2, len(u))
}