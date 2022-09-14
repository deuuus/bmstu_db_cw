package sqlstore

import (
	"database/sql"

	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
)

type DiseaseRepository struct {
	store *Store
}

func (r *DiseaseRepository) Percent(spec_id int) ([]int, []float64, error) {
	rows, err := r.store.db.Query("SELECT disease_id, percent FROM percent($1)", spec_id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	ds := make([]int, 0)
	ps := make([]float64, 0)
	for rows.Next() {
		var d int
		var p float64
		err := rows.Scan(&d, &p)
		if err != nil {
			return nil, nil, err
		}
		ds = append(ds, d)
		ps = append(ps, p)
	}
	return ds, ps, nil
}

func (r *DiseaseRepository) Find(id int) (*model.Disease, error) {
	d := &model.Disease{}
	if err := r.store.db.QueryRow(
		"SELECT name, spec_id FROM diseases WHERE id = $1", id,
	).Scan(
		&d.Name,
		&d.Spec_id,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return d, nil
}

func (r *DiseaseRepository) GetAll() ([]*model.Disease, error) {
	ds := []*model.Disease{}
	rows, err := r.store.db.Query("SELECT id, name, spec_id FROM diseases")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		d := &model.Disease{}
		err = rows.Scan(
			&d.ID,
			&d.Name,
			&d.Spec_id,
		)
		if err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func (r *DiseaseRepository) GetBySpecialization(spec_id int) ([]*model.Disease, error) {
	ds := []*model.Disease{}
	rows, err := r.store.db.Query("SELECT id, name FROM diseases WHERE spec_id = $1", spec_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		d := &model.Disease{}
		err = rows.Scan(
			&d.ID,
			&d.Name,
		)
		if err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}