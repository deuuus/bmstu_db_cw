package sqlstore

import "github.com/deuuus/bmstu_db_cw/src/internal/app/model"

type RecordRepository struct {
	store *Store
}

func (r *RecordRepository) Create(rec *model.Record) error {
	return r.store.db.QueryRow(
		"INSERT INTO records (visit_id, disease_id, medicine_id) VALUES ($1, $2, $3) RETURNING id",
		&rec.Visit_id,
		&rec.Disease_id,
		&rec.Medicine_id,
	).Scan(&rec.ID)
}

func (r *RecordRepository) GetAllByPatient(patient_id int) ([]*model.Record, error) {
	recs := []*model.Record{}
	rows, err := r.store.db.Query("SELECT records.id, records.visit_id, records.disease_id, records.medicine_id FROM records JOIN visits ON records.visit_id = visits.id WHERE visits.patient_id = $1", patient_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rec := &model.Record{}
		err = rows.Scan(
			&rec.ID,
			&rec.Visit_id,
			&rec.Disease_id,
			&rec.Medicine_id,
		)
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}

func (r *RecordRepository) GetAllByDoctor(doctor_id int) ([]*model.Record, error) {
	recs := []*model.Record{}
	rows, err := r.store.db.Query("SELECT records.id, records.visit_id, records.disease_id, records.medicine_id FROM records JOIN visits ON records.visit_id = visits.id WHERE visits.doctor_id = $1", doctor_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rec := &model.Record{}
		err = rows.Scan(
			&rec.ID,
			&rec.Visit_id,
			&rec.Disease_id,
			&rec.Medicine_id,
		)
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}