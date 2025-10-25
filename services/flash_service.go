package services

import (
	"net/http"

	"github.com/hypercommithq/hypercommit/httputil"
)

const (
	CookieNameFlash = "hypercommit_flash"
)

type FlashType string

const (
	FlashCelebration FlashType = "celebration"
	FlashSuccess     FlashType = "success"
	FlashError       FlashType = "error"
	FlashInfo        FlashType = "info"
)

type FlashMessage struct {
	Type FlashType
}

type FlashService interface {
	Set(w http.ResponseWriter, r *http.Request, flashType FlashType)
	Get(r *http.Request) *FlashMessage
	Clear(w http.ResponseWriter)
}

type flashService struct{}

func NewFlashService() FlashService {
	return &flashService{}
}

func (s *flashService) Set(w http.ResponseWriter, r *http.Request, flashType FlashType) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieNameFlash,
		Value:    string(flashType),
		Path:     "/",
		HttpOnly: true,
		Secure:   httputil.IsHTTPS(r),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   10,
	})
}

func (s *flashService) Get(r *http.Request) *FlashMessage {
	cookie, err := r.Cookie(CookieNameFlash)
	if err != nil {
		return nil
	}

	return &FlashMessage{
		Type: FlashType(cookie.Value),
	}
}

func (s *flashService) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   CookieNameFlash,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
