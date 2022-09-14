package sqlstore

import (
	"database/sql"

	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
)

type AdminRepository struct {
	store *Store
}

func (r *AdminRepository) Create(a *model.Admin) error {
	if err := a.Validate(); err != nil {
		return err
	}

	if err := a.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO admins (name, surname, email, encrypted_password) VALUES ($1, $2, $3, $4) RETURNING id",
		&a.Name,
		&a.Surname,
		&a.Email,
		&a.EncryptedPassword,
	).Scan(&a.ID)
}

func (r *AdminRepository) Find(id int) (*model.Admin, error) {
	a := &model.Admin{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, surname, email, encrypted_password FROM admins WHERE id = $1", id,
	).Scan(
		&a.ID,
		&a.Name,
		&a.Surname,
		&a.Email,
		&a.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return a, nil
}

func (r *AdminRepository) FindByEmail(email string) (*model.Admin, error) {
	a := &model.Admin{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, surname, email, encrypted_password FROM admins WHERE email = $1",
		email,
	).Scan(
		&a.ID,
		&a.Name,
		&a.Surname,
		&a.Email,
		&a.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return a, nil
}

func (r *AdminRepository) GetAll() ([]*model.Admin, error) {
	admins := []*model.Admin{}
	rows, err := r.store.db.Query("SELECT id, name, surname, email, encrypted_password FROM admins")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.Admin{}
		err = rows.Scan(
			&a.ID,
			&a.Name,
			&a.Surname,
			&a.Email,
			&a.EncryptedPassword,
		)
		if err != nil {
			return nil, err
		}
		admins = append(admins, a)
	}
	return admins, nil
}