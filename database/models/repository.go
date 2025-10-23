package models

type Repository struct {
	ID            int64
	Name          string
	Description   *string
	DefaultBranch string
	Visibility    string
	OwnerUserID   *int64
	OwnerOrgID    *int64
	CreatedAt     int64
	UpdatedAt     int64
}
