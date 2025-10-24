package repositories

import (
	"database/sql"
	"errors"

	"github.com/hypercommithq/hypercommit/database/models"
)

type AccessTokensRepository interface {
	Create(userID int64, name, tokenHash string) (*models.AccessToken, error)
	FindByID(id int64) (*models.AccessToken, error)
	FindByTokenHash(tokenHash string) (*models.AccessToken, error)
	FindByUserID(userID int64) ([]*models.AccessToken, error)
	UpdateLastUsed(id int64) error
	Delete(id int64) error
}

type accessTokensRepository struct {
	db *sql.DB
}

func NewAccessTokensRepository(db *sql.DB) AccessTokensRepository {
	return &accessTokensRepository{db: db}
}

func (r *accessTokensRepository) Create(userID int64, name, tokenHash string) (*models.AccessToken, error) {
	query := `
		INSERT INTO access_tokens (user_id, name, token_hash)
		VALUES (?, ?, ?)
		RETURNING id, user_id, name, token_hash, last_used_at, created_at
	`

	token := &models.AccessToken{}
	var lastUsedAt sql.NullInt64
	err := r.db.QueryRow(query, userID, name, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.Name,
		&token.TokenHash,
		&lastUsedAt,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Convert sql.NullInt64 to *int64
	if lastUsedAt.Valid {
		token.LastUsedAt = &lastUsedAt.Int64
	}

	return token, nil
}

func (r *accessTokensRepository) FindByID(id int64) (*models.AccessToken, error) {
	query := `
		SELECT id, user_id, name, token_hash, last_used_at, created_at
		FROM access_tokens
		WHERE id = ?
	`

	token := &models.AccessToken{}
	var lastUsedAt sql.NullInt64
	err := r.db.QueryRow(query, id).Scan(
		&token.ID,
		&token.UserID,
		&token.Name,
		&token.TokenHash,
		&lastUsedAt,
		&token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Convert sql.NullInt64 to *int64
	if lastUsedAt.Valid {
		token.LastUsedAt = &lastUsedAt.Int64
	}

	return token, nil
}

func (r *accessTokensRepository) FindByTokenHash(tokenHash string) (*models.AccessToken, error) {
	query := `
		SELECT id, user_id, name, token_hash, last_used_at, created_at
		FROM access_tokens
		WHERE token_hash = ?
	`

	token := &models.AccessToken{}
	var lastUsedAt sql.NullInt64
	err := r.db.QueryRow(query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.Name,
		&token.TokenHash,
		&lastUsedAt,
		&token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Convert sql.NullInt64 to *int64
	if lastUsedAt.Valid {
		token.LastUsedAt = &lastUsedAt.Int64
	}

	return token, nil
}

func (r *accessTokensRepository) FindByUserID(userID int64) ([]*models.AccessToken, error) {
	query := `
		SELECT id, user_id, name, token_hash, last_used_at, created_at
		FROM access_tokens
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*models.AccessToken
	for rows.Next() {
		token := &models.AccessToken{}
		var lastUsedAt sql.NullInt64
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Name,
			&token.TokenHash,
			&lastUsedAt,
			&token.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert sql.NullInt64 to *int64
		if lastUsedAt.Valid {
			token.LastUsedAt = &lastUsedAt.Int64
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (r *accessTokensRepository) UpdateLastUsed(id int64) error {
	query := `
		UPDATE access_tokens
		SET last_used_at = unixepoch()
		WHERE id = ?
	`

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

func (r *accessTokensRepository) Delete(id int64) error {
	query := `DELETE FROM access_tokens WHERE id = ?`

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
