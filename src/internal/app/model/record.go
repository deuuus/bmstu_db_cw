package model

type Record struct {
	ID          int    `json:"id"`
	Visit_id    int    `json:"visit_id"`
	Disease_id  int    `json:"disease_id"`
	Medicine_id int    `json:"medicine_id"`
}

type RecordView struct {
	Record         *Record
	Visit          *Visit
	Doctor         *Doctor
	Specialization *Specialization
	Disease        *Disease
	Medicine       *Medicine
}

type RecordDoctorView struct {
	Record         *Record
	Visit          *Visit
	Patient        *Patient
	Disease        *Disease
	Medicine       *Medicine
}