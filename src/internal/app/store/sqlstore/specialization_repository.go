package sqlstore

import (
	"database/sql"

	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
)

type SpecializationRepository struct {
	store *Store
}

func (r *SpecializationRepository) Find(id int) (*model.Specialization, error) {
	s := &model.Specialization{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, salary FROM specializations WHERE id = $1", id,
	).Scan(
		&s.ID,
		&s.Name,
		&s.Salary,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return s, nil
}

func (r *SpecializationRepository) FindByName(name string) (*model.Specialization, error) {
	s := &model.Specialization{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, salary FROM specializations WHERE name = $1",
		name,
	).Scan(
		&s.ID,
		&s.Name,
		&s.Salary,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return s, nil
}

func (r *SpecializationRepository) GetAll() ([]*model.Specialization, error) {
	specs := []*model.Specialization{}
	rows, err := r.store.db.Query("SELECT id, name, salary FROM specializations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		spec := &model.Specialization{}
		err = rows.Scan(
			&spec.ID,
			&spec.Name,
			&spec.Salary,
		)
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
	}
	return specs, nil
}