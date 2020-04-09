package sqlstore

import (
	"database/sql"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store"
)

// User rerpository
type UserRepository struct {
	store *Store
}

// Create new user
func (repo *UserRepository) Create(newUser *models.User) error {
	if err := newUser.Validate(); err != nil {
		return err
	}

	if err := newUser.BeforeCreate(); err != nil {
		return err
	}

	return repo.store.db.QueryRow(
		"insert into users (username, email, encrypted_password, created_at, token, contacts, role, is_active) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id",
		newUser.Login,
		newUser.Email,
		newUser.EncryptedPassword,
		newUser.CreatedAt,
		newUser.Token,
		newUser.Contacts,
		newUser.Role,
		newUser.IsActive,
	).Scan(&newUser.ID)
}

// Find user by email
func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}

	if err := repo.store.db.QueryRow(
		"select * from users where email = $1",
		email,
	).Scan(&user.ID, &user.Login, &user.Email, &user.EncryptedPassword, &user.CreatedAt,
		&user.Token, &user.Contacts, &user.Role, &user.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}

// Find user by id
func (repo *UserRepository) FindById(id int) (*models.User, error) {
	user := &models.User{}

	if err := repo.store.db.QueryRow(
		"select * from users where id = $1",
		id,
	).Scan(&user.ID, &user.Login, &user.Email, &user.EncryptedPassword, &user.CreatedAt,
		&user.Token, &user.Contacts, &user.Role, &user.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}

// Update user info by new model
func (repo *UserRepository) Update(user *models.User) error {
	_, err := repo.store.db.Exec(
		`update users set 
		username = $2, email = $3, encrypted_password = $4, created_at = $5, 
		token = $6, contacts = $7, role = $8, is_active = $9 
		where id = $1`,
		user.ID,
		user.Login,
		user.Email,
		user.EncryptedPassword,
		user.CreatedAt,
		user.Token,
		user.Contacts,
		user.Role,
		user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}

		return err
	}

	return nil
}
