package controllers

import (
	"crypto/sha256"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hypercodehq/hypercode/database/models"
	"github.com/hypercodehq/hypercode/database/repositories"
	"github.com/hypercodehq/hypercode/middleware"
	"github.com/hypercodehq/hypercode/services"
)

type GitController interface {
	UploadPack(w http.ResponseWriter, r *http.Request) error
	ReceivePack(w http.ResponseWriter, r *http.Request) error
	InfoRefs(w http.ResponseWriter, r *http.Request) error
}

type gitController struct {
	users         repositories.UsersRepository
	orgs          repositories.OrganizationsRepository
	repos         repositories.RepositoriesRepository
	contributors  repositories.ContributorsRepository
	accessTokens  repositories.AccessTokensRepository
	authService   services.AuthService
	reposBasePath string
}

func NewGitController(
	users repositories.UsersRepository,
	orgs repositories.OrganizationsRepository,
	repos repositories.RepositoriesRepository,
	contributors repositories.ContributorsRepository,
	accessTokens repositories.AccessTokensRepository,
	authService services.AuthService,
	reposBasePath string,
) GitController {
	return &gitController{
		users:         users,
		orgs:          orgs,
		repos:         repos,
		contributors:  contributors,
		accessTokens:  accessTokens,
		authService:   authService,
		reposBasePath: reposBasePath,
	}
}

func (c *gitController) UploadPack(w http.ResponseWriter, r *http.Request) error {
	return c.handleGitOperation(w, r, "upload-pack", false)
}

func (c *gitController) ReceivePack(w http.ResponseWriter, r *http.Request) error {
	return c.handleGitOperation(w, r, "receive-pack", true)
}

func (c *gitController) InfoRefs(w http.ResponseWriter, r *http.Request) error {
	service := r.URL.Query().Get("service")
	isWriteOp := strings.Contains(service, "receive-pack")
	return c.handleGitOperation(w, r, "info-refs", isWriteOp)
}

