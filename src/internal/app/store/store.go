package store

type Store interface {
	Patient()        PatientRepository
	Doctor()         DoctorRepository
	Admin()          AdminRepository

	Specialization() SpecializationRepository
	Visit()          VisitRepository
	Disease()        DiseaseRepository
	Medicine()       MedicineRepository
	Record()         RecordRepository
}
