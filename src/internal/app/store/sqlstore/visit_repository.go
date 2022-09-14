package sqlstore

import (
	"database/sql"

	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
)

type VisitRepository struct {
	store *Store
}

func (r *VisitRepository) Find(id int) (*model.Visit, error) {
	v := &model.Visit{}
	if err := r.store.db.QueryRow(
		"SELECT id, status, patient_id, doctor_id FROM visits WHERE id = $1", id,
	).Scan(
		&v.ID,
		&v.Status,
		&v.Patient_id,
		&v.Doctor_id,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return v, nil
}

func (r *VisitRepository) Create(v *model.Visit) error {
	return r.store.db.QueryRow(
		"INSERT INTO visits (status, patient_id, doctor_id) VALUES ('Active', $1, $2) RETURNING id",
		&v.Patient_id,
		&v.Doctor_id,
	).Scan(&v.ID)
}

func (r *VisitRepository) GetActiveVisitsByDoctor(doctor_id int) ([]*model.Visit, error) {
	vs := []*model.Visit{}
	rows, err := r.store.db.Query("SELECT id, patient_id FROM visits WHERE doctor_id = $1 AND status = 'Active'", doctor_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v := &model.Visit{}
		err = rows.Scan(
			&v.ID,
			&v.Patient_id,
		)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}

func (r *VisitRepository) GetDoneVisitsByDoctor(doctor_id int) ([]*model.Visit, error) {
	vs := []*model.Visit{}
	rows, err := r.store.db.Query("SELECT id, patient_id FROM visits WHERE doctor_id = $1 AND status = 'Done'", doctor_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v := &model.Visit{}
		err = rows.Scan(
			&v.ID,
			&v.Patient_id,
		)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}

func (r *VisitRepository) GetAllVisitsByDoctor(doctor_id int) ([]*model.Visit, error) {
	vs := []*model.Visit{}
	rows, err := r.store.db.Query("SELECT id, patient_id FROM visits WHERE doctor_id = $1", doctor_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v := &model.Visit{}
		err = rows.Scan(
			&v.ID,
			&v.Patient_id,
		)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}

func (r *VisitRepository) GetActiveVisitsByPatient(patient_id int) ([]*model.Visit, error) {
	vs := []*model.Visit{}
	rows, err := r.store.db.Query("SELECT id, doctor_id FROM visits WHERE patient_id = $1 AND status = 'Active'", patient_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v := &model.Visit{}
		err = rows.Scan(
			&v.ID,
			&v.Doctor_id,
		)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}

func (r *VisitRepository) CommitVisit(id int) error {
	_, err := r.store.db.Exec("UPDATE visits SET status = 'Done' WHERE id = $1;", id)
	return err
}