func (c *gitController) handleGitOperation(w http.ResponseWriter, r *http.Request, operation string, isWriteOp bool) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")

	if len(repoName) >= 4 && repoName[len(repoName)-4:] == ".git" {
		repoName = repoName[:len(repoName)-4]
	}

	ownerType, _ := middleware.GetOwnerType(r.Context())
	ownerID, _ := middleware.GetOwnerID(r.Context())

	var repo *models.Repository
	var err error

	if ownerType == middleware.OwnerTypeUser {
		repo, err = c.repos.FindByUserAndName(ownerID, repoName)
	} else {
		repo, err = c.repos.FindByOrgAndName(ownerID, repoName)
	}

	if err != nil {
		http.NotFound(w, r)
		return nil
	}

	// Get the owner ID for constructing the repository path
	var ownerIDForPath string
	if ownerType == middleware.OwnerTypeUser {
		ownerIDForPath = fmt.Sprintf("%d", ownerID)
	} else {
		ownerIDForPath = fmt.Sprintf("org_%d", ownerID)
	}

	slog.Info("git http request",
		"owner", owner,
		"repo", repoName,
		"operation", operation,
		"path", r.URL.Path,
		"method", r.Method,
		"visibility", repo.Visibility,
		"isWriteOp", isWriteOp)

	if repo.Visibility != "public" || isWriteOp {
		user, _ := c.authService.GetUserFromCookie(r)

		if user == nil {
			username, password, ok := r.BasicAuth()
			slog.Info("basic auth attempt", "username", username, "hasPassword", password != "", "ok", ok)

			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Git Repository"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				slog.Warn("no credentials provided")
				return nil
			}

			authenticatedUser, err := c.users.FindByUsername(username)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="Git Repository"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				slog.Warn("user not found", "username", username)
				return nil
			}

			// Try password authentication first
			valid := authenticatedUser.Password != nil && c.authService.CheckPassword(password, *authenticatedUser.Password)

			// If password auth fails, try access token authentication
			if !valid {
				valid, err = c.authenticateWithAccessToken(authenticatedUser.ID, password)
				if err != nil {
					slog.Warn("token authentication error", "username", username, "error", err)
				}
			}

			slog.Info("authentication check", "username", username, "valid", valid)
			if !valid {
				w.Header().Set("WWW-Authenticate", `Basic realm="Git Repository"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				slog.Warn("invalid password or token", "username", username)
				return nil
			}

			user = authenticatedUser
			slog.Info("basic auth successful", "username", username)
		}

		hasAccess := false

		if repo.OwnerUserID != nil && *repo.OwnerUserID == user.ID {
			hasAccess = true
		} else if repo.OwnerOrgID != nil {
			hasAccess = true
		}

		if !hasAccess && isWriteOp {
			contributor, err := c.contributors.FindByRepositoryAndUser(repo.ID, user.ID)
			if err == nil && contributor != nil {
				if contributor.Role == "write" || contributor.Role == "admin" {
					hasAccess = true
				}
			}
		}

		if !hasAccess && !isWriteOp {
			contributor, err := c.contributors.FindByRepositoryAndUser(repo.ID, user.ID)
			if err == nil && contributor != nil {
				hasAccess = true
			}
		}

		if !hasAccess {
			http.Error(w, "Forbidden", http.StatusForbidden)
			slog.Warn("user does not have access", "user", user.Username, "owner", owner, "isWriteOp", isWriteOp)
			return nil
		}
	}

	repoPath := filepath.Join(c.reposBasePath, ownerIDForPath, fmt.Sprintf("%d", repo.ID))

	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return err
	}

	pathPrefix := fmt.Sprintf("/%s/%s", owner, repoName)
	gitPath := r.URL.Path
	if strings.HasPrefix(gitPath, pathPrefix+".git/") {
		gitPath = strings.TrimPrefix(gitPath, pathPrefix+".git")
	} else if strings.HasPrefix(gitPath, pathPrefix+"/") {
		gitPath = strings.TrimPrefix(gitPath, pathPrefix)
	}

	cmd := exec.Command("git", "http-backend")
	cmd.Dir = absRepoPath

	env := os.Environ()
	env = append(env,
		"GIT_PROJECT_ROOT="+absRepoPath,
		"GIT_HTTP_EXPORT_ALL=1",
		"PATH_INFO="+gitPath,
		"REQUEST_METHOD="+r.Method,
		"QUERY_STRING="+r.URL.RawQuery,
		"CONTENT_TYPE="+r.Header.Get("Content-Type"),
	)

	for key, values := range r.Header {
		for _, value := range values {
			cgiKey := strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
			env = append(env, "HTTP_"+cgiKey+"="+value)
		}
	}
	cmd.Env = env

	cmd.Stdin = r.Body
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Failed to execute git command", http.StatusInternalServerError)
		return err
	}

	parts := strings.SplitN(string(output), "\r\n\r\n", 2)
	if len(parts) < 2 {
		parts = strings.SplitN(string(output), "\n\n", 2)
	}

	if len(parts) == 2 {
		headerLines := strings.Split(parts[0], "\n")
		for _, line := range headerLines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			colonIdx := strings.Index(line, ":")
			if colonIdx > 0 {
				headerName := strings.TrimSpace(line[:colonIdx])
				headerValue := strings.TrimSpace(line[colonIdx+1:])

				if strings.ToLower(headerName) == "status" {
					statusParts := strings.SplitN(headerValue, " ", 2)
					if len(statusParts) > 0 {
						if statusCode := statusParts[0]; statusCode != "" && statusCode != "200" {
							w.WriteHeader(200)
						}
					}
				} else {
					w.Header().Set(headerName, headerValue)
				}
			}
		}
		w.Write([]byte(parts[1]))
	} else {
		w.Write(output)
	}

	return nil
}

// authenticateWithAccessToken checks if the provided token is valid for the user
func (c *gitController) authenticateWithAccessToken(userID int64, rawToken string) (bool, error) {
	// Hash the provided token the same way we did when storing it
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := fmt.Sprintf("%x", hash)

	// Find the token
	token, err := c.accessTokens.FindByTokenHash(tokenHash)
	if err != nil {
		return false, err
	}

	if token == nil {
		return false, nil
	}

	// Verify the token belongs to the user
	if token.UserID != userID {
		return false, nil
	}

	// Update last used timestamp (ignore errors as this is not critical)
	_ = c.accessTokens.UpdateLastUsed(token.ID)

	return true, nil
}
