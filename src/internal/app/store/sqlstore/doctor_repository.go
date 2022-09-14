package sqlstore

import (
	"database/sql"

	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
)

type DoctorRepository struct {
	store *Store
}

func (r *DoctorRepository) Create(d *model.Doctor) error {
	if err := d.Validate(); err != nil {
		return err
	}

	if err := d.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO doctors (name, surname, work_since, spec_id, email, encrypted_password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		&d.Name,
		&d.Surname,
		&d.Work_since,
		&d.Spec_id,
		&d.Email,
		&d.EncryptedPassword,
	).Scan(&d.ID)
}

func (r *DoctorRepository) Find(id int) (*model.Doctor, error) {
	d := &model.Doctor{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, surname, work_since, spec_id, email, encrypted_password FROM doctors WHERE id = $1", id,
	).Scan(
		&d.ID,
		&d.Name,
		&d.Surname,
		&d.Work_since,
		&d.Spec_id,
		&d.Email,
		&d.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return d, nil
}

func (r *DoctorRepository) FindByEmail(email string) (*model.Doctor, error) {
	d := &model.Doctor{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, surname, work_since, spec_id, email, encrypted_password FROM doctors WHERE email = $1",
		email,
	).Scan(
		&d.ID,
		&d.Name,
		&d.Surname,
		&d.Work_since,
		&d.Spec_id,
		&d.Email,
		&d.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return d, nil
}

func (r *DoctorRepository) GetAll() ([]*model.Doctor, error) {
	docs := []*model.Doctor{}
	rows, err := r.store.db.Query("SELECT id, name, surname, work_since, spec_id, email, encrypted_password FROM doctors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		d := &model.Doctor{}
		err = rows.Scan(
			&d.ID,
			&d.Name,
			&d.Surname,
			&d.Work_since,
			&d.Spec_id,
			&d.Email,
			&d.EncryptedPassword,
		)
		if err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}