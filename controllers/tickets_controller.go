package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/hypercommit/database/repositories"
	"github.com/hypercommithq/hypercommit/httperror"
	custommiddleware "github.com/hypercommithq/hypercommit/middleware"
	"github.com/hypercommithq/hypercommit/services"
	"github.com/hypercommithq/hypercommit/views/pages"
)

type TicketsController interface {
	List(w http.ResponseWriter, r *http.Request) error
	Show(w http.ResponseWriter, r *http.Request) error
	New(w http.ResponseWriter, r *http.Request) error
	Create(w http.ResponseWriter, r *http.Request) error
	Close(w http.ResponseWriter, r *http.Request) error
	Reopen(w http.ResponseWriter, r *http.Request) error
	CreateComment(w http.ResponseWriter, r *http.Request) error
}

type ticketsController struct {
	tickets      repositories.TicketsRepository
	repos        repositories.RepositoriesRepository
	users        repositories.UsersRepository
	stars        repositories.StarsRepository
	contributors repositories.ContributorsRepository
	authService  services.AuthService
}

func NewTicketsController(
	tickets repositories.TicketsRepository,
	repos repositories.RepositoriesRepository,
	users repositories.UsersRepository,
	stars repositories.StarsRepository,
	contributors repositories.ContributorsRepository,
	authService services.AuthService,
) TicketsController {
	return &ticketsController{
		tickets:      tickets,
		repos:        repos,
		users:        users,
		stars:        stars,
		contributors: contributors,
		authService:  authService,
	}
}

func (c *ticketsController) List(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.New(500, "failed to find repository")
	}
	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	currentUser := custommiddleware.GetUserFromContext(r)

	// Get filter from query params (default to "open")
	statusFilter := r.URL.Query().Get("status")
	if statusFilter == "" {
		statusFilter = "open"
	}

	tickets, err := c.tickets.FindAllByRepository(repo.ID, statusFilter)
	if err != nil {
		slog.Error("failed to fetch tickets", "error", err)
		tickets = []*models.Ticket{}
	}

	openCount, _ := c.tickets.CountByRepository(repo.ID, "open")
	closedCount, _ := c.tickets.CountByRepository(repo.ID, "closed")

	// Check permissions
	canManage := false
	if currentUser != nil {
		if repo.OwnerUserID != nil && *repo.OwnerUserID == currentUser.ID {
			canManage = true
		} else {
			// Check if user is an admin contributor
			contributor, err := c.contributors.FindByRepositoryAndUser(repo.ID, currentUser.ID)
			if err == nil && contributor != nil && contributor.Role == "admin" {
				canManage = true
			}
		}
	}

	// Get star info
	starCount, _ := c.stars.CountByRepository(repo.ID)
	hasStarred := false
	if currentUser != nil {
		star, _ := c.stars.FindByUserAndRepository(repo.ID, currentUser.ID)
		hasStarred = star != nil
	}

	cloneURL := "https://" + r.Host + "/" + owner + "/" + repoName
	repositoryURL := cloneURL

	return pages.TicketsList(r, &pages.TicketsListData{
		User:          currentUser,
		Repository:    repo,
		OwnerUsername: owner,
		Tickets:       tickets,
		StatusFilter:  statusFilter,
		OpenCount:     openCount,
		ClosedCount:   closedCount,
		CanManage:     canManage,
		StarCount:     starCount,
		HasStarred:    hasStarred,
		CloneURL:      cloneURL,
		RepositoryURL: repositoryURL,
	}).Render(w, r)
}

