package store

import "github.com/deuuus/bmstu_db_cw/src/internal/app/model"

type PatientRepository interface {
	Create(*model.Patient) error
	Find(int)              (*model.Patient, error)
	FindByEmail(string)    (*model.Patient, error)
	GetAll()               ([]*model.Patient, error)
}

type DoctorRepository interface {
	Create(*model.Doctor) error
	Find(int) (*model.Doctor, error)
	FindByEmail(string) (*model.Doctor, error)
	GetAll() ([]*model.Doctor, error)
}

type AdminRepository interface {
	Create(*model.Admin) error
	Find(int) (*model.Admin, error)
	FindByEmail(string) (*model.Admin, error)
	GetAll() ([]*model.Admin, error)
}

type VisitRepository interface {
	Create(*model.Visit) error
	Find(int) (*model.Visit, error)
	GetActiveVisitsByDoctor(int) ([]*model.Visit, error)
	GetDoneVisitsByDoctor(int) ([]*model.Visit, error)
	GetAllVisitsByDoctor(int) ([]*model.Visit, error)
	GetActiveVisitsByPatient(patient_id int) ([]*model.Visit, error) 
	CommitVisit(int) error
}

type MedicineRepository interface {
	Find(int) (*model.Medicine, error)
	GetAll() ([]*model.Medicine, error)
}

type DiseaseRepository interface {
	Find(int) (*model.Disease, error)
	GetAll() ([]*model.Disease, error)
	GetBySpecialization(int) ([]*model.Disease, error)
	Percent(spec_id int) ([]int, []float64, error)
}

type SpecializationRepository interface {
	Find(int)               (*model.Specialization, error)
	FindByName(name string) (*model.Specialization, error)
	GetAll() ([]*model.Specialization, error)
}

type RecordRepository interface {
	Create(*model.Record) error
	GetAllByPatient(patient_id int) ([]*model.Record, error) 
	GetAllByDoctor(doctor_id int) ([]*model.Record, error) 
}