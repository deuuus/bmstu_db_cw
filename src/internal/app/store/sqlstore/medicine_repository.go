package sqlstore

import (
	"database/sql"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
)

type MedicineRepository struct {
	store *Store
}

func (r *MedicineRepository) Find(id int) (*model.Medicine, error) {
	m := &model.Medicine{}
	if err := r.store.db.QueryRow(
		"SELECT name, cost FROM medicines WHERE id = $1", id,
	).Scan(
		&m.Name,
		&m.Cost,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return m, nil
}

func (r *MedicineRepository) GetAll() ([]*model.Medicine, error) {
	ms := []*model.Medicine{}
	rows, err := r.store.db.Query("SELECT id, name, cost FROM medicines")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Medicine{}
		err = rows.Scan(
			&m.ID,
			&m.Name,
			&m.Cost,
		)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}