func (c *ticketsController) Show(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")
	numberStr := chi.URLParam(r, "number")

	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		return httperror.BadRequest("invalid ticket number")
	}

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.New(500, "failed to find repository")
	}
	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	ticket, err := c.tickets.FindByRepositoryAndNumber(repo.ID, number)
	if err != nil {
		return httperror.New(500, "failed to find ticket")
	}
	if ticket == nil {
		return httperror.NotFound("ticket not found")
	}

	// Get ticket author
	author, err := c.users.FindByID(ticket.AuthorID)
	if err != nil {
		slog.Error("failed to fetch ticket author", "error", err)
	}

	// Get comments
	comments, err := c.tickets.FindCommentsByTicket(ticket.ID)
	if err != nil {
		slog.Error("failed to fetch comments", "error", err)
		comments = []*models.TicketComment{}
	}

	// Get comment authors
	commentAuthors := make(map[int64]*models.User)
	for _, comment := range comments {
		if _, exists := commentAuthors[comment.AuthorID]; !exists {
			user, err := c.users.FindByID(comment.AuthorID)
			if err == nil && user != nil {
				commentAuthors[comment.AuthorID] = user
			}
		}
	}

	currentUser := custommiddleware.GetUserFromContext(r)

	// Check permissions
	canManage := false
	if currentUser != nil {
		if repo.OwnerUserID != nil && *repo.OwnerUserID == currentUser.ID {
			canManage = true
		} else {
			// Check if user is an admin contributor
			contributor, err := c.contributors.FindByRepositoryAndUser(repo.ID, currentUser.ID)
			if err == nil && contributor != nil && contributor.Role == "admin" {
				canManage = true
			}
		}
	}

	// Get star info
	starCount, _ := c.stars.CountByRepository(repo.ID)
	hasStarred := false
	if currentUser != nil {
		star, _ := c.stars.FindByUserAndRepository(repo.ID, currentUser.ID)
		hasStarred = star != nil
	}

	cloneURL := "https://" + r.Host + "/" + owner + "/" + repoName
	repositoryURL := cloneURL

	return pages.ShowTicket(r, &pages.ShowTicketData{
		User:           currentUser,
		Repository:     repo,
		OwnerUsername:  owner,
		Ticket:         ticket,
		Author:         author,
		Comments:       comments,
		CommentAuthors: commentAuthors,
		CanManage:      canManage,
		StarCount:      starCount,
		HasStarred:     hasStarred,
		CloneURL:       cloneURL,
		RepositoryURL:  repositoryURL,
	}).Render(w, r)
}

func (c *ticketsController) New(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.New(500, "failed to find repository")
	}
	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	currentUser := custommiddleware.GetUserFromContext(r)
	if currentUser == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	// Check permissions
	canManage := false
	if repo.OwnerUserID != nil && *repo.OwnerUserID == currentUser.ID {
		canManage = true
	} else {
		// Check if user is an admin contributor
		contributor, err := c.contributors.FindByRepositoryAndUser(repo.ID, currentUser.ID)
		if err == nil && contributor != nil && contributor.Role == "admin" {
			canManage = true
		}
	}

	// Get star info
	starCount, _ := c.stars.CountByRepository(repo.ID)
	star, _ := c.stars.FindByUserAndRepository(repo.ID, currentUser.ID)
	hasStarred := star != nil

	cloneURL := "https://" + r.Host + "/" + owner + "/" + repoName
	repositoryURL := cloneURL

	return pages.NewTicket(r, &pages.NewTicketData{
		User:          currentUser,
		Repository:    repo,
		OwnerUsername: owner,
		CanManage:     canManage,
		StarCount:     starCount,
		HasStarred:    hasStarred,
		CloneURL:      cloneURL,
		RepositoryURL: repositoryURL,
	}).Render(w, r)
}

