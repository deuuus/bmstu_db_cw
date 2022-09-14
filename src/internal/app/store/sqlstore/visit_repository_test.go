package sqlstore_test

import (
	"testing"

	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestVisitRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("visits, doctors, patients")

	s := sqlstore.New(db)

	v := model.TestVisit(t)

	assert.Error(t, s.Visit().Create(v))

	p := model.TestPatient(t)
	s.Patient().Create(p)
	d := model.TestDoctor(t)
	s.Doctor().Create(d)

	assert.Error(t, s.Visit().Create(v))

	v.Doctor_id = d.ID
	v.Patient_id = p.ID

	assert.NoError(t, s.Visit().Create(v))
	assert.NotNil(t, v)

	v.Doctor_id = 2
	assert.Error(t, s.Visit().Create(v))
}

func TestVisitRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("visits, doctors, patients")

	s := sqlstore.New(db)
	v1 := model.TestVisit(t)

	p := model.TestPatient(t)
	s.Patient().Create(p)
	d := model.TestDoctor(t)
	s.Doctor().Create(d)

	v1.Doctor_id = d.ID
	v1.Patient_id = p.ID

	s.Visit().Create(v1)

	v2, err := s.Visit().Find(v1.ID)

	assert.NoError(t, err)
	assert.NotNil(t, v2)
	assert.Equal(t, v2.Status, "Active")

	v2, err = s.Visit().Find(v1.ID + 1)
	
	assert.Error(t, err)
}

func TestVisitRepository_CommitVisit(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("visits, doctors, patients")

	s := sqlstore.New(db)

	p := model.TestPatient(t)
	s.Patient().Create(p)
	d := model.TestDoctor(t)
	s.Doctor().Create(d)

	v := model.TestVisit(t)
	v.Doctor_id = d.ID
	v.Patient_id = p.ID
	s.Visit().Create(v)

	assert.NoError(t, s.Visit().CommitVisit(v.ID))
	v, _ = s.Visit().Find(v.ID)
	assert.Equal(t, v.Status, "Done")
}

func TestVisitRepository_GetAllVisitsByDoctor(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("visits, doctors, patients")

	s := sqlstore.New(db)

	p1 := model.TestPatient(t)
	s.Patient().Create(p1)
	d1 := model.TestDoctor(t)
	s.Doctor().Create(d1)

	p2 := model.TestPatient(t)
	p2.Email = "patient@mail.ru"
	s.Patient().Create(p2)
	d2 := model.TestDoctor(t)
	d2.Email = "doctor@mail.ru"
	s.Doctor().Create(d2)

	v := model.TestVisit(t)
	v.Doctor_id = d1.ID
	v.Patient_id = p1.ID
	s.Visit().Create(v)

	v = model.TestVisit(t)
	v.Doctor_id = d1.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)
	s.Visit().CommitVisit(v.ID)
	v, _ = s.Visit().Find(v.ID)

	v = model.TestVisit(t)
	v.Doctor_id = d2.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)

	v = model.TestVisit(t)
	v.Doctor_id = d2.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)

	u, err := s.Visit().GetAllVisitsByDoctor(d1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 2)

	u, err = s.Visit().GetAllVisitsByDoctor(d2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 2)
}

func TestVisitRepository_GetActiveVisitsByDoctor(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("visits, doctors, patients")

	s := sqlstore.New(db)

	p1 := model.TestPatient(t)
	s.Patient().Create(p1)
	d1 := model.TestDoctor(t)
	s.Doctor().Create(d1)

	p2 := model.TestPatient(t)
	p2.Email = "patient@mail.ru"
	s.Patient().Create(p2)
	d2 := model.TestDoctor(t)
	d2.Email = "doctor@mail.ru"
	s.Doctor().Create(d2)

	v := model.TestVisit(t)
	v.Doctor_id = d1.ID
	v.Patient_id = p1.ID
	s.Visit().Create(v)

	v = model.TestVisit(t)
	v.Doctor_id = d1.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)
	s.Visit().CommitVisit(v.ID)
	v, _ = s.Visit().Find(v.ID)

	v = model.TestVisit(t)
	v.Doctor_id = d2.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)

	v = model.TestVisit(t)
	v.Doctor_id = d2.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)

	u, err := s.Visit().GetActiveVisitsByDoctor(d1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 1)

	u, err = s.Visit().GetActiveVisitsByDoctor(d2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 2)
}

func TestVisitRepository_GetDoneVisitsByDoctor(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("visits, doctors, patients")

	s := sqlstore.New(db)

	p1 := model.TestPatient(t)
	s.Patient().Create(p1)
	d1 := model.TestDoctor(t)
	s.Doctor().Create(d1)

	p2 := model.TestPatient(t)
	p2.Email = "patient@mail.ru"
	s.Patient().Create(p2)
	d2 := model.TestDoctor(t)
	d2.Email = "doctor@mail.ru"
	s.Doctor().Create(d2)

	v := model.TestVisit(t)
	v.Doctor_id = d1.ID
	v.Patient_id = p1.ID
	s.Visit().Create(v)

	v = model.TestVisit(t)
	v.Doctor_id = d1.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)
	s.Visit().CommitVisit(v.ID)
	v, _ = s.Visit().Find(v.ID)

	v = model.TestVisit(t)
	v.Doctor_id = d2.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)

	v = model.TestVisit(t)
	v.Doctor_id = d2.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)

	u, err := s.Visit().GetDoneVisitsByDoctor(d1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 1)

	u, err = s.Visit().GetDoneVisitsByDoctor(d2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 0)
}

func TestVisitRepository_GetActiveVisitsByPatient(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("visits, doctors, patients")

	s := sqlstore.New(db)

	p1 := model.TestPatient(t)
	s.Patient().Create(p1)
	d1 := model.TestDoctor(t)
	s.Doctor().Create(d1)

	p2 := model.TestPatient(t)
	p2.Email = "patient@mail.ru"
	s.Patient().Create(p2)
	d2 := model.TestDoctor(t)
	d2.Email = "doctor@mail.ru"
	s.Doctor().Create(d2)

	v := model.TestVisit(t)
	v.Doctor_id = d1.ID
	v.Patient_id = p1.ID
	s.Visit().Create(v)

	v = model.TestVisit(t)
	v.Doctor_id = d1.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)
	s.Visit().CommitVisit(v.ID)
	v, _ = s.Visit().Find(v.ID)

	v = model.TestVisit(t)
	v.Doctor_id = d2.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)

	v = model.TestVisit(t)
	v.Doctor_id = d2.ID
	v.Patient_id = p2.ID
	s.Visit().Create(v)

	u, err := s.Visit().GetActiveVisitsByPatient(p1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 1)

	u, err = s.Visit().GetActiveVisitsByPatient(p2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, len(u), 2)
}