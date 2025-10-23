package models

type Ticket struct {
	ID           int64
	RepositoryID int64
	Number       int64
	Title        string
	Body         *string
	Status       string // 'open' or 'closed'
	AuthorID     int64
	ClosedAt     *int64
	ClosedByID   *int64
	CreatedAt    int64
	UpdatedAt    int64
}

type TicketComment struct {
	ID        int64
	TicketID  int64
	AuthorID  int64
	Body      string
	CreatedAt int64
	UpdatedAt int64
}

type TicketLabel struct {
	ID           int64
	RepositoryID int64
	Name         string
	Color        string
	Description  *string
	CreatedAt    int64
}

type TicketLabelAssignment struct {
	ID        int64
	TicketID  int64
	LabelID   int64
	CreatedAt int64
}

type TicketAssignee struct {
	ID           int64
	TicketID     int64
	UserID       int64
	AssignedByID int64
	CreatedAt    int64
}

type TicketReaction struct {
	ID        int64
	TicketID  *int64
	CommentID *int64
	UserID    int64
	Emoji     string
	CreatedAt int64
}
