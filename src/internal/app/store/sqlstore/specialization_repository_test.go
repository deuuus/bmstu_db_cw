package sqlstore_test

import (
	"testing"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestSpecializationRepository_Find(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)

	s := sqlstore.New(db)
	u, err := s.Specialization().Find(5)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestSpecializationRepository_FindByName(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)

	s := sqlstore.New(db)
	_, err := s.Specialization().FindByName("Aaaa")
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u, err := s.Specialization().FindByName("Офтальмолог")
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestSpecializationRepository_GetAll(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)

	s := sqlstore.New(db)

	u, err := s.Specialization().GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, u)
}