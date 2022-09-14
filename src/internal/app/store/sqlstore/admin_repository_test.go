package sqlstore_test

import (
	"testing"

	"github.com/gopherschool/http-rest-api/internal/app/model"
	"github.com/gopherschool/http-rest-api/internal/app/store"
	"github.com/gopherschool/http-rest-api/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestAdminRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("admins")

	s := sqlstore.New(db)
	u := model.TestAdmin(t)
	assert.NoError(t, s.Admin().Create(u))
	assert.NotNil(t, u.ID)
}

func TestAdminRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("admins")

	s := sqlstore.New(db)
	u1 := model.TestAdmin(t)
	s.Admin().Create(u1)

	u2, err := s.Admin().Find(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)

	u1.ID = -100
	u2, err = s.Admin().Find(u1.ID)
	assert.Error(t, err)
	assert.Nil(t, u2)
}

func TestAdminRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("admins")

	s := sqlstore.New(db)
	u1 := model.TestAdmin(t)
	_, err := s.Admin().FindByEmail(u1.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Admin().Create(u1)
	u2, err := s.Admin().FindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestAdminRepository_GetAll(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("admins")

	s := sqlstore.New(db)
	u1 := model.TestAdmin(t)
	s.Admin().Create(u1)

	u2 := model.TestAdmin(t)
	u2.Email = "other@mail.ru"
	s.Admin().Create(u2)

	u, err := s.Admin().GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, 2, len(u))
}