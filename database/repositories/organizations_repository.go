package repositories

import (
	"database/sql"
	"errors"

	"github.com/hypercodehq/hypercode/database/models"
)

type OrganizationsRepository interface {
	Create(username, displayName string) (*models.Organization, error)
	FindByID(id int64) (*models.Organization, error)
	FindByUsername(username string) (*models.Organization, error)
	FindAll() ([]*models.Organization, error)
	Update(org *models.Organization) error
	Delete(id int64) error
}

type organizationsRepository struct {
	db *sql.DB
}

func NewOrganizationsRepository(db *sql.DB) OrganizationsRepository {
	return &organizationsRepository{db: db}
}

func (r *organizationsRepository) Create(username, displayName string) (*models.Organization, error) {
	query := `
		INSERT INTO organizations (username, display_name)
		VALUES (?, ?)
		RETURNING id, username, display_name, created_at, updated_at
	`

	org := &models.Organization{}
	err := r.db.QueryRow(query, username, displayName).Scan(
		&org.ID,
		&org.Username,
		&org.DisplayName,
		&org.CreatedAt,
		&org.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (r *organizationsRepository) FindByID(id int64) (*models.Organization, error) {
	query := `
		SELECT id, username, display_name, created_at, updated_at
		FROM organizations
		WHERE id = ?
	`

	org := &models.Organization{}
	err := r.db.QueryRow(query, id).Scan(
		&org.ID,
		&org.Username,
		&org.DisplayName,
		&org.CreatedAt,
		&org.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return org, nil
}

func (r *organizationsRepository) FindByUsername(username string) (*models.Organization, error) {
	query := `
		SELECT id, username, display_name, created_at, updated_at
		FROM organizations
		WHERE username = ?
	`

	org := &models.Organization{}
	err := r.db.QueryRow(query, username).Scan(
		&org.ID,
		&org.Username,
		&org.DisplayName,
		&org.CreatedAt,
		&org.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return org, nil
}

func (r *organizationsRepository) FindAll() ([]*models.Organization, error) {
	query := `
		SELECT id, username, display_name, created_at, updated_at
		FROM organizations
		ORDER BY username ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []*models.Organization
	for rows.Next() {
		org := &models.Organization{}
		err := rows.Scan(
			&org.ID,
			&org.Username,
			&org.DisplayName,
			&org.CreatedAt,
			&org.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		organizations = append(organizations, org)
	}

	return organizations, nil
}

func (r *organizationsRepository) Update(org *models.Organization) error {
	query := `
		UPDATE organizations
		SET username = ?, display_name = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, org.Username, org.DisplayName, org.ID)
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

func (r *organizationsRepository) Delete(id int64) error {
	query := `DELETE FROM organizations WHERE id = ?`

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
