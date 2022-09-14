package sqlstore

import (
	"database/sql"

	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
)

type PatientRepository struct {
	store *Store
}

func (r *PatientRepository) Create(p *model.Patient) error {
	if err := p.Validate(); err != nil {
		return err
	}

	if err := p.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO patients (name, surname, birth_year, gender, phone, email, encrypted_password) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		&p.Name,
		&p.Surname,
		&p.Birth_year,
		&p.Gender,
		&p.Phone,
		&p.Email,
		&p.EncryptedPassword,
	).Scan(&p.ID)
}

func (r *PatientRepository) Find(id int) (*model.Patient, error) {
	p := &model.Patient{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, surname, birth_year, gender, phone, email, encrypted_password FROM patients WHERE id = $1", id,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Surname,
		&p.Birth_year,
		&p.Gender,
		&p.Phone,
		&p.Email,
		&p.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return p, nil
}

func (r *PatientRepository) FindByEmail(email string) (*model.Patient, error) {
	p := &model.Patient{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, surname, birth_year, gender, phone, email, encrypted_password FROM patients WHERE email = $1",
		email,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Surname,
		&p.Birth_year,
		&p.Gender,
		&p.Phone,
		&p.Email,
		&p.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return p, nil
}

func (r *PatientRepository) GetAll() ([]*model.Patient, error) {
	ps := []*model.Patient{}
	rows, err := r.store.db.Query("SELECT id, name, surname, gender, birth_year, phone, email, encrypted_password FROM patients")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := &model.Patient{}
		err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.Surname,
			&p.Gender,
			&p.Birth_year,
			&p.Phone,
			&p.Email,
			&p.EncryptedPassword,
		)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}