func (c *ticketsController) Create(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.New(500, "failed to find repository")
	}
	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	currentUser := custommiddleware.GetUserFromContext(r)
	if currentUser == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return httperror.BadRequest("invalid form data")
	}

	title := r.FormValue("title")
	body := r.FormValue("body")

	var bodyPtr *string
	if body != "" {
		bodyPtr = &body
	}

	// Validate
	if title == "" {
		// Return to form with error
		canManage := false
		if repo.OwnerUserID != nil && *repo.OwnerUserID == currentUser.ID {
			canManage = true
		} else {
			// Check if user is an admin contributor
			contributor, err := c.contributors.FindByRepositoryAndUser(repo.ID, currentUser.ID)
			if err == nil && contributor != nil && contributor.Role == "admin" {
				canManage = true
			}
		}
		starCount, _ := c.stars.CountByRepository(repo.ID)
		star, _ := c.stars.FindByUserAndRepository(repo.ID, currentUser.ID)
		hasStarred := star != nil
		cloneURL := "https://" + r.Host + "/" + owner + "/" + repoName
		repositoryURL := cloneURL

		return pages.NewTicket(r, &pages.NewTicketData{
			User:          currentUser,
			Repository:    repo,
			OwnerUsername: owner,
			Title:         title,
			Body:          body,
			TitleError:    "Title is required",
			CanManage:     canManage,
			StarCount:     starCount,
			HasStarred:    hasStarred,
			CloneURL:      cloneURL,
			RepositoryURL: repositoryURL,
		}).Render(w, r)
	}

	ticket, err := c.tickets.Create(repo.ID, currentUser.ID, title, bodyPtr)
	if err != nil {
		slog.Error("failed to create ticket", "error", err)
		return httperror.New(500, "failed to create ticket")
	}

	http.Redirect(w, r, "/"+owner+"/"+repoName+"/tickets/"+strconv.FormatInt(ticket.Number, 10), http.StatusSeeOther)
	return nil
}

func (c *ticketsController) Close(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")
	numberStr := chi.URLParam(r, "number")

	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		return httperror.BadRequest("invalid ticket number")
	}

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.New(500, "failed to find repository")
	}
	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	currentUser := custommiddleware.GetUserFromContext(r)
	if currentUser == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	ticket, err := c.tickets.FindByRepositoryAndNumber(repo.ID, number)
	if err != nil || ticket == nil {
		return httperror.NotFound("ticket not found")
	}

	err = c.tickets.Close(ticket.ID, currentUser.ID)
	if err != nil {
		slog.Error("failed to close ticket", "error", err)
		return httperror.New(500, "failed to close ticket")
	}

	http.Redirect(w, r, "/"+owner+"/"+repoName+"/tickets/"+numberStr, http.StatusSeeOther)
	return nil
}

func (c *ticketsController) Reopen(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")
	numberStr := chi.URLParam(r, "number")

	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		return httperror.BadRequest("invalid ticket number")
	}

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.New(500, "failed to find repository")
	}
	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	currentUser := custommiddleware.GetUserFromContext(r)
	if currentUser == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	ticket, err := c.tickets.FindByRepositoryAndNumber(repo.ID, number)
	if err != nil || ticket == nil {
		return httperror.NotFound("ticket not found")
	}

	err = c.tickets.Reopen(ticket.ID)
	if err != nil {
		slog.Error("failed to reopen ticket", "error", err)
		return httperror.New(500, "failed to reopen ticket")
	}

	http.Redirect(w, r, "/"+owner+"/"+repoName+"/tickets/"+numberStr, http.StatusSeeOther)
	return nil
}

func (c *ticketsController) CreateComment(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")
	numberStr := chi.URLParam(r, "number")

	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		return httperror.BadRequest("invalid ticket number")
	}

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.New(500, "failed to find repository")
	}
	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	currentUser := custommiddleware.GetUserFromContext(r)
	if currentUser == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	ticket, err := c.tickets.FindByRepositoryAndNumber(repo.ID, number)
	if err != nil || ticket == nil {
		return httperror.NotFound("ticket not found")
	}

	if err := r.ParseForm(); err != nil {
		return httperror.BadRequest("invalid form data")
	}

	body := r.FormValue("body")
	if body == "" {
		// Just redirect back, don't create empty comment
		http.Redirect(w, r, "/"+owner+"/"+repoName+"/tickets/"+numberStr, http.StatusSeeOther)
		return nil
	}

	_, err = c.tickets.CreateComment(ticket.ID, currentUser.ID, body)
	if err != nil {
		slog.Error("failed to create comment", "error", err)
		return httperror.New(500, "failed to create comment")
	}

	http.Redirect(w, r, "/"+owner+"/"+repoName+"/tickets/"+numberStr, http.StatusSeeOther)
	return nil
}
