package models

type Contributor struct {
	ID           int64
	RepositoryID int64
	UserID       int64
	Role         string
	CreatedAt    int64
}
