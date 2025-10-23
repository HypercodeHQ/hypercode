package repositories

import (
	"database/sql"
	"errors"

	"github.com/hypercodehq/hypercode/database/models"
)

type DeviceAuthSessionsRepository interface {
	Create(id, code string, expiresAt int64) (*models.DeviceAuthSession, error)
	FindByID(id string) (*models.DeviceAuthSession, error)
	FindByCode(code string) (*models.DeviceAuthSession, error)
	Confirm(id string, userID int64, accessToken string) error
	UpdateStatus(id, status string) error
	DeleteExpired() error
}

type deviceAuthSessionsRepository struct {
	db *sql.DB
}

func NewDeviceAuthSessionsRepository(db *sql.DB) DeviceAuthSessionsRepository {
	return &deviceAuthSessionsRepository{db: db}
}

func (r *deviceAuthSessionsRepository) Create(id, code string, expiresAt int64) (*models.DeviceAuthSession, error) {
	query := `
		INSERT INTO device_auth_sessions (id, code, expires_at)
		VALUES (?, ?, ?)
		RETURNING id, code, user_id, access_token, status, created_at, expires_at
	`

	session := &models.DeviceAuthSession{}
	var userID sql.NullInt64
	var accessToken sql.NullString

	err := r.db.QueryRow(query, id, code, expiresAt).Scan(
		&session.ID,
		&session.Code,
		&userID,
		&accessToken,
		&session.Status,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	if userID.Valid {
		session.UserID = &userID.Int64
	}
	if accessToken.Valid {
		session.AccessToken = &accessToken.String
	}

	return session, nil
}

func (r *deviceAuthSessionsRepository) FindByID(id string) (*models.DeviceAuthSession, error) {
	query := `
		SELECT id, code, user_id, access_token, status, created_at, expires_at
		FROM device_auth_sessions
		WHERE id = ?
	`

	session := &models.DeviceAuthSession{}
	var userID sql.NullInt64
	var accessToken sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&session.ID,
		&session.Code,
		&userID,
		&accessToken,
		&session.Status,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if userID.Valid {
		session.UserID = &userID.Int64
	}
	if accessToken.Valid {
		session.AccessToken = &accessToken.String
	}

	return session, nil
}

func (r *deviceAuthSessionsRepository) FindByCode(code string) (*models.DeviceAuthSession, error) {
	query := `
		SELECT id, code, user_id, access_token, status, created_at, expires_at
		FROM device_auth_sessions
		WHERE code = ?
	`

	session := &models.DeviceAuthSession{}
	var userID sql.NullInt64
	var accessToken sql.NullString

	err := r.db.QueryRow(query, code).Scan(
		&session.ID,
		&session.Code,
		&userID,
		&accessToken,
		&session.Status,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if userID.Valid {
		session.UserID = &userID.Int64
	}
	if accessToken.Valid {
		session.AccessToken = &accessToken.String
	}

	return session, nil
}

func (r *deviceAuthSessionsRepository) Confirm(id string, userID int64, accessToken string) error {
	query := `
		UPDATE device_auth_sessions
		SET user_id = ?, access_token = ?, status = 'confirmed'
		WHERE id = ? AND status = 'pending'
	`

	result, err := r.db.Exec(query, userID, accessToken, id)
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

func (r *deviceAuthSessionsRepository) UpdateStatus(id, status string) error {
	query := `
		UPDATE device_auth_sessions
		SET status = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, status, id)
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

func (r *deviceAuthSessionsRepository) DeleteExpired() error {
	query := `
		DELETE FROM device_auth_sessions
		WHERE expires_at < unixepoch() OR status = 'expired'
	`

	_, err := r.db.Exec(query)
	return err
}
