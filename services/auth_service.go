package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/hypercodehq/hypercode/database/models"
	"github.com/hypercodehq/hypercode/database/repositories"
)

const (
	CookieNameUserID = "hypercode_user_id"
)

type AuthService interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
	SetUserCookie(w http.ResponseWriter, userID int64)
	GetUserFromCookie(r *http.Request) (*models.User, error)
	ClearUserCookie(w http.ResponseWriter)
}

type authService struct {
	users         repositories.UsersRepository
	signingSecret string
}

func NewAuthService(users repositories.UsersRepository, signingSecret string) AuthService {
	return &authService{
		users:         users,
		signingSecret: signingSecret,
	}
}

func (s *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *authService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *authService) signCookieValue(userID int64) string {
	timestamp := time.Now().Unix()
	payload := fmt.Sprintf("%d|%d", userID, timestamp)

	mac := hmac.New(sha256.New, []byte(s.signingSecret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("%s|%s", payload, signature)
}

func (s *authService) verifyCookieValue(signedValue string) (int64, error) {
	parts := strings.Split(signedValue, "|")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid cookie format")
	}

	userIDStr, timestampStr, providedSig := parts[0], parts[1], parts[2]

	payload := fmt.Sprintf("%s|%s", userIDStr, timestampStr)
	mac := hmac.New(sha256.New, []byte(s.signingSecret))
	mac.Write([]byte(payload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(providedSig), []byte(expectedSig)) != 1 {
		return 0, fmt.Errorf("invalid signature")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID")
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid timestamp")
	}

	if time.Now().Unix()-timestamp > 86400*365 {
		return 0, fmt.Errorf("cookie expired")
	}

	return userID, nil
}

func (s *authService) SetUserCookie(w http.ResponseWriter, userID int64) {
	signedValue := s.signCookieValue(userID)
	http.SetCookie(w, &http.Cookie{
		Name:     CookieNameUserID,
		Value:    signedValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 365,
	})
}

func (s *authService) GetUserFromCookie(r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie(CookieNameUserID)
	if err != nil {
		return nil, err
	}

	userID, err := s.verifyCookieValue(cookie.Value)
	if err != nil {
		return nil, err
	}

	return s.users.FindByID(userID)
}

func (s *authService) ClearUserCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   CookieNameUserID,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
