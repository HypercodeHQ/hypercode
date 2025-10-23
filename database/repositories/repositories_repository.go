package repositories

import (
	"database/sql"
	"errors"

	"github.com/hypercodehq/hypercode/database/models"
)

type RepositoriesRepository interface {
	CreateForUser(userID int64, name, visibility, defaultBranch string, description *string) (*models.Repository, error)
	CreateForOrg(orgID int64, name, visibility, defaultBranch string, description *string) (*models.Repository, error)
	FindByID(id int64) (*models.Repository, error)
	FindByUserAndName(userID int64, name string) (*models.Repository, error)
	FindByOrgAndName(orgID int64, name string) (*models.Repository, error)
	FindByOwnerAndName(ownerUsername, repoName string) (*models.Repository, error)
	FindAllByUser(userID int64) ([]*models.Repository, error)
	FindAllByOrg(orgID int64) ([]*models.Repository, error)
	FindPublic() ([]*models.Repository, error)
	FindAll() ([]*models.Repository, error)
	Update(repo *models.Repository) error
	Delete(id int64) error
}

type repositoriesRepository struct {
	db *sql.DB
}

func NewRepositoriesRepository(db *sql.DB) RepositoriesRepository {
	return &repositoriesRepository{db: db}
}

func (r *repositoriesRepository) CreateForUser(userID int64, name, visibility, defaultBranch string, description *string) (*models.Repository, error) {
	query := `
		INSERT INTO repositories (name, description, default_branch, visibility, owner_user_id)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, name, description, default_branch, visibility, owner_user_id, owner_org_id, created_at, updated_at
	`

	repo := &models.Repository{}
	err := r.db.QueryRow(query, name, description, defaultBranch, visibility, userID).Scan(
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

	return repo, nil
}

func (r *repositoriesRepository) CreateForOrg(orgID int64, name, visibility, defaultBranch string, description *string) (*models.Repository, error) {
	query := `
		INSERT INTO repositories (name, description, default_branch, visibility, owner_org_id)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, name, description, default_branch, visibility, owner_user_id, owner_org_id, created_at, updated_at
	`

	repo := &models.Repository{}
	err := r.db.QueryRow(query, name, description, defaultBranch, visibility, orgID).Scan(
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

	return repo, nil
}

func (r *repositoriesRepository) FindByID(id int64) (*models.Repository, error) {
	query := `
		SELECT id, name, description, default_branch, visibility, owner_user_id, owner_org_id, created_at, updated_at
		FROM repositories
		WHERE id = ?
	`

	repo := &models.Repository{}
	err := r.db.QueryRow(query, id).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return repo, nil
}

func (r *repositoriesRepository) FindByUserAndName(userID int64, name string) (*models.Repository, error) {
	query := `
		SELECT id, name, description, default_branch, visibility, owner_user_id, owner_org_id, created_at, updated_at
		FROM repositories
		WHERE owner_user_id = ? AND name = ?
	`

	repo := &models.Repository{}
	err := r.db.QueryRow(query, userID, name).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return repo, nil
}

func (r *repositoriesRepository) FindByOrgAndName(orgID int64, name string) (*models.Repository, error) {
	query := `
		SELECT id, name, description, default_branch, visibility, owner_user_id, owner_org_id, created_at, updated_at
		FROM repositories
		WHERE owner_org_id = ? AND name = ?
	`

	repo := &models.Repository{}
	err := r.db.QueryRow(query, orgID, name).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return repo, nil
}

func (r *repositoriesRepository) FindByOwnerAndName(ownerUsername, repoName string) (*models.Repository, error) {
	query := `
		SELECT r.id, r.name, r.description, r.default_branch, r.visibility, r.owner_user_id, r.owner_org_id, r.created_at, r.updated_at
		FROM repositories r
		LEFT JOIN users u ON r.owner_user_id = u.id
		LEFT JOIN organizations o ON r.owner_org_id = o.id
		WHERE (u.username = ? OR o.username = ?) AND r.name = ?
	`

	repo := &models.Repository{}
	err := r.db.QueryRow(query, ownerUsername, ownerUsername, repoName).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return repo, nil
}

func (r *repositoriesRepository) FindAllByUser(userID int64) ([]*models.Repository, error) {
	query := `
		SELECT id, name, description, default_branch, visibility, owner_user_id, owner_org_id, created_at, updated_at
		FROM repositories
		WHERE owner_user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*models.Repository
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
		repos = append(repos, repo)
	}

	return repos, nil
}

func (r *repositoriesRepository) FindAllByOrg(orgID int64) ([]*models.Repository, error) {
	query := `
		SELECT id, name, description, default_branch, visibility, owner_user_id, owner_org_id, created_at, updated_at
		FROM repositories
		WHERE owner_org_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*models.Repository
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
		repos = append(repos, repo)
	}

	return repos, nil
}

func (r *repositoriesRepository) FindPublic() ([]*models.Repository, error) {
	query := `
		SELECT id, name, description, default_branch, visibility, owner_user_id, owner_org_id, created_at, updated_at
		FROM repositories
		WHERE visibility = 'public'
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*models.Repository
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
		repos = append(repos, repo)
	}

	return repos, nil
}

func (r *repositoriesRepository) FindAll() ([]*models.Repository, error) {
	query := `
		SELECT id, name, description, default_branch, visibility, owner_user_id, owner_org_id, created_at, updated_at
		FROM repositories
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*models.Repository
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
		repos = append(repos, repo)
	}

	return repos, nil
}

func (r *repositoriesRepository) Update(repo *models.Repository) error {
	query := `
		UPDATE repositories
		SET name = ?, description = ?, default_branch = ?, visibility = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, repo.Name, repo.Description, repo.DefaultBranch, repo.Visibility, repo.ID)
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

func (r *repositoriesRepository) Delete(id int64) error {
	query := `DELETE FROM repositories WHERE id = ?`

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
