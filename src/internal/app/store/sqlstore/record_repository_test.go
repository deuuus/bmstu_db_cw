package sqlstore_test

import (
	"testing"

	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestRecordRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("records, doctors, patients, visits")

	rec := model.TestRecord(t)

	s := sqlstore.New(db)

	p := model.TestPatient(t)
	s.Patient().Create(p)

	d := model.TestDoctor(t)
	s.Doctor().Create(d)

	v := model.TestVisit(t)
	v.Doctor_id = d.ID
	v.Patient_id = p.ID

	assert.Error(t, s.Record().Create(rec))

	s.Visit().Create(v)

	rec.Visit_id = v.ID

	assert.NoError(t, s.Record().Create(rec))
	assert.NotNil(t, rec)

	rec.Visit_id = 100
	assert.Error(t, s.Record().Create(rec))
}

func TestRecordRepository_GetAllByDoctor(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("records, doctors, patients, visits")

	s := sqlstore.New(db)

	p := model.TestPatient(t)
	s.Patient().Create(p)

	d1 := model.TestDoctor(t)
	s.Doctor().Create(d1)

	d2 := model.TestDoctor(t)
	d2.Email = "other@mail.ru"
	s.Doctor().Create(d2)

	v := model.TestVisit(t)
	v.Status = "Done"
	v.Doctor_id = d1.ID
	v.Patient_id = p.ID

	rec := model.TestRecord(t)

	u, err := s.Record().GetAllByDoctor(d1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 0)

	s.Visit().Create(v)
	rec.Visit_id = v.ID
	s.Record().Create(rec)

	s.Visit().Create(v)
	rec.Visit_id = v.ID
	s.Record().Create(rec)

	u, err = s.Record().GetAllByDoctor(d2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 0)

	v.Doctor_id = d2.ID
	s.Visit().Create(v)
	rec.Visit_id = v.ID
	s.Record().Create(rec)

	s.Visit().Create(v)
	rec.Visit_id = v.ID
	s.Record().Create(rec)

	u, err = s.Record().GetAllByDoctor(d1.ID)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 2)
}

func TestRecordRepository_GetAllByPatient(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("records, doctors, patients, visits")

	s := sqlstore.New(db)

	p1 := model.TestPatient(t)
	s.Patient().Create(p1)

	p2 := model.TestPatient(t)
	p2.Email = "other@mail.ru"
	s.Patient().Create(p2)

	d := model.TestDoctor(t)
	s.Doctor().Create(d)

	v := model.TestVisit(t)
	v.Status = "Done"
	v.Doctor_id = d.ID
	v.Patient_id = p1.ID

	rec := model.TestRecord(t)

	u, err := s.Record().GetAllByPatient(p1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 0)

	s.Visit().Create(v)
	rec.Visit_id = v.ID
	s.Record().Create(rec)

	s.Visit().Create(v)
	rec.Visit_id = v.ID
	s.Record().Create(rec)

	u, err = s.Record().GetAllByPatient(p2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 0)

	v.Patient_id = p2.ID
	s.Visit().Create(v)
	rec.Visit_id = v.ID
	s.Record().Create(rec)

	s.Visit().Create(v)
	rec.Visit_id = v.ID
	s.Record().Create(rec)

	u, err = s.Record().GetAllByPatient(p1.ID)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 2)
}