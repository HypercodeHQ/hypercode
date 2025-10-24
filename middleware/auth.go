package middleware

import (
	"context"
	"net/http"

	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/hypercommit/services"
)

type contextKey string

const (
	ContextKeyUser  contextKey = "user"
	ContextKeyFlash contextKey = "flash"
)

func InjectUser(authService services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := authService.GetUserFromCookie(r)
			if err == nil && user != nil {
				ctx := context.WithValue(r.Context(), ContextKeyUser, user)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(ContextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func Auth(authService services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := authService.GetUserFromCookie(r)
			if err != nil || user == nil {
				http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
				return
			}
			ctx := context.WithValue(r.Context(), ContextKeyUser, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func InjectFlash(flashService services.FlashService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			flash := flashService.Get(r)
			if flash != nil {
				ctx := context.WithValue(r.Context(), ContextKeyFlash, flash)
				r = r.WithContext(ctx)

				flashService.Clear(w)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetFlashFromContext(r *http.Request) *services.FlashMessage {
	flash, ok := r.Context().Value(ContextKeyFlash).(*services.FlashMessage)
	if !ok {
		return nil
	}
	return flash
}
