package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hyperstitieux/hypercode/config"
	"github.com/hyperstitieux/hypercode/controllers"
	"github.com/hyperstitieux/hypercode/database"
	"github.com/hyperstitieux/hypercode/database/repositories"
	"github.com/hyperstitieux/hypercode/httperror"
	custommiddleware "github.com/hyperstitieux/hypercode/middleware"
	"github.com/hyperstitieux/hypercode/public"
	"github.com/hyperstitieux/hypercode/services"
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

	authService := services.NewAuthService(users, cfg.SigningSecret)
	flashService := services.NewFlashService()

	homeController := controllers.NewHomeController(repos, users, orgs)
	signUpController := controllers.NewSignUpController(users, authService, flashService)
	signInController := controllers.NewSignInController(users, authService)
	signOutController := controllers.NewSignOutController(authService)
	settingsController := controllers.NewSettingsController(users, authService)
	forgotPasswordController := controllers.NewForgotPasswordController()
	resetPasswordController := controllers.NewResetPasswordController()
	orgsController := controllers.NewOrganizationsController()
	reposController := controllers.NewRepositoriesController(repos, users, contributors, authService, cfg.ReposBasePath)
	gitController := controllers.NewGitController(users, orgs, repos, contributors, authService, cfg.ReposBasePath)

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

	r.Get("/sign-up", wrapHandler(signUpController.Show))
	r.Post("/sign-up", wrapHandler(signUpController.Handle))

	r.Get("/sign-in", wrapHandler(signInController.Show))
	r.Post("/sign-in", wrapHandler(signInController.Handle))

	r.Get("/sign-out", wrapHandler(signOutController.Handle))

	r.Get("/settings", wrapHandler(settingsController.Show))
	r.Post("/settings/general", wrapHandler(settingsController.UpdateGeneral))
	r.Post("/settings/password", wrapHandler(settingsController.UpdatePassword))

	r.Get("/forgot-password", wrapHandler(forgotPasswordController.Show))
	r.Post("/forgot-password", wrapHandler(forgotPasswordController.Handle))

	r.Get("/reset-password", wrapHandler(resetPasswordController.Show))
	r.Post("/reset-password", wrapHandler(resetPasswordController.Handle))

	r.Get("/repositories/new", wrapHandler(reposController.Create))
	r.Post("/repositories/new", wrapHandler(reposController.Store))

	r.Get("/organizations/new", wrapHandler(orgsController.Create))
	r.Post("/organizations/new", wrapHandler(orgsController.Store))

	r.Route("/{owner}", func(r chi.Router) {
		r.Use(custommiddleware.OwnerResolver(users, orgs))

		r.Get("/", wrapHandler(orgsController.Show))
		r.Get("/settings", wrapHandler(orgsController.Settings))
		r.Put("/settings", wrapHandler(orgsController.Update))
		r.Delete("/", wrapHandler(orgsController.Delete))

		r.Route("/{repo}", func(r chi.Router) {
			r.Get("/", wrapHandler(reposController.Show))
			r.Delete("/", wrapHandler(reposController.Delete))

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
