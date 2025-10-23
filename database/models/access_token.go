package models

type AccessToken struct {
	ID         int64
	UserID     int64
	Name       string
	TokenHash  string
	LastUsedAt *int64
	CreatedAt  int64
}
