package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hyperstitieux/hypercode/database/repositories"
)

const (
	OwnerTypeKey contextKey = "ownerType"
	OwnerIDKey   contextKey = "ownerID"
)

type OwnerType string

const (
	OwnerTypeUser OwnerType = "user"
	OwnerTypeOrg  OwnerType = "org"
)

func OwnerResolver(usersRepo repositories.UsersRepository, orgsRepo repositories.OrganizationsRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			owner := chi.URLParam(r, "owner")
			if owner == "" {
				next.ServeHTTP(w, r)
				return
			}

			user, err := usersRepo.FindByUsername(owner)
			if err == nil && user != nil {
				ctx := context.WithValue(r.Context(), OwnerTypeKey, OwnerTypeUser)
				ctx = context.WithValue(ctx, OwnerIDKey, user.ID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			org, err := orgsRepo.FindByUsername(owner)
			if err == nil && org != nil {
				ctx := context.WithValue(r.Context(), OwnerTypeKey, OwnerTypeOrg)
				ctx = context.WithValue(ctx, OwnerIDKey, org.ID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			http.Error(w, "owner not found", http.StatusNotFound)
		})
	}
}

func GetOwnerType(ctx context.Context) (OwnerType, bool) {
	ownerType, ok := ctx.Value(OwnerTypeKey).(OwnerType)
	return ownerType, ok
}

func GetOwnerID(ctx context.Context) (int64, bool) {
	ownerID, ok := ctx.Value(OwnerIDKey).(int64)
	return ownerID, ok
}
