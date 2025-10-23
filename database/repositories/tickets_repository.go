package repositories

import (
	"database/sql"
	"errors"

	"github.com/hypercodehq/hypercode/database/models"
)

type TicketsRepository interface {
	Create(repositoryID, authorID int64, title string, body *string) (*models.Ticket, error)
	FindByID(id int64) (*models.Ticket, error)
	FindByRepositoryAndNumber(repositoryID, number int64) (*models.Ticket, error)
	FindAllByRepository(repositoryID int64, status string) ([]*models.Ticket, error)
	CountByRepository(repositoryID int64, status string) (int64, error)
	Update(ticket *models.Ticket) error
	Close(ticketID, closedByID int64) error
	Reopen(ticketID int64) error
	Delete(id int64) error

	// Comments
	CreateComment(ticketID, authorID int64, body string) (*models.TicketComment, error)
	FindCommentsByTicket(ticketID int64) ([]*models.TicketComment, error)
	UpdateComment(comment *models.TicketComment) error
	DeleteComment(id int64) error
}

type ticketsRepository struct {
	db *sql.DB
}

func NewTicketsRepository(db *sql.DB) TicketsRepository {
	return &ticketsRepository{db: db}
}

func (r *ticketsRepository) Create(repositoryID, authorID int64, title string, body *string) (*models.Ticket, error) {
	// Get the next ticket number for this repository
	var number int64
	err := r.db.QueryRow(`
		SELECT COALESCE(MAX(number), 0) + 1
		FROM tickets
		WHERE repository_id = ?
	`, repositoryID).Scan(&number)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO tickets (repository_id, number, title, body, author_id)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, repository_id, number, title, body, status, author_id, closed_at, closed_by_id, created_at, updated_at
	`

	ticket := &models.Ticket{}
	err = r.db.QueryRow(query, repositoryID, number, title, body, authorID).Scan(
		&ticket.ID,
		&ticket.RepositoryID,
		&ticket.Number,
		&ticket.Title,
		&ticket.Body,
		&ticket.Status,
		&ticket.AuthorID,
		&ticket.ClosedAt,
		&ticket.ClosedByID,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (r *ticketsRepository) FindByID(id int64) (*models.Ticket, error) {
	query := `
		SELECT id, repository_id, number, title, body, status, author_id, closed_at, closed_by_id, created_at, updated_at
		FROM tickets
		WHERE id = ?
	`

	ticket := &models.Ticket{}
	err := r.db.QueryRow(query, id).Scan(
		&ticket.ID,
		&ticket.RepositoryID,
		&ticket.Number,
		&ticket.Title,
		&ticket.Body,
		&ticket.Status,
		&ticket.AuthorID,
		&ticket.ClosedAt,
		&ticket.ClosedByID,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return ticket, nil
}

func (r *ticketsRepository) FindByRepositoryAndNumber(repositoryID, number int64) (*models.Ticket, error) {
	query := `
		SELECT id, repository_id, number, title, body, status, author_id, closed_at, closed_by_id, created_at, updated_at
		FROM tickets
		WHERE repository_id = ? AND number = ?
	`

	ticket := &models.Ticket{}
	err := r.db.QueryRow(query, repositoryID, number).Scan(
		&ticket.ID,
		&ticket.RepositoryID,
		&ticket.Number,
		&ticket.Title,
		&ticket.Body,
		&ticket.Status,
		&ticket.AuthorID,
		&ticket.ClosedAt,
		&ticket.ClosedByID,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return ticket, nil
}

func (r *ticketsRepository) FindAllByRepository(repositoryID int64, status string) ([]*models.Ticket, error) {
	query := `
		SELECT id, repository_id, number, title, body, status, author_id, closed_at, closed_by_id, created_at, updated_at
		FROM tickets
		WHERE repository_id = ?
	`

	args := []interface{}{repositoryID}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY number DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*models.Ticket
	for rows.Next() {
		ticket := &models.Ticket{}
		err := rows.Scan(
			&ticket.ID,
			&ticket.RepositoryID,
			&ticket.Number,
			&ticket.Title,
			&ticket.Body,
			&ticket.Status,
			&ticket.AuthorID,
			&ticket.ClosedAt,
			&ticket.ClosedByID,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (r *ticketsRepository) CountByRepository(repositoryID int64, status string) (int64, error) {
	query := `SELECT COUNT(*) FROM tickets WHERE repository_id = ?`
	args := []interface{}{repositoryID}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (r *ticketsRepository) Update(ticket *models.Ticket) error {
	query := `
		UPDATE tickets
		SET title = ?, body = ?, updated_at = unixepoch()
		WHERE id = ?
	`

	_, err := r.db.Exec(query, ticket.Title, ticket.Body, ticket.ID)
	return err
}

func (r *ticketsRepository) Close(ticketID, closedByID int64) error {
	query := `
		UPDATE tickets
		SET status = 'closed', closed_at = unixepoch(), closed_by_id = ?, updated_at = unixepoch()
		WHERE id = ?
	`

	_, err := r.db.Exec(query, closedByID, ticketID)
	return err
}

func (r *ticketsRepository) Reopen(ticketID int64) error {
	query := `
		UPDATE tickets
		SET status = 'open', closed_at = NULL, closed_by_id = NULL, updated_at = unixepoch()
		WHERE id = ?
	`

	_, err := r.db.Exec(query, ticketID)
	return err
}

func (r *ticketsRepository) Delete(id int64) error {
	query := `DELETE FROM tickets WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// Comments

func (r *ticketsRepository) CreateComment(ticketID, authorID int64, body string) (*models.TicketComment, error) {
	query := `
		INSERT INTO ticket_comments (ticket_id, author_id, body)
		VALUES (?, ?, ?)
		RETURNING id, ticket_id, author_id, body, created_at, updated_at
	`

	comment := &models.TicketComment{}
	err := r.db.QueryRow(query, ticketID, authorID, body).Scan(
		&comment.ID,
		&comment.TicketID,
		&comment.AuthorID,
		&comment.Body,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Update ticket's updated_at timestamp
	_, _ = r.db.Exec(`UPDATE tickets SET updated_at = unixepoch() WHERE id = ?`, ticketID)

	return comment, nil
}

func (r *ticketsRepository) FindCommentsByTicket(ticketID int64) ([]*models.TicketComment, error) {
	query := `
		SELECT id, ticket_id, author_id, body, created_at, updated_at
		FROM ticket_comments
		WHERE ticket_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.TicketComment
	for rows.Next() {
		comment := &models.TicketComment{}
		err := rows.Scan(
			&comment.ID,
			&comment.TicketID,
			&comment.AuthorID,
			&comment.Body,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *ticketsRepository) UpdateComment(comment *models.TicketComment) error {
	query := `
		UPDATE ticket_comments
		SET body = ?, updated_at = unixepoch()
		WHERE id = ?
	`

	_, err := r.db.Exec(query, comment.Body, comment.ID)
	if err != nil {
		return err
	}

	// Update ticket's updated_at timestamp
	_, _ = r.db.Exec(`UPDATE tickets SET updated_at = unixepoch() WHERE id = ?`, comment.TicketID)

	return nil
}

func (r *ticketsRepository) DeleteComment(id int64) error {
	// Get ticket ID before deleting
	var ticketID int64
	err := r.db.QueryRow(`SELECT ticket_id FROM ticket_comments WHERE id = ?`, id).Scan(&ticketID)
	if err != nil {
		return err
	}

	query := `DELETE FROM ticket_comments WHERE id = ?`
	_, err = r.db.Exec(query, id)
	if err != nil {
		return err
	}

	// Update ticket's updated_at timestamp
	_, _ = r.db.Exec(`UPDATE tickets SET updated_at = unixepoch() WHERE id = ?`, ticketID)

	return nil
}
