package sqlstore_test

import (
	"testing"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestMedicineRepository_Find(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)

	s := sqlstore.New(db)
	u2, err := s.Medicine().Find(1)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestMedicineRepository_GetAll(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)

	s := sqlstore.New(db)

	u, err := s.Medicine().GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, u)
}
