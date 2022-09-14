package sqlstore_test

import (
	"testing"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestDiseaseRepository_Find(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)

	s := sqlstore.New(db) 
	u, err := s.Disease().Find(5)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestDiseaseRepository_GetAll(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)

	s := sqlstore.New(db)
	u, err := s.Disease().GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestDiseaseRepository_GetBySpecialization(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)

	s := sqlstore.New(db)
	u, err := s.Disease().GetBySpecialization(1)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, 3, len(u))
}