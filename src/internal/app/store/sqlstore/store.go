package sqlstore

import (
	"database/sql"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
)

type Store struct {
	db                       *sql.DB
	patientRepository        *PatientRepository
	doctorRepository         *DoctorRepository
	adminRepository          *AdminRepository
	diseaseRepository        *DiseaseRepository
	medicineRepository       *MedicineRepository
	recordRepository         *RecordRepository
	visitRepository          *VisitRepository
	specializationRepository *SpecializationRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Patient() store.PatientRepository {
	if s.patientRepository != nil {
		return s.patientRepository
	}
	s.patientRepository = &PatientRepository{
		store: s,
	}
	return s.patientRepository
}

func (s *Store) Doctor() store.DoctorRepository {
	if s.doctorRepository != nil {
		return s.doctorRepository
	}
	s.doctorRepository = &DoctorRepository{
		store: s,
	}
	return s.doctorRepository
}

func (s *Store) Admin() store.AdminRepository {
	if s.adminRepository != nil {
		return s.adminRepository
	}
	s.adminRepository = &AdminRepository{
		store: s,
	}
	return s.adminRepository
}

func (s *Store) Specialization() store.SpecializationRepository {
	if s.specializationRepository != nil {
		return s.specializationRepository
	}
	s.specializationRepository = &SpecializationRepository{
		store: s,
	}
	return s.specializationRepository
}

func (s *Store) Visit() store.VisitRepository {
	if s.visitRepository != nil {
		return s.visitRepository
	}
	s.visitRepository = &VisitRepository{
		store: s,
	}
	return s.visitRepository
}

func (s *Store) Disease() store.DiseaseRepository {
	if s.diseaseRepository != nil {
		return s.diseaseRepository
	}
	s.diseaseRepository = &DiseaseRepository{
		store: s,
	}
	return s.diseaseRepository
}

func (s *Store) Medicine() store.MedicineRepository {
	if s.medicineRepository != nil {
		return s.medicineRepository
	}
	s.medicineRepository = &MedicineRepository{
		store: s,
	}
	return s.medicineRepository
}

func (s *Store) Record() store.RecordRepository {
	if s.recordRepository != nil {
		return s.recordRepository
	}
	s.recordRepository = &RecordRepository{
		store: s,
	}
	return s.recordRepository
}
