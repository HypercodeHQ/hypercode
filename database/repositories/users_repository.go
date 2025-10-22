package repositories

import (
	"database/sql"
	"errors"

	"github.com/hyperstitieux/hypercode/database/models"
)

type UsersRepository interface {
	Create(username, email, displayName, password string) (*models.User, error)
	FindByID(id int64) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int64) error
}

type usersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) UsersRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) Create(username, email, displayName, password string) (*models.User, error) {
	query := `
		INSERT INTO users (username, email, display_name, password)
		VALUES (?, ?, ?, ?)
		RETURNING id, username, email, display_name, password, created_at, updated_at
	`

	user := &models.User{}
	err := r.db.QueryRow(query, username, email, displayName, password).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *usersRepository) FindByID(id int64) (*models.User, error) {
	query := `
		SELECT id, username, email, display_name, password, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *usersRepository) FindByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, display_name, password, created_at, updated_at
		FROM users
		WHERE username = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *usersRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, display_name, password, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *usersRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET username = ?, email = ?, display_name = ?, password = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, user.Username, user.Email, user.DisplayName, user.Password, user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *usersRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
