package repositories

import (
	"database/sql"
	"errors"

	"github.com/hypercodehq/hypercode/database/models"
)

type StarsRepository interface {
	Create(repositoryID, userID int64) (*models.Star, error)
	Delete(repositoryID, userID int64) error
	FindByUserAndRepository(repositoryID, userID int64) (*models.Star, error)
	CountByRepository(repositoryID int64) (int64, error)
	FindStarredRepositoriesByUser(userID int64) ([]*models.Repository, error)
}

type starsRepository struct {
	db *sql.DB
}

func NewStarsRepository(db *sql.DB) StarsRepository {
	return &starsRepository{db: db}
}

func (r *starsRepository) Create(repositoryID, userID int64) (*models.Star, error) {
	query := `
		INSERT INTO stars (repository_id, user_id)
		VALUES (?, ?)
		RETURNING id, repository_id, user_id, created_at
	`

	star := &models.Star{}
	err := r.db.QueryRow(query, repositoryID, userID).Scan(
		&star.ID,
		&star.RepositoryID,
		&star.UserID,
		&star.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return star, nil
}

func (r *starsRepository) Delete(repositoryID, userID int64) error {
	query := `DELETE FROM stars WHERE repository_id = ? AND user_id = ?`

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

func (r *starsRepository) FindByUserAndRepository(repositoryID, userID int64) (*models.Star, error) {
	query := `
		SELECT id, repository_id, user_id, created_at
		FROM stars
		WHERE repository_id = ? AND user_id = ?
	`

	star := &models.Star{}
	err := r.db.QueryRow(query, repositoryID, userID).Scan(
		&star.ID,
		&star.RepositoryID,
		&star.UserID,
		&star.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return star, nil
}

func (r *starsRepository) CountByRepository(repositoryID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM stars WHERE repository_id = ?`

	var count int64
	err := r.db.QueryRow(query, repositoryID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *starsRepository) FindStarredRepositoriesByUser(userID int64) ([]*models.Repository, error) {
	query := `
		SELECT r.id, r.name, r.description, r.default_branch, r.visibility, r.owner_user_id, r.owner_org_id, r.created_at, r.updated_at
		FROM repositories r
		INNER JOIN stars s ON s.repository_id = r.id
		WHERE s.user_id = ?
		ORDER BY s.created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repositories []*models.Repository
	for rows.Next() {
		repo := &models.Repository{}
		err := rows.Scan(
			&repo.ID,
			&repo.Name,
			&repo.Description,
			&repo.DefaultBranch,
			&repo.Visibility,
			&repo.OwnerUserID,
			&repo.OwnerOrgID,
			&repo.CreatedAt,
			&repo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, repo)
	}

	return repositories, nil
}
