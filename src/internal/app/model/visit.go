package model

import validation "github.com/go-ozzo/ozzo-validation"

type Visit struct {
	ID                int    `json:"id"`
	Status            string `json:"status"`
	Doctor_id         int    `json:"doctor_id"`
	Patient_id        int    `json:"patient_id"`
}

type VisitView struct {
	Visit   *Visit
	Patient *Patient
}

type ActiveVisitView struct {
	VPs             []VisitView
	Diseases        []*Disease
	Medicines       []*Medicine
}

type ActivePatientVisitView struct {
	Visit          *Visit
	Doctor         *Doctor
	Specialization *Specialization
}

func (p *Visit) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.Status, validation.In("Active", "Done")),
	)
}
