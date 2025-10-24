package repositories

import (
	"database/sql"
	"errors"

	"github.com/hypercommithq/hypercommit/database/models"
)

type ContributorsRepository interface {
	Create(repositoryID, userID int64, role string) (*models.Contributor, error)
	FindByID(id int64) (*models.Contributor, error)
	FindByRepositoryAndUser(repositoryID, userID int64) (*models.Contributor, error)
	FindAllByRepository(repositoryID int64) ([]*models.Contributor, error)
	FindAllByUser(userID int64) ([]*models.Contributor, error)
	UpdateRole(id int64, role string) error
	Delete(id int64) error
	DeleteByRepositoryAndUser(repositoryID, userID int64) error
}

type contributorsRepository struct {
	db *sql.DB
}

func NewContributorsRepository(db *sql.DB) ContributorsRepository {
	return &contributorsRepository{db: db}
}

func (r *contributorsRepository) Create(repositoryID, userID int64, role string) (*models.Contributor, error) {
	query := `
		INSERT INTO contributors (repository_id, user_id, role)
		VALUES (?, ?, ?)
		RETURNING id, repository_id, user_id, role, created_at
	`

	contributor := &models.Contributor{}
	err := r.db.QueryRow(query, repositoryID, userID, role).Scan(
		&contributor.ID,
		&contributor.RepositoryID,
		&contributor.UserID,
		&contributor.Role,
		&contributor.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return contributor, nil
}

func (r *contributorsRepository) FindByID(id int64) (*models.Contributor, error) {
	query := `
		SELECT id, repository_id, user_id, role, created_at
		FROM contributors
		WHERE id = ?
	`

	contributor := &models.Contributor{}
	err := r.db.QueryRow(query, id).Scan(
		&contributor.ID,
		&contributor.RepositoryID,
		&contributor.UserID,
		&contributor.Role,
		&contributor.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return contributor, nil
}

func (r *contributorsRepository) FindByRepositoryAndUser(repositoryID, userID int64) (*models.Contributor, error) {
	query := `
		SELECT id, repository_id, user_id, role, created_at
		FROM contributors
		WHERE repository_id = ? AND user_id = ?
	`

	contributor := &models.Contributor{}
	err := r.db.QueryRow(query, repositoryID, userID).Scan(
		&contributor.ID,
		&contributor.RepositoryID,
		&contributor.UserID,
		&contributor.Role,
		&contributor.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return contributor, nil
}

func (r *contributorsRepository) FindAllByRepository(repositoryID int64) ([]*models.Contributor, error) {
	query := `
		SELECT id, repository_id, user_id, role, created_at
		FROM contributors
		WHERE repository_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, repositoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contributors []*models.Contributor
	for rows.Next() {
		contributor := &models.Contributor{}
		err := rows.Scan(
			&contributor.ID,
			&contributor.RepositoryID,
			&contributor.UserID,
			&contributor.Role,
			&contributor.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		contributors = append(contributors, contributor)
	}

	return contributors, nil
}

func (r *contributorsRepository) FindAllByUser(userID int64) ([]*models.Contributor, error) {
	query := `
		SELECT id, repository_id, user_id, role, created_at
		FROM contributors
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contributors []*models.Contributor
	for rows.Next() {
		contributor := &models.Contributor{}
		err := rows.Scan(
			&contributor.ID,
			&contributor.RepositoryID,
			&contributor.UserID,
			&contributor.Role,
			&contributor.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		contributors = append(contributors, contributor)
	}

	return contributors, nil
}

func (r *contributorsRepository) UpdateRole(id int64, role string) error {
	query := `
		UPDATE contributors
		SET role = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, role, id)
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

func (r *contributorsRepository) Delete(id int64) error {
	query := `DELETE FROM contributors WHERE id = ?`

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

func (r *contributorsRepository) DeleteByRepositoryAndUser(repositoryID, userID int64) error {
	query := `DELETE FROM contributors WHERE repository_id = ? AND user_id = ?`

	result, err := r.db.Exec(query, repositoryID, userID)
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
