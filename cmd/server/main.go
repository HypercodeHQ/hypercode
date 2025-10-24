package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hypercommithq/hypercommit/config"
	"github.com/hypercommithq/hypercommit/controllers"
	"github.com/hypercommithq/hypercommit/database"
	"github.com/hypercommithq/hypercommit/database/repositories"
	"github.com/hypercommithq/hypercommit/httperror"
	custommiddleware "github.com/hypercommithq/hypercommit/middleware"
	"github.com/hypercommithq/hypercommit/public"
	"github.com/hypercommithq/hypercommit/services"
)

func main() {
	cfg := config.New()

	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	users := repositories.NewUsersRepository(db.DB)
	orgs := repositories.NewOrganizationsRepository(db.DB)
	repos := repositories.NewRepositoriesRepository(db.DB)
	contributors := repositories.NewContributorsRepository(db.DB)
	stars := repositories.NewStarsRepository(db.DB)
	tickets := repositories.NewTicketsRepository(db.DB)
	accessTokens := repositories.NewAccessTokensRepository(db.DB)
	deviceAuthSessions := repositories.NewDeviceAuthSessionsRepository(db.DB)

	authService := services.NewAuthService(users, cfg.SigningSecret)
	flashService := services.NewFlashService()
	gitService := services.NewGitService(cfg.ReposBasePath)
	githubOAuthService := services.NewGitHubOAuthService(cfg.GitHubClientID, cfg.GitHubClientSecret, cfg.GitHubCallbackURL)

	homeController := controllers.NewHomeController(repos, users, orgs, stars)
	signUpController := controllers.NewSignUpController(users, authService, flashService)
	signInController := controllers.NewSignInController(users, authService)
	signOutController := controllers.NewSignOutController(authService)
	githubAuthController := controllers.NewGitHubAuthController(users, authService, githubOAuthService)
	settingsController := controllers.NewSettingsController(users, accessTokens, authService)
	accessTokensController := controllers.NewAccessTokensController(accessTokens)
	deviceAuthController := controllers.NewDeviceAuthController(deviceAuthSessions, accessTokens, users)
	forgotPasswordController := controllers.NewForgotPasswordController()
	resetPasswordController := controllers.NewResetPasswordController()
	orgsController := controllers.NewOrganizationsController(orgs, users, repos, stars, authService)
	reposController := controllers.NewRepositoriesController(repos, users, contributors, stars, orgs, authService, gitService, cfg.ReposBasePath)
	gitController := controllers.NewGitController(users, orgs, repos, contributors, accessTokens, authService, cfg.ReposBasePath)
	exploreController := controllers.NewExploreController(repos, users, orgs, stars, authService)
	ticketsController := controllers.NewTicketsController(tickets, repos, users, stars, contributors, authService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(custommiddleware.InjectUser(authService))
	r.Use(custommiddleware.InjectFlash(flashService))
	r.Use(custommiddleware.StaticFileServer(public.FileServer()))

	r.Get("/", wrapHandler(homeController.Show))

	r.Get("/auth/sign-up", wrapHandler(signUpController.Show))
	r.Post("/auth/sign-up", wrapHandler(signUpController.Handle))

	r.Get("/auth/sign-in", wrapHandler(signInController.Show))
	r.Post("/auth/sign-in", wrapHandler(signInController.Handle))

	r.Get("/auth/sign-out", wrapHandler(signOutController.Handle))

	r.Get("/auth/github", wrapHandler(githubAuthController.Login))
	r.Get("/auth/github/callback", wrapHandler(githubAuthController.Callback))

	r.Post("/api/auth/device/code", wrapHandler(deviceAuthController.InitiateDeviceAuth))
	r.Get("/api/auth/device/poll", wrapHandler(deviceAuthController.PollDeviceAuth))
	r.Get("/auth/device", wrapHandler(deviceAuthController.ShowDeviceAuthPage))
	r.Post("/auth/device/confirm", wrapHandler(deviceAuthController.ConfirmDeviceAuth))

	r.Get("/settings", wrapHandler(settingsController.Show))
	r.Post("/settings/general", wrapHandler(settingsController.UpdateGeneral))
	r.Post("/settings/password", wrapHandler(settingsController.UpdatePassword))
	r.Post("/settings/access-tokens", wrapHandler(accessTokensController.Create))
	r.Post("/settings/access-tokens/{id}/delete", wrapHandler(accessTokensController.Delete))

	r.Get("/forgot-password", wrapHandler(forgotPasswordController.Show))
	r.Post("/forgot-password", wrapHandler(forgotPasswordController.Handle))

	r.Get("/reset-password", wrapHandler(resetPasswordController.Show))
	r.Post("/reset-password", wrapHandler(resetPasswordController.Handle))

	r.Get("/repositories/new", wrapHandler(reposController.Create))
	r.Post("/repositories/new", wrapHandler(reposController.Store))

	r.Get("/organizations/new", wrapHandler(orgsController.Create))
	r.Post("/organizations/new", wrapHandler(orgsController.Store))

	r.Get("/explore/repositories", wrapHandler(exploreController.Repositories))
	r.Get("/explore/users", wrapHandler(exploreController.Users))
	r.Get("/explore/organizations", wrapHandler(exploreController.Organizations))

	r.Route("/{owner}", func(r chi.Router) {
		r.Use(custommiddleware.OwnerResolver(users, orgs))

		r.Get("/", wrapHandler(orgsController.Show))
		r.Get("/repositories", wrapHandler(orgsController.Repositories))
		r.Get("/stars", wrapHandler(orgsController.Stars))
		r.Get("/settings", wrapHandler(orgsController.Settings))
		r.Put("/settings", wrapHandler(orgsController.Update))
		r.Delete("/", wrapHandler(orgsController.Delete))

		r.Route("/{repo}", func(r chi.Router) {
			r.Get("/", wrapHandler(reposController.Show))
			r.Post("/star", wrapHandler(reposController.Star))
			r.Post("/unstar", wrapHandler(reposController.Unstar))
			r.Get("/settings", wrapHandler(reposController.Settings))
			r.Post("/settings/general", wrapHandler(reposController.UpdateSettings))
			r.Post("/settings/collaborators/add", wrapHandler(reposController.AddCollaborator))
			r.Post("/settings/collaborators/remove", wrapHandler(reposController.RemoveCollaborator))
			r.Post("/settings/collaborators/update", wrapHandler(reposController.UpdateCollaboratorRole))
			r.Post("/settings/delete", wrapHandler(reposController.Delete))

			// Tree routes - handle both with and without ref
			r.Get("/tree", wrapHandler(reposController.Tree))
			r.Get("/tree/{ref}", wrapHandler(reposController.Tree))
			r.Get("/tree/{ref}/*", wrapHandler(reposController.Tree))

			// Tickets routes
			r.Get("/tickets", wrapHandler(ticketsController.List))
			r.Get("/tickets/new", wrapHandler(ticketsController.New))
			r.Post("/tickets/new", wrapHandler(ticketsController.Create))
			r.Get("/tickets/{number}", wrapHandler(ticketsController.Show))
			r.Post("/tickets/{number}/close", wrapHandler(ticketsController.Close))
			r.Post("/tickets/{number}/reopen", wrapHandler(ticketsController.Reopen))
			r.Post("/tickets/{number}/comments", wrapHandler(ticketsController.CreateComment))

			r.Get("/info/refs", wrapHandler(gitController.InfoRefs))
			r.Post("/git-upload-pack", wrapHandler(gitController.UploadPack))
			r.Post("/git-receive-pack", wrapHandler(gitController.ReceivePack))
		})
	})

	slog.Info("starting server", "addr", cfg.HTTPAddr)

	if err := http.ListenAndServe(cfg.HTTPAddr, r); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func wrapHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			var httpErr httperror.HTTPError
			if errors.As(err, &httpErr) {
				http.Error(w, httpErr.Message, httpErr.StatusCode)
				if httpErr.StatusCode >= 500 {
					slog.Error("handler error", "error", err, "path", r.URL.Path, "method", r.Method)
				}
				return
			}

			slog.Error("handler error", "error", err, "path", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
