package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GitHubOAuthService interface {
	GetAuthURL(state string) string
	ExchangeCode(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*GitHubUser, error)
}

type GitHubUser struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
}

type githubOAuthService struct {
	config *oauth2.Config
}

func NewGitHubOAuthService(clientID, clientSecret, callbackURL string) GitHubOAuthService {
	return &githubOAuthService{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  callbackURL,
			Scopes:       []string{"user:email", "read:user"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (s *githubOAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (s *githubOAuthService) ExchangeCode(code string) (*oauth2.Token, error) {
	ctx := context.Background()
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

func (s *githubOAuthService) GetUserInfo(token *oauth2.Token) (*GitHubUser, error) {
	client := s.config.Client(context.Background(), token)

	// Get user info
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var user GitHubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	// If email is not public, fetch from emails endpoint
	if user.Email == "" {
		emailsResp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			return nil, fmt.Errorf("failed to get user emails: %w", err)
		}
		defer emailsResp.Body.Close()

		emailsBody, err := io.ReadAll(emailsResp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read emails response: %w", err)
		}

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		if err := json.Unmarshal(emailsBody, &emails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal emails: %w", err)
		}

		// Find primary verified email
		for _, email := range emails {
			if email.Primary && email.Verified {
				user.Email = email.Email
				break
			}
		}
	}

	return &user, nil
}
