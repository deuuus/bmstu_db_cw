package sqlstore_test

import (
	"testing"

	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestDoctorRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("doctors")

	s := sqlstore.New(db)
	u := model.TestDoctor(t)
	assert.NoError(t, s.Doctor().Create(u))
	assert.NotNil(t, u.ID)
}

func TestDoctorRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("doctors")

	s := sqlstore.New(db)
	u1 := model.TestDoctor(t)
	s.Doctor().Create(u1)
	u2, err := s.Doctor().Find(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestDoctorRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("doctors")

	s := sqlstore.New(db)
	u1 := model.TestDoctor(t)
	_, err := s.Doctor().FindByEmail(u1.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Doctor().Create(u1)
	u2, err := s.Doctor().FindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestDoctorRepository_GetAll(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("doctors")

	s := sqlstore.New(db)
	u1 := model.TestDoctor(t)
	s.Doctor().Create(u1)

	u2 := model.TestDoctor(t)
	u2.Email = "other@mail.ru"
	s.Doctor().Create(u2)

	u, err := s.Doctor().GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, 2, len(u))
}