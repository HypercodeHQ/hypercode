package models

type DeviceAuthSession struct {
	ID          string // UUID
	Code        string // User-friendly code like "ABCD-1234"
	UserID      *int64
	AccessToken *string
	Status      string // pending, confirmed, expired
	CreatedAt   int64
	ExpiresAt   int64
}
