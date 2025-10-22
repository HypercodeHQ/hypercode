package models

type User struct {
	ID          int64
	Username    string
	Email       string
	DisplayName string
	Password    string
	CreatedAt   int64
	UpdatedAt   int64
}
